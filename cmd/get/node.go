package get

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

type nodeCommand struct {
	config  ConfigProvider
	client  GetClient
	display GetDisplayer

	usingID bool

	writer io.Writer
}

func (c *nodeCommand) run(cmd *cobra.Command, args []string) error {
	timeout, err := c.config.CommandTimeout()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	switch len(args) {
	case 1:
		n, err := c.getNode(ctx, args[0])
		if err != nil {
			return err
		}

		return c.display.GetNode(ctx, c.writer, n)
	default:
		nodes, err := c.listNodes(ctx, args)
		if err != nil {
			return err
		}

		return c.display.GetNodeList(ctx, c.writer, nodes)
	}
}

// getNode retrieves a single node resource using the API client, determining
// whether to retrieve the node by name or ID based on the current command configuration.
func (c *nodeCommand) getNode(ctx context.Context, ref string) (*node.Resource, error) {
	if !c.usingID {
		return c.client.GetNodeByName(ctx, ref)
	}

	uid := id.Node(ref)
	return c.client.GetNode(ctx, uid)
}

// listNodes retrieves a list of node resources using the API client, determining
// whether to retrieve nodes by names by name or ID based on the current
// command configuration.
func (c *nodeCommand) listNodes(ctx context.Context, refs []string) ([]*node.Resource, error) {
	if !c.usingID {
		return c.client.GetListNodesByName(ctx, refs...)
	}

	uids := make([]id.Node, len(refs))
	for i, a := range refs {
		uids[i] = id.Node(a)
	}

	return c.client.GetListNodes(ctx, uids...)
}

func newNode(w io.Writer, client GetClient, config ConfigProvider) *cobra.Command {
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
		Short:   "node retrieves basic information about StorageOS nodes",
		Example: `
$ storageos get node banana
`,

		RunE: c.run,

		// If a legitimate error occurs as part of the VERB node command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	cobraCommand.Flags().BoolVar(&c.usingID, "use-id", false, "request StorageOS nodes by ID instead of by name")

	return cobraCommand
}