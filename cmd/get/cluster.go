package get

import (
	"io"

	"github.com/spf13/cobra"
)

type clusterCommand struct {
	client  GetClient
	display GetDisplayer

	writer io.Writer
}

func (c *clusterCommand) run(cmd *cobra.Command, _ []string) error {
	cluster, err := c.client.GetCluster()
	if err != nil {
		return err
	}

	return c.display.WriteGetCluster(c.writer, cluster)
}

func newCluster(w io.Writer, client GetClient, display GetDisplayer) *cobra.Command {
	c := &clusterCommand{
		client:  client,
		display: display,

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
