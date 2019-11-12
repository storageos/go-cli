package get

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/config"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
)

type clusterCommand struct {
	client  GetClient
	display GetDisplayer

	writer io.Writer
}

func (c *clusterCommand) run(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultDialTimeout)
	defer cancel()

	cluster, err := c.client.GetCluster(ctx)
	if err != nil {
		return err
	}

	return c.display.GetCluster(c.writer, cluster)
}

func newCluster(w io.Writer, client GetClient) *cobra.Command {
	c := &clusterCommand{
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
	}

	return cobraCommand
}
