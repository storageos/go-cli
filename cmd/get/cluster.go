package get

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output"
)

type clusterCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	writer io.Writer
}

func (c *clusterCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, _ []string) error {
	cluster, err := c.client.GetCluster(ctx)
	if err != nil {
		return err
	}

	nodes, err := c.client.ListNodes(ctx)
	if err != nil {
		return err
	}

	return c.display.GetCluster(ctx, c.writer, output.NewCluster(cluster, nodes))
}

func newCluster(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &clusterCommand{
		config: config,
		client: client,
		writer: w,
	}
	cobraCommand := &cobra.Command{
		Use:   "cluster",
		Short: "Fetch cluster-wide configuration details",
		Example: `
$ storageos get cluster
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

		// If a legitimate error occurs as part of the VERB cluster command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	return cobraCommand
}
