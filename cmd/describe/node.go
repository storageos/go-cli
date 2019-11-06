package describe

import (
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/pkg/id"
)

type nodeCommand struct {
	client  DescribeClient
	display DescribeDisplayer

	writer io.Writer
}

func (c *nodeCommand) run(cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 1:
		return c.describeNode(cmd, args)
	default:
		return c.listNodes(cmd, args)
	}
}

func (c *nodeCommand) describeNode(_ *cobra.Command, args []string) error {
	uid := id.Node(args[0])

	node, err := c.client.DescribeNode(uid)
	if err != nil {
		return err
	}

	return c.display.WriteDescribeNode(c.writer, node)
}

func (c *nodeCommand) listNodes(_ *cobra.Command, args []string) error {
	uids := make([]id.Node, len(args))
	for i, a := range args {
		uids[i] = id.Node(a)
	}

	nodes, err := c.client.DescribeListNodes(uids...)
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
