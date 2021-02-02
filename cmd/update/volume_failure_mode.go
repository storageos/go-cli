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
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

type volumeFailureModeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string
	volumeID  string

	intent       string
	threshold    uint64
	useThreshold bool

	// useCAS determines whether the command makes the set replicas request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string

	writer io.Writer
}

func (c *volumeFailureModeCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
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

	params := &apiclient.SetFailureModeRequestParams{}
	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}

	if c.useThreshold {
		updatedVol, err := c.client.SetFailureThreshold(ctx, nsID, volID, c.threshold, params)
		if err != nil {
			return err
		}

		return c.display.SetFailureMode(ctx, c.writer, output.NewVolumeUpdate(updatedVol))
	}

	updatedVol, err := c.client.SetFailureModeIntent(ctx, nsID, volID, c.intent, params)
	if err != nil {
		return err
	}

	return c.display.SetFailureMode(ctx, c.writer, output.NewVolumeUpdate(updatedVol))
}

func newVolumeFailureMode(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeFailureModeCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "failure-mode [volume name] [intent|threshold]",
		Short: "Updates a volume's failure mode behaviour",
		Example: `
$ storageos update volume failure-mode my-volume-name hard --namespace my-namespace-name
$ storageos update volume failure-mode my-volume-name 2 --namespace my-namespace-name
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {

			if len(args) != 2 {
				return clierr.NewErrInvalidArgNum(args, 2, "storageos update volume failure-mode [volume] [intent|threshold]")
			}

			c.volumeID = args[0]

			// Optimistically parse a failure threshold. If this fails then
			// the request should treat
			failureThreshold, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				c.intent = args[1]
			} else {
				c.threshold = failureThreshold
				c.useThreshold = true
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

	return cobraCommand
}
