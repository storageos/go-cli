package get

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/selectors"
	"code.storageos.net/storageos/c2-cli/volume"
)

var errNoNamespaceSpecified = errors.New("must specify a namespace to get volumes from")

type volumeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	allNamespaces bool
	namespace     string
	selectors     []string

	writer io.Writer
}

func (c *volumeCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	set, err := selectors.NewSetFromStrings(c.selectors...)
	if err != nil {
		return err
	}

	if c.allNamespaces {
		vols, err := c.client.GetAllVolumes(ctx)
		if err != nil {
			return err
		}

		return c.display.GetListVolumes(
			ctx,
			c.writer,
			set.FilterVolumes(vols),
		)
	}

	nsID, err := c.getNamespaceID(ctx)
	if err != nil {
		return err
	}

	switch len(args) {
	case 1:
		// If only one volume is requested then API requests can be minimised.
		v, err := c.getVolume(ctx, nsID, args[0])
		if err != nil {
			return err
		}

		return c.display.GetVolume(ctx, c.writer, v)

	default:
		vols, err := c.listVolumes(ctx, nsID, args)
		if err != nil {
			return err
		}

		return c.display.GetListVolumes(
			ctx,
			c.writer,
			set.FilterVolumes(vols),
		)
	}
}

// getNamespaceID retrieves the namespace ID to use for API client calls based
// on the configuration of c.
//
// If no namespace is specified on c, the ID of the default namespace is
// retrieved.
func (c *volumeCommand) getNamespaceID(ctx context.Context) (id.Namespace, error) {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return "", err
	}

	if !useIDs {
		ns, err := c.client.GetNamespaceByName(ctx, c.namespace)
		if err != nil {
			return "", err
		}

		return ns.ID, nil
	}
	return id.Namespace(c.namespace), nil
}

// getVolume requests a volume resource from ns, treating vol as an ID or a
// name based on c's configuration.
func (c *volumeCommand) getVolume(ctx context.Context, ns id.Namespace, vol string) (*volume.Resource, error) {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return nil, err
	}

	if !useIDs {
		return c.client.GetVolumeByName(ctx, ns, vol)
	}

	return c.client.GetVolume(ctx, ns, id.Volume(vol))
}

// listVolumes requests a list of volume resources using the configured API
// client, filtering using vols (if provided) as c's configuration dictates.
func (c *volumeCommand) listVolumes(ctx context.Context, ns id.Namespace, vols []string) ([]*volume.Resource, error) {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return nil, err
	}

	if !useIDs {
		return c.client.GetNamespaceVolumesByName(ctx, ns, vols...)
	}

	volIDs := []id.Volume{}
	for _, uid := range vols {
		volIDs = append(volIDs, id.Volume(uid))
	}

	return c.client.GetNamespaceVolumes(ctx, ns, volIDs...)
}

func newVolume(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Aliases: []string{"volumes"},
		Use:     "volume [volume names...]",
		Short:   "volume retrieves basic information about StorageOS volumes",
		Example: `
$ storageos get volumes --all-namespaces

$ storageos get volume --namespace my-namespace-name my-volume-name
`,
		Args: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, args []string) error {
			if c.allNamespaces && len(args) > 0 {
				return errors.New("volumes cannot be retrieved by name or ID across all namespaces")
			}
			return nil
		}),

		PreRunE: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, args []string) error {
			c.display = SelectDisplayer(c.config)

			ns, err := c.config.Namespace()
			if err != nil {
				return err
			}

			if ns == "" && !c.allNamespaces {
				return errNoNamespaceSpecified
			}
			c.namespace = ns

			return nil
		}),

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.EnsureTargetOrSelectors(&c.selectors),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		// If a legitimate error occurs as part of the VERB volume command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	cobraCommand.Flags().BoolVar(&c.allNamespaces, "all-namespaces", false, "retrieves volumes from all accessible namespaces. This option overrides the namespace configuration")

	flagutil.SupportSelectors(cobraCommand.Flags(), &c.selectors)

	return cobraCommand
}
