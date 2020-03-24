package describe

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
)

var (
	errNoNamespaceSpecified = errors.New("must specify a namespace to get a volume from")
	errMissingVolume        = errors.New("must specify a volume to describe")
)

type volumeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string
	selectors []string

	writer io.Writer
}

func (c *volumeCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
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

	volumes := make([]*volume.Resource, 0)
	var vol *volume.Resource
	for _, v := range args {
		if useIDs {
			vol, err = c.client.GetVolume(ctx, ns.ID, id.Volume(v))
			if err != nil {
				return err
			}
		} else {
			vol, err = c.client.GetVolumeByName(ctx, ns.ID, v)
			if err != nil {
				return err
			}
		}

		volumes = append(volumes, vol)
	}

	nodes, err := c.getNodeMapping(ctx)
	if err != nil {
		return err
	}

	outputVols := make([]*output.Volume, 0, len(volumes))

	for _, vol := range volumes {
		outputVol, err := output.NewVolume(vol, ns, nodes)
		if err != nil {
			return err
		}

		outputVols = append(outputVols, outputVol)
	}

	switch len(outputVols) {
	case 1:
		return c.display.DescribeVolume(ctx, c.writer, outputVols[0])
	default:
		return c.display.DescribeListVolumes(ctx, c.writer, outputVols)
	}
}

// getNodeMapping fetches the list of nodes from the API and builds a map from
// their ID to the full resource.
func (c *volumeCommand) getNodeMapping(ctx context.Context) (map[id.Node]*node.Resource, error) {
	nodeList, err := c.client.GetListNodes(ctx)
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
		Short:   "Show detailed information for a volume",
		Example: `
$ storageos describe volume --namespace my-namespace-name my-volume-name
`,

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

			if len(args) < 1 {
				return errMissingVolume
			}

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

	return cobraCommand
}
