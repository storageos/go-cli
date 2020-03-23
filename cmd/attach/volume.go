package attach

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

var (
	errArguments            = errors.New("must specify exactly two arguments <volume> <node>")
	errNoNamespaceSpecified = errors.New("must specify the namespace of the volume to attach")
)

type volumeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string

	writer io.Writer
}

func (c *volumeCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
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
		node, err := c.client.GetNodeByName(ctx, nodeName)
		if err != nil {
			return err
		}
		nodeID = node.ID
	}

	err = c.client.AttachVolume(ctx, namespaceID, volumeID, nodeID)
	if err != nil {
		return err
	}

	return c.display.AttachVolume(ctx, c.writer)
}

func newVolume(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "volume",
		Short: "Attach a volume to a node",
		Example: `
$ storageos attach volume --namespace my-namespace-name my-volume my-node
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
			)(c.runWithCtx)

			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	return cobraCommand
}
