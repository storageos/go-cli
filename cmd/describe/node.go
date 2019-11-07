package describe

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/config"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

type nodeCommand struct {
	client  DescribeClient
	display DescribeDisplayer

	writer io.Writer
}

func (c *nodeCommand) run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultCommandTimeout)
	defer cancel()

	switch len(args) {
	case 1:
		return c.describeNode(ctx, cmd, args)
	default:
		return c.listNodes(ctx, cmd, args)
	}
}

func (c *nodeCommand) describeNode(ctx context.Context, _ *cobra.Command, args []string) error {
	uid := id.Node(args[0])

	node, err := c.client.DescribeNode(ctx, uid)
	if err != nil {
		return err
	}

	return c.display.WriteDescribeNode(c.writer, node)
}

func (c *nodeCommand) listNodes(ctx context.Context, _ *cobra.Command, args []string) error {
	uids := make([]id.Node, len(args))
	for i, a := range args {
		uids[i] = id.Node(a)
	}

	nodes, err := c.client.DescribeListNodes(ctx, uids...)
	if err != nil {
		return err
	}

	return c.display.WriteDescribeNodeList(c.writer, nodes)
}

func newNode(w io.Writer, client DescribeClient, display DescribeDisplayer) *cobra.Command {
	c := &nodeCommand{
		client:  client,
		display: display,

		writer: w,
	}
	cobraCommand := &cobra.Command{
		Aliases: []string{"nodes"},
		Use:     "node [node ids...]",
		Short:   "node retrieves detailed information about StorageOS nodes",
		Example: `
$ storageos describe node banana
`,

		RunE: c.run,
	}

	return cobraCommand
}
