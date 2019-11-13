package get

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/output/jsonformat"
)

type clusterCommand struct {
	config  ConfigProvider
	client  GetClient
	display GetDisplayer

	writer io.Writer
}

func (c *clusterCommand) run(cmd *cobra.Command, _ []string) error {
	timeout, err := c.config.DialTimeout()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cluster, err := c.client.GetCluster(ctx)
	if err != nil {
		return err
	}

	return c.display.GetCluster(ctx, c.writer, cluster)
}

func newCluster(w io.Writer, client GetClient, config ConfigProvider) *cobra.Command {
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

		RunE: c.run,

		// If a legitimate error occurs as part of the VERB cluster command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	return cobraCommand
}
