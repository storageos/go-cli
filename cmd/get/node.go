package get

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/selectors"
)

type nodeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	selectors []string

	writer io.Writer
}

func (c *nodeCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 1:
		n, err := c.getNode(ctx, args[0])
		if err != nil {
			return err
		}

		return c.display.GetNode(ctx, c.writer, n)
	default:
		set, err := selectors.NewSetFromStrings(c.selectors...)
		if err != nil {
			return err
		}

		nodes, err := c.listNodes(ctx, args)
		if err != nil {
			return err
		}

		return c.display.GetListNodes(ctx, c.writer, set.FilterNodes(nodes))
	}
}

// getNode retrieves a single node resource using the API client, determining
// whether to retrieve the node by name or ID based on the current command configuration.
func (c *nodeCommand) getNode(ctx context.Context, ref string) (*node.Resource, error) {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return nil, err
	}

	if !useIDs {
		return c.client.GetNodeByName(ctx, ref)
	}

	uid := id.Node(ref)
	return c.client.GetNode(ctx, uid)
}

// listNodes retrieves a list of node resources using the API client, determining
// whether to retrieve nodes by names by name or ID based on the current
// command configuration.
func (c *nodeCommand) listNodes(ctx context.Context, refs []string) ([]*node.Resource, error) {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return nil, err
	}

	if !useIDs {
		return c.client.GetListNodesByName(ctx, refs...)
	}

	uids := make([]id.Node, len(refs))
	for i, a := range refs {
		uids[i] = id.Node(a)
	}

	return c.client.GetListNodes(ctx, uids...)
}

func newNode(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &nodeCommand{
		config: config,
		client: client,
		writer: w,
	}
	cobraCommand := &cobra.Command{
		Aliases: []string{"nodes"},
		Use:     "node [node names...]",
		Short:   "node retrieves basic information about StorageOS nodes",
		Example: `
$ storageos get node my-node-name
`,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			c.display = SelectDisplayer(c.config)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.EnsureTargetOrSelectors(&c.selectors),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		// If a legitimate error occurs as part of the VERB node command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	flagutil.SupportSelectors(cobraCommand.Flags(), &c.selectors)

	return cobraCommand
}
