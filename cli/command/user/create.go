package user

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"syscall"

	"context"
	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"golang.org/x/crypto/ssh/terminal"
)

type stringSlice []string

func (s *stringSlice) Type() string {
	return "string"
}

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSlice) Set(val string) error {
	*s = append(*s, strings.Split(val, ",")...)
	return nil
}

type createOptions struct {
	username string
	password bool
	groups   stringSlice
	role     string
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := createOptions{}

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] [USERNAME]",
		Short: `Create a new User, E.g. "storageos user create --password alice" (interactive password prompt)`,
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				if opt.username != "" {
					fmt.Fprint(storageosCli.Err(), "Conflicting options: either specify --username or provide positional arg, not both\n")
					return cli.StatusError{StatusCode: 1}
				}
				opt.username = args[0]
			}
			return runCreate(storageosCli, opt)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opt.username, "username", "", "Username")
	flags.Lookup("username").Hidden = true
	flags.BoolVar(&opt.password, "password", false, "Prompt for password (interactive)")
	flags.StringVar(&opt.role, "role", "user", "Role")
	flags.Var(&opt.groups, "groups", "Groups")

	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opt createOptions) error {
	var password string

	if opt.password {
		var err error

		password, err = getPassword(storageosCli)
		if err != nil {
			return err
		}
	}

	opt.role = strings.ToLower(opt.role)

	if !verify(storageosCli, opt) {
		return errors.New("Input failed verification")
	}

	client := storageosCli.Client()

	params := types.UserCreateOptions{
		Username: opt.username,
		Password: password,
		Groups:   opt.groups,
		Role:     opt.role,
		Context:  context.Background(),
	}

	err := client.UserCreate(params)
	return err
}

func getPassword(storageosCli *command.StorageOSCli) (string, error) {
retry:
	fmt.Fprint(storageosCli.Out(), "Password: ")
	passBytes1, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	fmt.Fprint(storageosCli.Out(), "\nConfirm Password: ")
	passBytes2, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	if string(passBytes1) != string(passBytes2) {
		fmt.Fprintln(storageosCli.Err(), "\nPasswords don't match, please retry...")
		goto retry
	}

	return string(passBytes1), nil
}

func verify(storageosCli *command.StorageOSCli, opt createOptions) (verifies bool) {
	verifies = true

	if !verifyUsername(opt.username) {
		verifies = false
		fmt.Fprintln(storageosCli.Err(), `Username doesn't follow format "[a-zA-Z0-9]+"`)
	}

	if i, pass := verifyGroups(opt.groups); !pass {
		verifies = false
		fmt.Fprintf(storageosCli.Err(), `Group element %d doesn't follow format "[a-zA-Z0-9]+"\n`, i)
	}

	if !verifyRole(opt.role) {
		verifies = false
		fmt.Fprintf(storageosCli.Err(), `Role must me "user" or "admin", not %s\n`, opt.role)
	}

	return
}

func verifyUsername(username string) bool {
	return regexp.MustCompile(`[a-zA-Z0-9]+`).MatchString(username)
}

func verifyGroups(groups []string) (index int, pass bool) {
	matcher := regexp.MustCompile(`[a-zA-Z0-9]+`)
	for i, group := range groups {
		if !matcher.MatchString(group) {
			return i, false
		}
	}
	return -1, true
}

func verifyRole(role string) bool {
	return role == "admin" || role == "user"
}
