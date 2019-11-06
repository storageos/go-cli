package describe

import (
	"io"

	"github.com/spf13/cobra"
)

type clusterCommand struct {
	client  DescribeClient
	display DescribeDisplayer

	writer io.Writer
}

func (c *clusterCommand) run(cmd *cobra.Command, _ []string) error {
	cluster, err := c.client.DescribeCluster()
	if err != nil {
		return err
	}

	return c.display.WriteDescribeCluster(c.writer, cluster)
}

func newCluster(w io.Writer, client DescribeClient, display DescribeDisplayer) *cobra.Command {
	c := &clusterCommand{
		client:  client,
		display: display,

		writer: w,
	}
	cobraCommand := &cobra.Command{
		Use:   "cluster",
		Short: "cluster retrieves detailed information about the StorageOS cluster",
		Example: `
$ storageos describe cluster
`,

		RunE: c.run,
	}

	return cobraCommand
}
