package get

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
)

// ConfigProvider specifies the configuration settings which commands require
// access to.
type ConfigProvider interface {
	CommandTimeout() (time.Duration, error)
}

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
	GetCluster(context.Context, io.Writer, *cluster.Resource) error
	GetNode(context.Context, io.Writer, *node.Resource) error
	GetNodeList(context.Context, io.Writer, []*node.Resource) error
	GetVolume(context.Context, io.Writer, *volume.Resource) error
	GetVolumeList(context.Context, io.Writer, []*volume.Resource) error
}

// NewCommand configures the set of commands which are grouped by the "get" verb.
func NewCommand(client GetClient, config ConfigProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "get retrieves a StorageOS resource, displaying basic information about it",
	}

	command.AddCommand(
		newCluster(os.Stdout, client, config),
		newNode(os.Stdout, client, config),
		newVolume(os.Stdout, client, config),
	)

	return command
}
