package update

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/output"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/pkg/version"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

var (
	errNoVolumeSpecified        = errors.New("must specify the volume to update")
	errTooManyArguments         = errors.New("must specify only one volume")
	errNoNamespaceSpecified     = errors.New("must specify the namespace of the volume to update")
	errNoFieldToUpdateSpecified = errors.New("must specify one field to update")
	errOnlyOneFieldAllowed      = errors.New("must specify only one field to change per update")
)

const (
	replicasName    = "replicas"
	descriptionName = "description"
)

type volumeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string

	// useCAS determines whether the command makes the set replicas request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string
	useAsync   bool

	fieldChanged string
	numReplicas  uint64
	description  string

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

	switch c.fieldChanged {
	case replicasName:
		params := &apiclient.SetReplicasRequestParams{}
		if c.useCAS() {
			params.CASVersion = version.FromString(c.casVersion)
		}
		err = c.client.SetReplicas(ctx, nsID, volID, c.numReplicas, params)
		if err != nil {
			return err
		}

		return c.display.SetReplicas(ctx, c.writer)

	case descriptionName:
		params := &apiclient.UpdateVolumeRequestParams{}
		if c.useCAS() {
			params.CASVersion = version.FromString(c.casVersion)
		}
		vol, err := c.client.UpdateVolumeDescription(ctx, nsID, volID, c.description, params)
		if err != nil {
			return err
		}

		return c.display.UpdateVolumeDescription(ctx, c.writer, output.NewVolumeUpdate(vol))
	}

	return errNoFieldToUpdateSpecified
}

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
$ storageos update volume --namespace my-namespace-name my-volume-name --replicas 3
$ storageos update volume --namespace my-namespace-name my-volume-name -r 3

$ storageos update volume --namespace my-namespace-name my-volume-name --description "new description"
$ storageos update volume --namespace my-namespace-name my-volume-name -d "new description"
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errNoVolumeSpecified
			}
			if len(args) > 1 {
				return errTooManyArguments
			}

			changed := 0
			if cmd.Flags().Changed(replicasName) {
				changed++
				c.fieldChanged = replicasName
			}
			if cmd.Flags().Changed(descriptionName) {
				changed++
				c.fieldChanged = descriptionName
			}

			if changed == 0 {
				return errNoFieldToUpdateSpecified
			}

			if changed > 1 {
				return errOnlyOneFieldAllowed
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
	cobraCommand.Flags().StringVarP(&c.description, descriptionName, "d", "", "set a new description")

	return cobraCommand
}
