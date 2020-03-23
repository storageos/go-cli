package detach

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

var (
	errArguments            = errors.New("must specify exactly one volume to detach")
	errNoNamespaceSpecified = errors.New("must specify the namespace of the volume to detach")
)

type volumeCommand struct {
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

func (c *volumeCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
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

func newVolume(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "volume",
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
