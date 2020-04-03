package describe

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/output/textformat"
	"code.storageos.net/storageos/c2-cli/output/yamlformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/c2-cli/volume"
)

// ConfigProvider specifies the configuration settings which commands require
// access to.
type ConfigProvider interface {
	Username() (string, error)
	Password() (string, error)

	CommandTimeout() (time.Duration, error)
	UseIDs() (bool, error)
	OutputFormat() (output.Format, error)
	Namespace() (string, error)
}

// Client describes the functionality required by the CLI application
// to reasonably implement the "describe" verb commands.
type Client interface {
	Authenticate(ctx context.Context, username, password string) (*user.Resource, error)

	GetCluster(ctx context.Context) (*cluster.Resource, error)

	GetNode(ctx context.Context, uid id.Node) (*node.Resource, error)
	GetNodeByName(ctx context.Context, name string) (*node.Resource, error)
	ListNodes(ctx context.Context) ([]*node.Resource, error)
	GetListNodesByUID(ctx context.Context, uids ...id.Node) ([]*node.Resource, error)
	GetListNodesByName(ctx context.Context, names ...string) ([]*node.Resource, error)

	GetVolume(ctx context.Context, namespace id.Namespace, vid id.Volume) (*volume.Resource, error)
	GetVolumeByName(ctx context.Context, namespace id.Namespace, name string) (*volume.Resource, error)
	GetNamespaceVolumesByUID(ctx context.Context, namespaceID id.Namespace, volIDs ...id.Volume) ([]*volume.Resource, error)
	GetNamespaceVolumesByName(ctx context.Context, namespaceID id.Namespace, names ...string) ([]*volume.Resource, error)
	GetAllVolumes(ctx context.Context) ([]*volume.Resource, error)

	GetNamespace(ctx context.Context, namespaceID id.Namespace) (*namespace.Resource, error)
	GetNamespaceByName(ctx context.Context, name string) (*namespace.Resource, error)
	ListNamespaces(ctx context.Context) ([]*namespace.Resource, error)
}

// Displayer defines the functionality required by the CLI application
// to display the results gathered by the "describe" verb commands.
type Displayer interface {
	DescribeCluster(ctx context.Context, w io.Writer, c *output.Cluster) error
	DescribeNode(ctx context.Context, w io.Writer, node *output.NodeDescription) error
	DescribeListNodes(ctx context.Context, w io.Writer, nodes []*output.NodeDescription) error
	DescribeVolume(ctx context.Context, w io.Writer, volume *output.Volume) error
	DescribeListVolumes(ctx context.Context, w io.Writer, volumes []*output.Volume) error
}

// NewCommand configures the set of commands which are grouped by the "describe" verb.
func NewCommand(client Client, config ConfigProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "describe",
		Short: "Fetch extended details for resources",
	}

	command.AddCommand(
		newNode(os.Stdout, client, config),
		newVolume(os.Stdout, client, config),
		newCluster(os.Stdout, client, config),
	)

	return command
}

// SelectDisplayer returns the right command displayer specified in the
// config provider.
func SelectDisplayer(cp ConfigProvider) Displayer {
	out, err := cp.OutputFormat()
	if err != nil {
		return textformat.NewDisplayer(textformat.NewTimeFormatter())
	}

	switch out {
	case output.JSON:
		return jsonformat.NewDisplayer(jsonformat.DefaultEncodingIndent)
	case output.YAML:
		return yamlformat.NewDisplayer("")
	case output.Text:
		fallthrough
	default:
		return textformat.NewDisplayer(textformat.NewTimeFormatter())
	}
}
