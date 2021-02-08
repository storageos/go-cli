package get

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/diagnostics"
	"code.storageos.net/storageos/c2-cli/licence"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/output/textformat"
	"code.storageos.net/storageos/c2-cli/output/yamlformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/policygroup"
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
	Namespace() (string, error)
	OutputFormat() (output.Format, error)
}

// Client defines the functionality required by the CLI application to
// reasonably implement the "get" verb commands.
type Client interface {
	Authenticate(ctx context.Context, username, password string) (apiclient.AuthSession, error)

	GetCluster(ctx context.Context) (*cluster.Resource, error)
	GetLicence(ctx context.Context) (*licence.Resource, error)
	GetDiagnostics(ctx context.Context) (*diagnostics.BundleReadCloser, error)
	GetSingleNodeDiagnostics(ctx context.Context, id id.Node) (*diagnostics.BundleReadCloser, error)
	GetSingleNodeDiagnosticsByName(ctx context.Context, name string) (*diagnostics.BundleReadCloser, error)

	GetUser(ctx context.Context, userID id.User) (*user.Resource, error)
	GetUserByName(ctx context.Context, username string) (*user.Resource, error)
	ListUsers(ctx context.Context) ([]*user.Resource, error)
	GetListUsersByUID(ctx context.Context, uIDs []id.User) ([]*user.Resource, error)
	GetListUsersByUsername(ctx context.Context, usernames []string) ([]*user.Resource, error)

	GetPolicyGroup(ctx context.Context, pgID id.PolicyGroup) (*policygroup.Resource, error)
	GetPolicyGroupByName(ctx context.Context, name string) (*policygroup.Resource, error)
	GetListPolicyGroupsByName(ctx context.Context, names ...string) ([]*policygroup.Resource, error)
	GetListPolicyGroupsByUID(ctx context.Context, gids ...id.PolicyGroup) ([]*policygroup.Resource, error)

	ListNodes(ctx context.Context) ([]*node.Resource, error)
	GetNode(ctx context.Context, uid id.Node) (*node.Resource, error)
	GetNodeByName(ctx context.Context, name string) (*node.Resource, error)
	GetListNodesByUID(ctx context.Context, uids ...id.Node) ([]*node.Resource, error)
	GetListNodesByName(ctx context.Context, names ...string) ([]*node.Resource, error)

	GetVolume(ctx context.Context, namespaceID id.Namespace, uid id.Volume) (*volume.Resource, error)
	GetVolumeByName(ctx context.Context, namespaceID id.Namespace, name string) (*volume.Resource, error)
	GetAllVolumes(ctx context.Context) ([]*volume.Resource, error)
	GetNamespaceVolumesByUID(ctx context.Context, namespaceID id.Namespace, uids ...id.Volume) ([]*volume.Resource, error)
	GetNamespaceVolumesByName(ctx context.Context, namespaceID id.Namespace, names ...string) ([]*volume.Resource, error)

	GetNamespace(ctx context.Context, uid id.Namespace) (*namespace.Resource, error)
	GetNamespaceByName(ctx context.Context, name string) (*namespace.Resource, error)
	GetListNamespacesByUID(ctx context.Context, uids ...id.Namespace) ([]*namespace.Resource, error)
	GetListNamespacesByName(ctx context.Context, name ...string) ([]*namespace.Resource, error)
	ListNamespaces(ctx context.Context) ([]*namespace.Resource, error)
}

// Displayer defines the functionality required by the CLI application to
// display the results gathered by the "get" verb commands.
type Displayer interface {
	GetCluster(ctx context.Context, w io.Writer, cluster *output.Cluster) error
	GetLicence(ctx context.Context, w io.Writer, l *output.Licence) error
	GetUser(ctx context.Context, w io.Writer, user *output.User) error
	GetUsers(ctx context.Context, w io.Writer, users []*output.User) error
	GetDiagnostics(ctx context.Context, w io.Writer, outputPath string) error
	GetNode(ctx context.Context, w io.Writer, node *output.Node) error
	GetListNodes(ctx context.Context, w io.Writer, nodes []*output.Node) error
	GetNamespace(ctx context.Context, w io.Writer, namespace *output.Namespace) error
	GetPolicyGroup(ctx context.Context, w io.Writer, group *output.PolicyGroup) error
	GetListPolicyGroups(ctx context.Context, w io.Writer, groups []*output.PolicyGroup) error
	GetListNamespaces(ctx context.Context, w io.Writer, namespaces []*output.Namespace) error
	GetVolume(ctx context.Context, w io.Writer, volume *output.Volume) error
	GetListVolumes(ctx context.Context, w io.Writer, volumes []*output.Volume) error
}

// NewCommand configures the set of commands which are grouped by the "get" verb.
func NewCommand(client Client, config ConfigProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "Fetch basic details for resources",
	}

	command.AddCommand(
		newCluster(os.Stdout, client, config),
		newDiagnostics(os.Stdout, client, config),
		newNode(os.Stdout, client, config),
		newNamespace(os.Stdout, client, config),
		newVolume(os.Stdout, client, config),
		newUser(os.Stdout, client, config),
		newPolicyGroup(os.Stdout, client, config),
		newLicence(os.Stdout, client, config),
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
