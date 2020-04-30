package create

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/licence"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/output/textformat"
	"code.storageos.net/storageos/c2-cli/output/yamlformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
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

// Client describes the functionality required by the CLI application
// to reasonably implement the "create" verb commands.
type Client interface {
	Authenticate(ctx context.Context, username, password string) (apiclient.AuthSession, error)

	CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error)
	CreateVolume(ctx context.Context, namespace id.Namespace, name, description string, fs volume.FsType, sizeBytes uint64, labelSet labels.Set, params *apiclient.CreateVolumeRequestParams) (*volume.Resource, error)
	CreateNamespace(ctx context.Context, name string, labelSet labels.Set) (*namespace.Resource, error)
	CreatePolicyGroup(ctx context.Context, name string, specs []*policygroup.Spec) (*policygroup.Resource, error)

	GetLicence(ctx context.Context) (*licence.Resource, error)
	GetNamespace(ctx context.Context, uid id.Namespace) (*namespace.Resource, error)
	GetNamespaceByName(ctx context.Context, name string) (*namespace.Resource, error)
	GetListNodesByUID(ctx context.Context, uids ...id.Node) ([]*node.Resource, error)
	GetListPolicyGroupsByUID(ctx context.Context, gids ...id.PolicyGroup) ([]*policygroup.Resource, error)
	GetListPolicyGroupsByName(ctx context.Context, names ...string) ([]*policygroup.Resource, error)

	ListNamespaces(ctx context.Context) ([]*namespace.Resource, error)
}

// Displayer describes the functionality required by the CLI application
// to display the resources produced by the "create" verb commands.
type Displayer interface {
	CreateUser(ctx context.Context, w io.Writer, user *output.User) error
	CreateVolume(ctx context.Context, w io.Writer, volume *output.Volume) error
	CreateVolumeAsync(ctx context.Context, w io.Writer) error
	CreateNamespace(ctx context.Context, w io.Writer, namespace *output.Namespace) error
	CreatePolicyGroup(ctx context.Context, w io.Writer, group *output.PolicyGroup) error
}

// NewCommand configures the set of commands which are grouped by the "create"
// verb.
func NewCommand(client Client, config ConfigProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "Create new resources",
	}

	command.AddCommand(
		newUser(os.Stdout, client, config),
		newVolume(os.Stdout, client, config),
		newNamespace(os.Stdout, client, config),
		newPolicyGroup(os.Stdout, client, config),
	)

	return command
}

// SelectDisplayer instantiates the appropriate Displayer for the settings
// given by config.
func SelectDisplayer(config ConfigProvider) Displayer {
	out, err := config.OutputFormat()
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
