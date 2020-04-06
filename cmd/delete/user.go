package delete

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

type userCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	// useCAS determines whether the command makes the delete request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string

	writer io.Writer
}

func (c *userCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {

	var userID id.User

	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	if useIDs {
		userID = id.User(args[0])
	} else {
		u, err := c.client.GetUserByName(ctx, args[0])
		if err != nil {
			return err
		}
		userID = u.ID
	}

	params := &apiclient.DeleteUserRequestParams{}
	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}

	err = c.client.DeleteUser(
		ctx,
		userID,
		params,
	)
	if err != nil {
		return err
	}

	return c.display.DeleteUser(ctx, c.writer, output.UserDeletion{ID: userID})
}

func newUser(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &userCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "user [user name]",
		Short: "Delete a user",
		Example: `
$ storageos delete user my-unneeded-user
$ storageos delete user --use-ids my-user-id
`,

		Args: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("must specify exactly one user for deletion")
			}
			return nil
		}),

		PreRunE: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, args []string) error {
			c.display = SelectDisplayer(c.config)

			return nil
		}),

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	c.useCAS = flagutil.SupportCAS(cobraCommand.Flags(), &c.casVersion)

	return cobraCommand
}
