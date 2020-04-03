package get

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/policygroup"
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
	case 0:
		users, err := c.client.ListUsers(ctx)
		if err != nil {
			return err
		}

		outputUsers, err := c.toOutputUsers(ctx, users)
		if err != nil {
			return err
		}

		return c.display.GetUsers(ctx, c.writer, outputUsers)

	case 1:
		var u *user.Resource

		if useIDs {
			uID := id.User(args[0])
			u, err = c.client.GetUser(ctx, uID)
			if err != nil {
				return err
			}
		} else {
			u, err = c.client.GetUserByName(ctx, args[0])
			if err != nil {
				return err
			}
		}

		policyGroups, err := c.client.GetListPolicyGroupsByUID(ctx, u.Groups...)
		if err != nil {
			return err
		}

		groupMapping := make(map[id.PolicyGroup]*policygroup.Resource)
		for _, p := range policyGroups {
			groupMapping[p.ID] = p
		}

		outputUser, err := output.NewUser(u, groupMapping)
		if err != nil {
			return err
		}

		return c.display.GetUser(ctx, c.writer, outputUser)

	default: // more than 1
		var users []*user.Resource

		if useIDs {

			ids := make([]id.User, 0, len(args))
			for _, arg := range args {
				ids = append(ids, id.User(arg))
			}

			users, err = c.client.GetListUsersByUID(ctx, ids)
			if err != nil {
				return err
			}

		} else {
			users, err = c.client.GetListUsersByUsername(ctx, args)
			if err != nil {
				return err
			}
		}

		outputUsers, err := c.toOutputUsers(ctx, users)
		if err != nil {
			return err
		}

		return c.display.GetUsers(ctx, c.writer, outputUsers)
	}
}

func (c *userCommand) toOutputUsers(ctx context.Context, users []*user.Resource) ([]*output.User, error) {
	groups := make([]id.PolicyGroup, 0)
	for _, u := range users {
		groups = append(groups, u.Groups...)
	}

	policyGroups, err := c.client.GetListPolicyGroupsByUID(ctx, groups...)
	if err != nil {
		return nil, err
	}

	groupMapping := make(map[id.PolicyGroup]*policygroup.Resource)
	for _, p := range policyGroups {
		groupMapping[p.ID] = p
	}

	outputUsers, err := output.NewUsers(users, groupMapping)
	if err != nil {
		return nil, err
	}

	return outputUsers, nil
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
		Short:   "Fetch user details",
		Example: `
$ storageos get user my-username
$ storageos get user my-username-1 my-username-2
$ storageos get user --use-ids my-userid
$ storageos get user --use-ids my-userid-1 my-userid-2
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

		// If a legitimate error occurs as part of the VERB user command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	return cobraCommand
}
