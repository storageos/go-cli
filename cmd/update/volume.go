package update

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/pkg/version"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

var (
	errNoVolumeSpecified      = errors.New("must specify the volume to update")
	errTooManyArguments       = errors.New("must specify only one volume")
	errNoNamespaceSpecified   = errors.New("must specify the namespace of the volume to update")
	errNoNumReplicasSpecified = errors.New("must specify the number of replicas requested")
)

type volumeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string

	// useCAS determines whether the command makes the set replicas request
	// constrained by the provided casVersion.
	useCAS      func() bool
	casVersion  string
	useAsync    bool
	numReplicas uint64

	writer io.Writer
}

func (c *volumeCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
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
		volID = id.Volume(args[0])
	} else {
		vol, err := c.client.GetVolumeByName(ctx, nsID, args[0])
		if err != nil {
			return err
		}
		volID = vol.ID
	}

	params := &apiclient.SetReplicasRequestParams{}
	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}
	err = c.client.SetReplicas(ctx, nsID, volID, c.numReplicas, params)
	if err != nil {
		return err
	}

	return c.display.SetReplicas(ctx, c.writer)
}

const replicasName = "replicas"

func newVolume(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "volume",
		Short: "Update a volume",
		Example: `
$ storageos update volume --namespace my-namespace-name my-volume-name --replicas 3`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errNoVolumeSpecified
			}
			if len(args) > 1 {
				return errTooManyArguments
			}

			if !cmd.Flags().Changed(replicasName) {
				return errNoNumReplicasSpecified
			}

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
	cobraCommand.Flags().Uint64VarP(&c.numReplicas, replicasName, "r", 0, "set number of replicas requested")

	return cobraCommand
}
