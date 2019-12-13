package get

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
)

type clusterCommand struct {
	config  ConfigProvider
	client  GetClient
	display GetDisplayer

	writer io.Writer
}

func (c *clusterCommand) run(cmd *cobra.Command, _ []string) error {
	timeout, err := c.config.CommandTimeout()
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

func newCluster(w io.Writer, initClient func() (*apiclient.Client, error), config ConfigProvider) *cobra.Command {
	c := &clusterCommand{
		config: config,
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

		PreRunE: func(_ *cobra.Command, _ []string) error {
			client, err := initClient()
			if err != nil {
				return fmt.Errorf("error initialising api client: %w", err)
			}
			c.client = client
			return nil
		},
		RunE: c.run,

		// If a legitimate error occurs as part of the VERB cluster command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	return cobraCommand
}
