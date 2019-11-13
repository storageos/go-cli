package describe

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

type nodeCommand struct {
	config  ConfigProvider
	client  DescribeClient
	display DescribeDisplayer

	writer io.Writer
}

func (c *nodeCommand) run(cmd *cobra.Command, args []string) error {
	timeout, err := c.config.DialTimeout()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	switch len(args) {
	case 1:
		return c.describeNode(ctx, args)
	default:
		return c.listNodes(ctx, args)
	}
}

func (c *nodeCommand) describeNode(ctx context.Context, args []string) error {
	uid := id.Node(args[0])

	node, err := c.client.DescribeNode(ctx, uid)
	if err != nil {
		return err
	}

	return c.display.DescribeNode(ctx, c.writer, node)
}

func (c *nodeCommand) listNodes(ctx context.Context, args []string) error {
	uids := make([]id.Node, len(args))
	for i, a := range args {
		uids[i] = id.Node(a)
	}

	nodes, err := c.client.DescribeListNodes(ctx, uids...)
	if err != nil {
		return err
	}

	return c.display.DescribeNodeList(ctx, c.writer, nodes)
}

func newNode(w io.Writer, client DescribeClient, config ConfigProvider) *cobra.Command {
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
		Use:     "node [node ids...]",
		Short:   "node retrieves detailed information about StorageOS nodes",
		Example: `
$ storageos describe node banana
`,

		RunE: c.run,
	}

	return cobraCommand
}
