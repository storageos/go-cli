package create

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

var (
	errUsernameArgRequired = errors.New("Username argument required")
	errConflictingUsername = errors.New("Conflicting usernames provided: specify either via the flag or the argument")

	errPasswordTooShort     = errors.New("Provided password must have at least 8 characters")
	errUserPasswordMismatch = errors.New("Provided passwords do not match")
)

type userCommand struct {
	config  ConfigProvider
	client  CreateClient
	display CreateDisplayer

	username      string
	password      string
	createAsAdmin bool
	groups        []string

	writer io.Writer
}

func (c *userCommand) run(cmd *cobra.Command, _ []string) error {
	timeout, err := c.config.CommandTimeout()
	if err != nil {
		return err
	}

	// If there is no password available when running the command interactively
	// prompt for one.
	if c.password == "" {
		p, err := c.promptForPassword()
		if err != nil {
			return err
		}
		c.password = p
	}

	if len(c.password) < 8 {
		return errPasswordTooShort
	}

	// Start the timeout after the prompt interaction has been completed.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return c.createUser(ctx)
}

func (c *userCommand) createUser(ctx context.Context) error {
	groupIDs := make([]id.PolicyGroup, len(c.groups))
	for i, g := range c.groups {
		groupIDs[i] = id.PolicyGroup(g)
	}

	user, err := c.client.CreateUser(
		ctx,
		c.username,
		c.password,
		c.createAsAdmin,
		groupIDs...,
	)
	if err != nil {
		return err
	}

	return c.display.CreateUser(ctx, c.writer, user)
}

// promptForPassword will interactively request a password from the user,
// rejecting blank responses.
func (c *userCommand) promptForPassword() (string, error) {
	fmt.Fprint(c.writer, "Password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(c.writer)
	if err != nil {
		return "", err
	}

	fmt.Fprint(c.writer, "Confirm Password: ")
	confirmation, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(c.writer)
	if err != nil {
		return "", err
	}
	if !bytes.Equal(password, confirmation) {
		return "", errUserPasswordMismatch
	}

	return string(password), nil
}

// newUser builds a cobra command from the provided arguments for requesting the
// creation of a StorageOS user account.
func newUser(w io.Writer, client CreateClient, config ConfigProvider) *cobra.Command {
	c := &userCommand{
		config: config,
		client: client,
		display: jsonformat.NewDisplayer(
			jsonformat.DefaultEncodingIndent,
		),

		writer: w,
	}
	cobraCommand := &cobra.Command{
		Use:   "user",
		Short: "user requests the creation of a new StorageOS user account",
		Example: `
$ storageos create user --with-username --with-password --with-admin --with-groups
		`,

		Args: func(_ *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				if c.username == "" {
					return errUsernameArgRequired
				}
			case 1:
				if c.username != "" {
					return errConflictingUsername
				}
				c.username = args[0]
			default:
				return errors.New("too many arguments")
			}

			return nil
		},

		RunE: c.run,

		SilenceUsage: true,
	}

	cobraCommand.Flags().StringVar(&c.username, "with-username", "", "the username to assign to the StorageOS user account being created")
	cobraCommand.Flags().StringVar(&c.password, "with-password", "", "the password to assign to the StorageOS user account being created. If not specified, this will be prompted for.")
	cobraCommand.Flags().BoolVar(&c.createAsAdmin, "with-admin", false, "controls whether the StorageOS user account being created is given administrative priviledges")
	cobraCommand.Flags().StringArrayVar(&c.groups, "with-groups", []string{}, "the list of policy group IDs to assign to the StorageOS user account being created")

	return cobraCommand
}
