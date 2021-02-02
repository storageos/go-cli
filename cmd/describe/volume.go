package describe

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/clierr"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/selectors"
	"code.storageos.net/storageos/c2-cli/volume"
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
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

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

		// If no volumes match the filter then return
		if len(volumes) == 0 {
			return c.display.DescribeVolume(ctx, c.writer, nil)
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
				return clierr.ErrNoNamespaceSpecified
			}

			outputVolumes = append(outputVolumes, output.NewVolume(v, ns, nodes))
		}

		return c.display.DescribeListVolumes(
			ctx,
			c.writer,
			outputVolumes,
		)
	}

	var ns *namespace.Resource

	if useIDs {
		ns, err = c.client.GetNamespace(ctx, id.Namespace(c.namespace))
		if err != nil {
			return err
		}
	} else {
		ns, err = c.client.GetNamespaceByName(ctx, c.namespace)
		if err != nil {
			return err
		}
	}

	switch len(args) {
	case 1:
		var vol *volume.Resource
		var err error

		if useIDs {
			vol, err = c.client.GetVolume(ctx, ns.ID, id.Volume(args[0]))
		} else {
			vol, err = c.client.GetVolumeByName(ctx, ns.ID, args[0])
		}
		if err != nil {
			return err
		}

		nodes, err := c.getNodeMapping(ctx)
		if err != nil {
			return err
		}

		return c.display.DescribeVolume(ctx, c.writer, output.NewVolume(vol, ns, nodes))

	default:
		volumes, err := c.listVolumes(ctx, ns.ID, args)
		if err != nil {
			return err
		}

		volumes = set.FilterVolumes(volumes)

		nodes, err := c.getNodeMapping(ctx)
		if err != nil {
			return err
		}

		outputVols := make([]*output.Volume, 0, len(volumes))

		for _, vol := range volumes {
			outputVols = append(outputVols, output.NewVolume(vol, ns, nodes))
		}

		return c.display.DescribeListVolumes(ctx, c.writer, outputVols)
	}
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

func newVolume(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Aliases: []string{"volumes"},
		Use:     "volume [volume names...]",
		Short:   "Show detailed information for volumes",
		Example: `
$ storageos describe volumes

$ storageos describe volume --namespace my-namespace-name my-volume-name
`,

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
				runwrappers.EnsureTargetOrSelectors(&c.selectors),
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		// If a legitimate error occurs as part of the VERB volume command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	cobraCommand.Flags().BoolVarP(&c.allNamespaces, "all-namespaces", "A", false, "retrieves volumes from all accessible namespaces. This option overrides the namespace configuration")

	flagutil.SupportSelectors(cobraCommand.Flags(), &c.selectors)

	return cobraCommand
}
