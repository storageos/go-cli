package attach

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/output/textformat"
	"code.storageos.net/storageos/c2-cli/output/yamlformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
)

var (
	errArguments            = errors.New("must specify exactly two arguments <volume> <node>")
	errNoNamespaceSpecified = errors.New("must specify the namespace of the volume to attach")
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
// to reasonably implement the "attach" verb commands.
type Client interface {
	Authenticate(ctx context.Context, username, password string) (apiclient.AuthSession, error)

	GetNamespaceByName(ctx context.Context, name string) (*namespace.Resource, error)
	GetVolumeByName(ctx context.Context, namespaceID id.Namespace, name string) (*volume.Resource, error)
	GetNodeByName(ctx context.Context, name string) (*node.Resource, error)
	AttachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, nodeID id.Node) error
}

// Displayer describes the functionality required by the CLI application
// to display the resources produced by the "attach" verb commands.
type Displayer interface {
	AttachVolume(ctx context.Context, w io.Writer) error
}

type attachCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string

	writer io.Writer
}

func (c *attachCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	var (
		namespaceID id.Namespace
		volumeID    id.Volume
		nodeID      id.Node
	)

	if useIDs {
		namespaceID = id.Namespace(c.namespace)
		volumeID = id.Volume(args[0])
		nodeID = id.Node(args[1])
	} else {
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

		nodeName := args[1]
		n, err := c.client.GetNodeByName(ctx, nodeName)
		if err != nil {
			return err
		}
		nodeID = n.ID
	}

	err = c.client.AttachVolume(ctx, namespaceID, volumeID, nodeID)
	if err != nil {
		return err
	}

	return c.display.AttachVolume(ctx, c.writer)
}

// NewCommand configures the "attach" verb command.
func NewCommand(client Client, config ConfigProvider) *cobra.Command {
	c := &attachCommand{
		config: config,
		client: client,
		writer: os.Stdout,
	}

	cobraCommand := &cobra.Command{
		Use:   "attach",
		Short: "Attach a volume to a node",
		Example: `
$ storageos attach --namespace my-namespace-name my-volume my-node
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
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
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)

			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

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
