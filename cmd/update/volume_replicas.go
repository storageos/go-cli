package update

import (
	"context"
	"io"
	"strconv"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/clierr"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

type volumeReplicasCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace      string
	volumeID       string
	targetReplicas uint64

	// useCAS determines whether the command makes the set replicas request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string
	useAsync   bool

	writer io.Writer
}

func (c *volumeReplicasCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
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

	params := &apiclient.SetReplicasRequestParams{}
	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}

	err = c.client.SetReplicas(ctx, nsID, volID, c.targetReplicas, params)
	if err != nil {
		return err
	}

	return c.display.SetReplicas(ctx, c.writer, c.targetReplicas)

}

func newVolumeReplicas(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeReplicasCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "replicas [volume name] [target number]",
		Short: "Updates a volume's target replica number",
		Example: `
$ storageos update volume replicas my-volume-name 3 --namespace my-namespace-name
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {

			if len(args) != 2 {
				return clierr.NewErrInvalidArgNum(args, 2, "storageos update volume replicas [volume] [replica num]")
			}

			c.volumeID = args[0]

			targetReplicas, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			if targetReplicas > 5 {
				return newErrInvalidReplicaNum(targetReplicas)
			}

			c.targetReplicas = targetReplicas

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
	flagutil.WarnAboutValueBeingOverwrittenByK8s(cobraCommand.Flags())

	return cobraCommand
}
