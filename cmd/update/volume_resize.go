package update

import (
	"context"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/version"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

type volumeSizeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string
	volumeID  string
	sizeBytes uint64

	// useCAS determines whether the command makes the set replicas request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string
	useAsync   bool

	writer io.Writer
}

func (c *volumeSizeCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	var nsID id.Namespace

	if useIDs {
		nsID = id.Namespace(c.namespace)
	} else {
		ns, err := c.client.GetNamespaceByName(ctx, c.namespace)
		if err != nil {
			return err
		}
		nsID = ns.ID
	}

	var volID id.Volume

	if useIDs {
		volID = id.Volume(c.volumeID)
	} else {
		vol, err := c.client.GetVolumeByName(ctx, nsID, c.volumeID)
		if err != nil {
			return err
		}
		volID = vol.ID
	}

	params := &apiclient.ResizeVolumeRequestParams{}
	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}

	updatedVol, err := c.client.ResizeVolume(ctx, nsID, volID, c.sizeBytes, params)
	if err != nil {
		return err
	}

	return c.display.ResizeVolume(ctx, c.writer, output.NewVolumeUpdate(updatedVol))

}

func newVolumeSize(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeSizeCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "size [volume name] [size]",
		Short: "Updates a volume's size",
		Example: `
$ storageos update volume size my-volume-name 42GiB --namespace my-namespace-name
$ storageos update volume size my-volume-name 42gib --namespace my-namespace-name
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {

			if len(args) != 2 {
				return newErrInvalidArgNum(args, 2)
			}

			c.volumeID = args[0]

			bytes, err := humanize.ParseBytes(args[1])
			if err != nil {
				return newErrInvalidSizeArg(args[1])
			}

			c.sizeBytes = bytes

			return nil
		}),

		PreRunE: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, args []string) error {
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

	c.useCAS = flagutil.SupportCAS(cobraCommand.Flags(), &c.casVersion)
	flagutil.SupportAsync(cobraCommand.Flags(), &c.useAsync)

	return cobraCommand
}
