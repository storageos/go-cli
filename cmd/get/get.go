package get

import (
	"context"
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
	GetCluster(context.Context) (*cluster.Resource, error)

	GetNode(context.Context, id.Node) (*node.Resource, error)
	GetListNodes(context.Context, ...id.Node) ([]*node.Resource, error)

	GetVolume(context.Context, id.Namespace, id.Volume) (*volume.Resource, error)
	GetAllVolumes(context.Context) ([]*volume.Resource, error)
	GetNamespaceVolumes(context.Context, id.Namespace, ...id.Volume) ([]*volume.Resource, error)
}

// GetDisplayer defines the functionality required by the CLI application to
// display the results gathered by the "get" verb commands.
type GetDisplayer interface {
	GetCluster(io.Writer, *cluster.Resource) error
	GetNode(io.Writer, *node.Resource) error
	GetNodeList(io.Writer, []*node.Resource) error
	GetVolume(io.Writer, *volume.Resource) error
	GetVolumeList(io.Writer, []*volume.Resource) error
}

// NewCommand configures the set of commands which are grouped by the "get" verb.
func NewCommand(client GetClient) *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "get retrieves a StorageOS resource, displaying basic information about it",
	}

	command.AddCommand(
		newCluster(os.Stdout, client),
		newNode(os.Stdout, client),
		newVolume(os.Stdout, client),
	)

	return command
}
