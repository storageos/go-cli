package get

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/policygroup"
)

type policyGroupCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	writer io.Writer
}

func (c *policyGroupCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	namespaces, err := c.client.ListNamespaces(ctx)
	if err != nil {
		return err
	}

	switch len(args) {

	case 1:
		group, err := c.getPolicyGroup(ctx, args[0], useIDs)
		if err != nil {
			return err
		}

		return c.display.GetPolicyGroup(ctx, c.writer, output.NewPolicyGroup(group, namespaces))

	default:
		groups, err := c.listPolicyGroups(ctx, args, useIDs)
		if err != nil {
			return err
		}

		return c.display.GetListPolicyGroups(
			ctx,
			c.writer,
			output.NewPolicyGroups(groups, namespaces),
		)
	}
}

// getPolicyGroup retrieves a single policy group resource using the API client,
// by ID or by name, depending on the useIDs bool
func (c *policyGroupCommand) getPolicyGroup(ctx context.Context, ref string, useIDs bool) (*policygroup.Resource, error) {
	if !useIDs {
		return c.client.GetPolicyGroupByName(ctx, ref)
	}

	uid := id.PolicyGroup(ref)
	return c.client.GetPolicyGroup(ctx, uid)
}

// listPolicyGroups retrieves a list of policy group resources using the API
// client, by ID or by name, depending on the useIDs bool
func (c *policyGroupCommand) listPolicyGroups(ctx context.Context, refs []string, useIDs bool) ([]*policygroup.Resource, error) {
	if !useIDs {
		return c.client.GetListPolicyGroupsByName(ctx, refs...)
	}

	uids := make([]id.PolicyGroup, len(refs))
	for i, ref := range refs {
		uids[i] = id.PolicyGroup(ref)
	}

	return c.client.GetListPolicyGroupsByUID(ctx, uids...)
}

func newPolicyGroup(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &policyGroupCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Aliases: []string{"policy-groups"},
		Use:     "policy-group [policy-group names...]",
		Short:   "Retrieve basic details of policy groups",
		Example: `
$ storageos get policy-group
$ storageos get policy-group my-policy-group-name
$ storageos get policy-group --use-ids my-policy-group-id
`,
		PreRun: func(_ *cobra.Command, _ []string) {
			c.display = SelectDisplayer(c.config)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	return cobraCommand
}
