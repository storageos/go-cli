package detach

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/output/textformat"
	"code.storageos.net/storageos/c2-cli/output/yamlformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/volume"
)

var (
	errArguments            = errors.New("must specify exactly one volume to detach")
	errNoNamespaceSpecified = errors.New("must specify the namespace of the volume to detach")
)

// ConfigProvider specifies the configuration settings which commands require
// access to.
type ConfigProvider interface {
	CommandTimeout() (time.Duration, error)
	UseIDs() (bool, error)
	Namespace() (string, error)
	OutputFormat() (output.Format, error)
}

// Client describes the functionality required by the CLI application
// to reasonably implement the "detach" verb commands.
type Client interface {
	GetNamespaceByName(ctx context.Context, name string) (*namespace.Resource, error)
	GetVolumeByName(ctx context.Context, namespaceID id.Namespace, name string) (*volume.Resource, error)
	DetachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *apiclient.DetachVolumeRequestParams) error
}

// Displayer describes the functionality required by the CLI application
// to display the resources produced by the "detach" verb commands.
type Displayer interface {
	DetachVolume(ctx context.Context, w io.Writer) error
}

type detachCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string

	// useCAS determines whether the command makes the detach request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string

	writer io.Writer
}

func (c *detachCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	namespaceID := id.Namespace(c.namespace)
	volumeID := id.Volume(args[0])

	// If not using IDs, get the appropriate IDs
	if !useIDs {
		ns, err := c.client.GetNamespaceByName(ctx, c.namespace)
		if err != nil {
			return err
		}
		namespaceID = ns.ID

		volName := args[0]
		vol, err := c.client.GetVolumeByName(ctx, namespaceID, volName)
		if err != nil {
			return err
		}
		volumeID = vol.ID
	}

	params := &apiclient.DetachVolumeRequestParams{}
	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}

	err = c.client.DetachVolume(
		ctx,
		namespaceID,
		volumeID,
		params,
	)
	if err != nil {
		return err
	}

	return c.display.DetachVolume(ctx, c.writer)
}

// NewCommand configures the "detach" verb command.
func NewCommand(client Client, config ConfigProvider) *cobra.Command {
	c := &detachCommand{
		config: config,
		client: client,
		writer: os.Stdout,
	}

	cobraCommand := &cobra.Command{
		Use:   "detach",
		Short: "Detach a volume from its current location",
		Example: `
$ storageos detach volume --namespace my-namespace-name my-volume
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errArguments
			}
			return nil
		}),
		PreRunE: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, _ []string) error {
			c.display = SelectDisplayer(c.config)

			ns, err := c.config.Namespace()
			if err != nil {
				return err
			}

			if ns == "" {
				return errNoNamespaceSpecified
			}
			c.namespace = ns

			return nil
		}),

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.EnsureNamespaceSetWhenUseIDs(c.config),
			)(c.runWithCtx)

			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	c.useCAS = flagutil.SupportCAS(cobraCommand.Flags(), &c.casVersion)

	return cobraCommand
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
		return jsonformat.NewDisplayer("")
	case output.YAML:
		return yamlformat.NewDisplayer("")
	case output.Text:
		fallthrough
	default:
		return textformat.NewDisplayer(textformat.NewTimeFormatter())
	}
}
