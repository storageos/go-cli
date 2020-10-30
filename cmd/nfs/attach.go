package nfs

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/clierr"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

type attachCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string

	// useCAS determines whether the command makes the set replicas request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string
	useAsync   bool

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
	)

	if useIDs {
		namespaceID = id.Namespace(c.namespace)
		volumeID = id.Volume(args[0])
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
	}

	params := &apiclient.AttachNFSVolumeRequestParams{}

	// If asynchrony is specified then source the timeout and initialise the
	// params.
	if c.useAsync {
		timeout, err := c.config.CommandTimeout()
		if err != nil {
			return err
		}

		params.AsyncMax = timeout
	}
	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}

	err = c.client.AttachNFSVolume(ctx, namespaceID, volumeID, params)

	if err != nil {
		return err
	}

	// Display the "request submitted" message if it was async, instead of
	// the deletion confirmation below.
	if c.useAsync {
		return c.display.AsyncRequest(ctx, c.writer)
	}

	return c.display.AttachVolume(ctx, c.writer)
}

// newAttach configures the "attach" subcommand for nfs.
func newAttach(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &attachCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "attach [volume]",
		Short: "Attach a NFS volume",
		Example: `
$ storageos nfs attach --namespace my-namespace-name my-volume 
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return clierr.NewErrInvalidArgNum(args, 1, "storageos nfs attach [volume]")
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
				return clierr.ErrNoNamespaceSpecified
			}
			c.namespace = ns

			return nil
		}),

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.EnsureNamespaceSetWhenUseIDs(c.config),
				runwrappers.AuthenticateClient(c.config, c.client),
				runwrappers.HandleLicenceError(c.client),
			)(c.runWithCtx)

			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	c.useCAS = flagutil.SupportCAS(cobraCommand.Flags(), &c.casVersion)
	flagutil.SupportAsync(cobraCommand.Flags(), &c.useAsync)

	return cobraCommand
}
