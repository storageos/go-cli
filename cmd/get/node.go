package get

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

type nodeCommand struct {
	config  ConfigProvider
	client  GetClient
	display GetDisplayer

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
		return c.getNode(ctx, args)
	default:
		return c.listNodes(ctx, args)
	}
}

func (c *nodeCommand) getNode(ctx context.Context, args []string) error {
	uid := id.Node(args[0])

	node, err := c.client.GetNode(ctx, uid)
	if err != nil {
		return err
	}

	return c.display.GetNode(ctx, c.writer, node)
}

func (c *nodeCommand) listNodes(ctx context.Context, args []string) error {
	uids := make([]id.Node, len(args))
	for i, a := range args {
		uids[i] = id.Node(a)
	}

	nodes, err := c.client.GetListNodes(ctx, uids...)
	if err != nil {
		return err
	}

	return c.display.GetNodeList(ctx, c.writer, nodes)
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
		Use:     "node [node ids...]",
		Short:   "node retrieves basic information about StorageOS nodes",
		Example: `
$ storageos get node banana
`,

		RunE: c.run,

		// If a legitimate error occurs as part of the VERB node command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	return cobraCommand
}
