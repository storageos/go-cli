package get

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
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
	GetCluster(ctx context.Context) (*cluster.Resource, error)

	GetNode(ctx context.Context, uid id.Node) (*node.Resource, error)
	GetNodeByName(ctx context.Context, name string) (*node.Resource, error)
	GetListNodes(ctx context.Context, uids ...id.Node) ([]*node.Resource, error)
	GetListNodesByName(ctx context.Context, names ...string) ([]*node.Resource, error)

	GetVolume(ctx context.Context, namespaceID id.Namespace, uid id.Volume) (*volume.Resource, error)
	GetVolumeByName(ctx context.Context, namespaceID id.Namespace, name string) (*volume.Resource, error)
	GetAllVolumes(ctx context.Context) ([]*volume.Resource, error)
	GetNamespaceVolumes(ctx context.Context, namespaceID id.Namespace, uids ...id.Volume) ([]*volume.Resource, error)
	GetNamespaceVolumesByName(ctx context.Context, namespaceID id.Namespace, names ...string) ([]*volume.Resource, error)

	GetNamespaceByName(ctx context.Context, name string) (*namespace.Resource, error)
	GetAllNamespaces(ctx context.Context) ([]*namespace.Resource, error)
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
func NewCommand(initClient func() (*apiclient.Client, error), config ConfigProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "get retrieves a StorageOS resource, displaying basic information about it",
	}

	command.AddCommand(
		newCluster(os.Stdout, initClient, config),
		newNode(os.Stdout, initClient, config),
		newVolume(os.Stdout, initClient, config),
	)

	return command
}
