package update

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/clierr"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

type volumeLabelsCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string
	volumeID  string
	labels    labels.Set

	// flag values
	upsertStr    string
	upsertLabels labels.Set
	deleteStr    string
	deleteLabels labels.Set

	// useCAS determines whether the command makes the set replicas request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string
	useAsync   bool

	writer io.Writer
}

func (c *volumeLabelsCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {

	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	// fetch ns

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

	// fetch volume and upsert/delete labels as desired

	var volID id.Volume

	if useIDs {
		volID = id.Volume(c.volumeID)

		vol, err := c.client.GetVolume(ctx, nsID, volID)
		if err != nil {
			return err
		}

		c.labels = labels.UpsertExisting(vol.Labels, c.upsertLabels)
		c.labels = labels.RemoveExisting(vol.Labels, c.deleteLabels)

	} else {
		vol, err := c.client.GetVolumeByName(ctx, nsID, c.volumeID)
		if err != nil {
			return err
		}

		volID = vol.ID

		c.labels = labels.UpsertExisting(vol.Labels, c.upsertLabels)
		c.labels = labels.RemoveExisting(vol.Labels, c.deleteLabels)
	}

	// perform the update

	params := &apiclient.UpdateVolumeRequestParams{}

	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}

	// If asynchrony is specified then source the timeout and initialise the params.
	if c.useAsync {
		timeout, err := c.config.CommandTimeout()
		if err != nil {
			return err
		}

		params.AsyncMax = timeout
	}

	updatedVol, err := c.client.UpdateVolumeLabels(ctx, nsID, volID, c.labels, params)
	if err != nil {
		return err
	}

	if c.useAsync {
		return c.display.AsyncRequest(ctx, c.writer)
	}

	return c.display.UpdateVolumeLabels(ctx, c.writer, output.NewVolumeUpdate(updatedVol))
}

func newVolumeLabels(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeLabelsCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "labels [volume name] [labels]",
		Short: "Updates a volume's labels",
		Example: `
$ storageos update volume labels my-volume-name app=myapp,tier=production --namespace my-namespace-name
$ storageos update volume labels my-volume-name --delete app=myapp --upsert tier=production --namespace my-namespace-name
$ storageos update volume labels my-volume-name --upsert tier=production --namespace my-namespace-name
$ storageos update volume labels my-volume-name --delete app=myapp --namespace my-namespace-name
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {

			var err error
			// parse upsert/delete strings if any
			c.upsertLabels, err = labels.FromString(c.upsertStr)
			if err != nil {
				return err
			}

			c.deleteLabels, err = labels.FromString(c.deleteStr)
			if err != nil {
				return err
			}

			if len(c.deleteLabels) != 0 {
				if len(args) != 1 {
					return clierr.NewErrInvalidArgNum(args, 1, "storageos update volume labels --delete a=1 [volume]")
				}

				c.volumeID = args[0]
				// the user used the --upsert or --delete flag
				return nil
			}

			if len(c.upsertLabels) != 0 {
				if len(args) != 1 {
					return clierr.NewErrInvalidArgNum(args, 1, "storageos update volume labels --upsert a=2 [volume]")
				}

				c.volumeID = args[0]
				// the user used the --upsert or --delete flag
				return nil
			}

			// the user used the "set" behavior, they cannot use both this and the
			// flags

			if len(args) != 2 {
				return clierr.NewErrInvalidArgNum(args, 2, "storageos update volume labels [volume] [labels]")
			}

			c.volumeID = args[0]

			labels, err := labels.FromString(args[1])
			if err != nil {
				return err
			}

			c.labels = labels

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

	// local flags
	cobraCommand.Flags().StringVar(&c.upsertStr, "upsert", "", "upsert labels to the volume's existing set")
	cobraCommand.Flags().StringVar(&c.deleteStr, "delete", "", "delete labels from the volume's existing set")

	return cobraCommand
}
