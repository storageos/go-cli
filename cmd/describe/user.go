package describe

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/user"
)

type userCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer
	writer  io.Writer
}

func (c *userCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	switch len(args) {
	case 1:
		var u *user.Resource
		var err error

		if useIDs {
			u, err = c.client.GetUser(ctx, id.User(args[0]))
		} else {
			u, err = c.client.GetUserByName(ctx, args[0])
		}
		if err != nil {
			return err
		}

		policyGroups, err := c.client.GetListPolicyGroupsByUID(ctx, u.Groups...)
		if err != nil {
			return err
		}

		return c.display.DescribeUser(ctx, c.writer, output.NewUser(u, policyGroups))

	default:
		// get all users
		users, err := c.client.ListUsers(ctx)
		if err != nil {
			return err
		}

		// get a merged list of all policy groups
		groups := make([]id.PolicyGroup, 0)
		for _, u := range users {
			groups = append(groups, u.Groups...)
		}

		policyGroups, err := c.client.GetListPolicyGroupsByUID(ctx, groups...)
		if err != nil {
			return err
		}

		outputUsers := make([]*output.User, 0, len(users))
		for _, u := range users {
			outputUsers = append(outputUsers, output.NewUser(u, policyGroups))
		}

		return c.display.DescribeListUsers(ctx, c.writer, outputUsers)
	}
}

func newUser(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &userCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Aliases: []string{"users"},
		Use:     "user [user names...]",
		Short:   "Show detailed information for users",
		Example: `
$ storageos describe users
$ storageos describe user my-username
$ storageos describe user my-username-1 my-username-2
$ storageos describe user --use-ids my-userid
$ storageos describe user --use-ids my-userid-1 my-userid-2
`,

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

		// If a legitimate error occurs as part of the VERB volume command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	return cobraCommand
}
