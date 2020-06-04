package update

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/output/textformat"
	"code.storageos.net/storageos/c2-cli/output/yamlformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/volume"
)

var (
	errNoNamespaceSpecified = errors.New("no namespace specified")
)

// ConfigProvider specifies the configuration
type ConfigProvider interface {
	Username() (string, error)
	Password() (string, error)

	CommandTimeout() (time.Duration, error)
	UseIDs() (bool, error)
	OutputFormat() (output.Format, error)
	Namespace() (string, error)
}

// Client defines the functionality required by the CLI application to
// reasonably implement the "update" verb commands.
type Client interface {
	Authenticate(ctx context.Context, username, password string) (apiclient.AuthSession, error)

	SetReplicas(ctx context.Context, nsID id.Namespace, volID id.Volume, numReplicas uint64, params *apiclient.SetReplicasRequestParams) error
	UpdateVolumeDescription(ctx context.Context, nsID id.Namespace, volID id.Volume, description string, params *apiclient.UpdateVolumeRequestParams) (*volume.Resource, error)
	UpdateVolumeLabels(ctx context.Context, nsID id.Namespace, volID id.Volume, labels labels.Set, params *apiclient.UpdateVolumeRequestParams) (*volume.Resource, error)
	ResizeVolume(ctx context.Context, nsID id.Namespace, volID id.Volume, sizeBytes uint64, params *apiclient.ResizeVolumeOptionalRequestParams) (*volume.Resource, error)

	GetVolumeByName(ctx context.Context, namespace id.Namespace, name string) (*volume.Resource, error)
	GetVolume(ctx context.Context, namespaceID id.Namespace, uid id.Volume) (*volume.Resource, error)
	GetNamespaceByName(ctx context.Context, name string) (*namespace.Resource, error)
}

// Displayer defines the functionality required by the CLI application to
// display the results returned by "update" verb operations.
type Displayer interface {
	UpdateVolume(ctx context.Context, w io.Writer, updatedVol output.VolumeUpdate) error
	SetReplicas(ctx context.Context, w io.Writer, new uint64) error
}

// NewCommand configures the set of commands which are grouped by the "update" verb.
func NewCommand(client Client, config ConfigProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "update",
		Short: "Make changes to existing resources",
	}

	command.AddCommand(
		newVolumeUpdate(client, config),
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
