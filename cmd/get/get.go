package get

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/pkg/cluster"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/node"
	"code.storageos.net/storageos/c2-cli/pkg/volume"
)

// GetClient defines the functionality required by the CLI application to
// reasonably implement the "get" verb commands.
type GetClient interface {
	GetCluster() (*cluster.Resource, error)
	GetNode(id.Node) (*node.Resource, error)
	GetListNodes(...id.Node) ([]*node.Resource, error)
	GetVolume(id.Namespace, id.Volume) (*volume.Resource, error)
}

// GetDisplayer defines the functionality required by the CLI application to
// display the results gathered by the "get" verb commands.
type GetDisplayer interface {
	WriteGetCluster(io.Writer, *cluster.Resource) error
	WriteGetNode(io.Writer, *node.Resource) error
	WriteGetNodeList(io.Writer, []*node.Resource) error
	WriteGetVolume(io.Writer, *volume.Resource) error
}

// NewCommand configures the set of commands which are grouped by the "get" verb.
func NewCommand(client GetClient, display GetDisplayer) *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "get retrieves a StorageOS resource, displaying basic information about it",
	}

	command.AddCommand(
		newCluster(os.Stdout, client, display),
		newNode(os.Stdout, client, display),
		newVolume(os.Stdout, client, display),
	)

	return command
}
