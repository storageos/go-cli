package get

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/config"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

type nodeCommand struct {
	client  GetClient
	display GetDisplayer

	writer io.Writer
}

func (c *nodeCommand) run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultCommandTimeout)
	defer cancel()

	switch len(args) {
	case 1:
		return c.getNode(ctx, cmd, args)
	default:
		return c.listNodes(ctx, cmd, args)
	}
}

func (c *nodeCommand) getNode(ctx context.Context, _ *cobra.Command, args []string) error {
	uid := id.Node(args[0])

	node, err := c.client.GetNode(ctx, uid)
	if err != nil {
		return err
	}

	return c.display.WriteGetNode(c.writer, node)
}

func (c *nodeCommand) listNodes(ctx context.Context, _ *cobra.Command, args []string) error {
	uids := make([]id.Node, len(args))
	for i, a := range args {
		uids[i] = id.Node(a)
	}

	nodes, err := c.client.GetListNodes(ctx, uids...)
	if err != nil {
		return err
	}

	return c.display.WriteGetNodeList(c.writer, nodes)
}

func newNode(w io.Writer, client GetClient, display GetDisplayer) *cobra.Command {
	c := &nodeCommand{
		client:  client,
		display: display,

		writer: w,
	}
	cobraCommand := &cobra.Command{
		Aliases: []string{"nodes"},
		Use:     "node [node ids...]",
		Short:   "node retrieves basic information about StorageOS nodes",
		Example: `
$ storageos get node banana
`,

		RunE: c.run,
	}

	return cobraCommand
}
