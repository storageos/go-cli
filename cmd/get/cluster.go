package get

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
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

	return c.display.GetCluster(ctx, c.writer, cluster)
}

func newCluster(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &clusterCommand{
		config: config,
		client: client,
		display: jsonformat.NewDisplayer(
			jsonformat.DefaultEncodingIndent,
		),

		writer: w,
	}
	cobraCommand := &cobra.Command{
		Use:   "cluster",
		Short: "cluster retrieves basic information about the StorageOS cluster",
		Example: `
$ storageos get cluster
`,

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.RunWithTimeout(c.config)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		// If a legitimate error occurs as part of the VERB cluster command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	return cobraCommand
}
