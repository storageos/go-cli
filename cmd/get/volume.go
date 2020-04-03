package get

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/selectors"
	"code.storageos.net/storageos/c2-cli/volume"
)

var (
	errNoNamespaceSpecified = errors.New("must specify a namespace to get volumes from")
	errMissingNamespace     = errors.New("volume contains a missing namespace reference")
)

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
		volumes, err := c.client.GetAllVolumes(ctx)
		if err != nil {
			return err
		}

		volumes = set.FilterVolumes(volumes)

		// Return early if there are no volumes to display
		if len(volumes) == 0 {
			return c.display.GetListVolumes(ctx, c.writer, nil)
		}

		namespaces, err := c.getNamespaces(ctx)
		if err != nil {
			return err
		}

		nodes, err := c.getNodeMapping(ctx)
		if err != nil {
			return err
		}

		outputVolumes := make([]*output.Volume, 0, len(volumes))
		for _, v := range volumes {

			ns, ok := namespaces[v.Namespace]
			if !ok {
				return errMissingNamespace
			}

			outputVolumes = append(outputVolumes, output.NewVolume(v, ns, nodes))
		}

		return c.display.GetListVolumes(
			ctx,
			c.writer,
			outputVolumes,
		)
	}

	namespaceID, err := c.getNamespaceID(ctx)
	if err != nil {
		return err
	}

	switch len(args) {
	case 1:
		// If only one volume is requested then API requests can be minimised.
		v, err := c.getVolume(ctx, namespaceID, args[0])
		if err != nil {
			return err
		}

		namespace, err := c.client.GetNamespace(ctx, v.Namespace)
		if err != nil {
			return err
		}

		nodes, err := c.getNodeMapping(ctx)
		if err != nil {
			return err
		}

		return c.display.GetVolume(ctx, c.writer, output.NewVolume(v, namespace, nodes))

	default:
		volumes, err := c.listVolumes(ctx, namespaceID, args)
		if err != nil {
			return err
		}

		volumes = set.FilterVolumes(volumes)

		// Return early if there are no volumes to display
		if len(volumes) == 0 {
			return c.display.GetListVolumes(ctx, c.writer, nil)
		}

		namespaces, err := c.getNamespaces(ctx)
		if err != nil {
			return err
		}

		nodes, err := c.getNodeMapping(ctx)
		if err != nil {
			return err
		}

		outputVolumes := make([]*output.Volume, 0, len(volumes))
		for _, v := range volumes {

			ns, ok := namespaces[v.Namespace]
			if !ok {
				return errMissingNamespace
			}

			outputVolumes = append(outputVolumes, output.NewVolume(v, ns, nodes))
		}

		return c.display.GetListVolumes(
			ctx,
			c.writer,
			outputVolumes,
		)
	}
}

// getNamespaceNames returns a map with namespace id as keys and namespace names
// (string) as value.
// List of volumes in input is used to filter out all unnecessary namespaces
func (c *volumeCommand) getNamespaces(ctx context.Context) (map[id.Namespace]*namespace.Resource, error) {
	namespaces, err := c.client.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	namespacesMap := make(map[id.Namespace]*namespace.Resource)
	for _, n := range namespaces {
		namespacesMap[n.ID] = n
	}

	return namespacesMap, nil
}

// getNamespaceID retrieves the namespace ID to use for API client calls based
// on the configuration of c.
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

	return c.client.GetNamespaceVolumesByUID(ctx, ns, volIDs...)
}

// getNodeMapping fetches the list of nodes from the API and builds a map from
// their ID to the full resource.
func (c *volumeCommand) getNodeMapping(ctx context.Context) (map[id.Node]*node.Resource, error) {
	nodeList, err := c.client.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	nodes := map[id.Node]*node.Resource{}
	for _, n := range nodeList {
		nodes[n.ID] = n
	}

	return nodes, nil
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
		Short:   "Retrieve basic details of volumes",
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
				runwrappers.EnsureNamespaceSetWhenUseIDs(c.config),
				runwrappers.EnsureTargetOrSelectors(&c.selectors),
				runwrappers.AuthenticateClient(c.config, c.client),
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
