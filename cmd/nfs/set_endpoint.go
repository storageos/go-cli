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

type volumeNFSMountEndpointCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string
	volumeID  string
	endpoint  string

	// useCAS determines whether the command makes the set replicas request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string
	useAsync   bool

	writer io.Writer
}

func (c *volumeNFSMountEndpointCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
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

	params := &apiclient.UpdateNFSVolumeMountEndpointRequestParams{}
	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}

	err = c.client.UpdateNFSVolumeMountEndpoint(ctx, nsID, volID, c.endpoint, params)
	if err != nil {
		return err
	}

	// Display the "request submitted" message if it was async, instead of
	// the deletion confirmation below.
	if c.useAsync {
		return c.display.AsyncRequest(ctx, c.writer)
	}

	return c.display.UpdateNFSVolumeMountEndpoint(ctx, c.writer, volID, c.endpoint)
}

func newSetEndpoint(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeNFSMountEndpointCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "endpoint [endpoint]",
		Short: "Updates a volume's NFS mount endpoint",
		Example: `
$ storageos nfs endpoint my-volume-name "10.0.0.1:/" --namespace my-namespace-name
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {

			if len(args) != 2 {
				return clierr.NewErrInvalidArgNum(args, 2, "storageos update volume nfs-endpoint [volume name] [endpoint]")
			}

			c.volumeID = args[0]

			c.endpoint = args[1]

			return nil
		}),

		PreRunE: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, args []string) error {
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
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	c.useCAS = flagutil.SupportCAS(cobraCommand.Flags(), &c.casVersion)
	flagutil.SupportAsync(cobraCommand.Flags(), &c.useAsync)

	return cobraCommand
}
