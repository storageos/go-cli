package describe

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
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
		n, err := c.describeNode(ctx, args[0])
		if err != nil {
			return err
		}

		return c.display.DescribeNode(ctx, c.writer, n)
	default:
		set, err := selectors.NewSetFromStrings(c.selectors...)
		if err != nil {
			return err
		}

		nodes, err := c.listNodes(ctx, args)
		if err != nil {
			return err
		}

		return c.display.DescribeListNodes(ctx, c.writer, set.FilterNodeStates(nodes))
	}
}

// describeNode retrieves a single node state using the API client, determining
// whether to retrieve the node by name or ID based on the current command configuration.
func (c *nodeCommand) describeNode(ctx context.Context, ref string) (*node.State, error) {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return nil, err
	}

	if !useIDs {
		return c.client.DescribeNodeByName(ctx, ref)
	}

	uid := id.Node(ref)
	return c.client.DescribeNode(ctx, uid)
}

// listNodes retrieves a list of node states using the API client, determining
// whether to retrieve nodes by names by name or ID based on the current
// command configuration.
func (c *nodeCommand) listNodes(ctx context.Context, refs []string) ([]*node.State, error) {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return nil, err
	}

	if !useIDs {
		return c.client.DescribeListNodesByName(ctx, refs...)
	}

	uids := make([]id.Node, len(refs))
	for i, ref := range refs {
		uids[i] = id.Node(ref)
	}

	return c.client.DescribeListNodes(ctx, uids...)
}

func newNode(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &nodeCommand{
		config: config,
		client: client,
		display: jsonformat.NewDisplayer(
			jsonformat.DefaultEncodingIndent,
		),

		writer: w,
	}
	cobraCommand := &cobra.Command{
		Aliases: []string{"nodes"},
		Use:     "node [node names...]",
		Short:   "Retrieve detailed information for nodes in the cluster",
		Example: `
$ storageos describe node my-node-name
`,

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.EnsureTargetOrSelectors(&c.selectors),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	flagutil.SupportSelectors(cobraCommand.Flags(), &c.selectors)

	return cobraCommand
}
