package create

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/alecthomas/units"
	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/volume"
)

var (
	errVolumeNameSpecifiedWrong = errors.New("must specify exactly one name for the volume")
	errNoNamespaceSpecified     = errors.New("must specify a namespace to create the volume in")
)

type volumeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string

	// Known label configuration aliases
	useCaching     bool
	useCompression bool
	useThrottle    bool
	withReplicas   uint64

	// Core volume configuration settings
	description string
	fsType      string
	sizeStr     string
	labelPairs  []string

	useAsync bool

	writer io.Writer
}

func (c *volumeCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {

	var createVolumeRequestParams *apiclient.CreateVolumeRequestParams

	// If asynchrony is specified then source the timeout and initialise the params.
	if c.useAsync {
		timeout, err := c.config.CommandTimeout()
		if err != nil {
			return err
		}

		createVolumeRequestParams = &apiclient.CreateVolumeRequestParams{
			AsyncMax: timeout,
		}
	}

	// Convert the flag values to the desired types/units
	labelSet, err := labels.NewSetFromPairs(c.labelPairs)
	if err != nil {
		return err
	}

	// Set any of the internal well-known labels by their specified flag values
	c.setKnownLabels(labelSet)

	sizeBytes, err := units.ParseStrictBytes(c.sizeStr)
	if err != nil {
		return err
	}
	if sizeBytes <= 0 {
		return fmt.Errorf("provided invalid volume size of %d bytes", sizeBytes)
	}

	name := args[0]
	namespaceID := id.Namespace(c.namespace)

	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	var ns *namespace.Resource

	if !useIDs {
		ns, err = c.client.GetNamespaceByName(ctx, c.namespace)
		if err != nil {
			return err
		}
		namespaceID = ns.ID
	} else {
		ns, err = c.client.GetNamespace(ctx, namespaceID)
		if err != nil {
			return err
		}
	}

	vol, err := c.client.CreateVolume(
		ctx,
		namespaceID,
		name,
		c.description,
		volume.FsTypeFromString(c.fsType),
		uint64(sizeBytes),
		labelSet,
		createVolumeRequestParams,
	)
	if err != nil {
		return err
	}

	// If the request was async then write our "request submitted" message
	// and return.
	if c.useAsync {
		return c.display.AsyncRequest(ctx, c.writer)
	}

	nodes, err := c.getNodeMapping(ctx)
	if err != nil {
		return err
	}

	return c.display.CreateVolume(ctx, c.writer, output.NewVolume(vol, ns, nodes))
}

// getNodeMapping fetches the list of nodes from the API and builds a map from
// their ID to the full resource.
func (c *volumeCommand) getNodeMapping(ctx context.Context) (map[id.Node]*node.Resource, error) {
	nodeList, err := c.client.GetListNodesByUID(ctx)
	if err != nil {
		return nil, err
	}

	nodes := map[id.Node]*node.Resource{}
	for _, n := range nodeList {
		nodes[n.ID] = n
	}

	return nodes, nil
}

func (c *volumeCommand) setKnownLabels(original labels.Set) {
	// Set the no-cache label
	original[volume.LabelNoCache] = strconv.FormatBool(!c.useCaching)

	// Set the no-compress label
	original[volume.LabelNoCompress] = strconv.FormatBool(!c.useCompression)

	// Set the throttle label if desired
	original[volume.LabelThrottle] = strconv.FormatBool(c.useThrottle)

	// Set the replication label if replicas are desired
	if c.withReplicas > 0 {
		original[volume.LabelReplicas] = strconv.FormatUint(c.withReplicas, 10)
	}
}

func newVolume(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "volume",
		Short: "Provision a new volume",
		Example: `
$ storageos create volume --description "This volume contains the data for my app" --fs-type "ext4" --labels env=prod,rack=db-1 --size 10GiB --namespace my-namespace-name my-app

$ storageos create volume --replicas 1 --namespace my-namespace-name my-replicated-app
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errVolumeNameSpecifiedWrong
			}
			return nil
		}),
		PreRunE: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, _ []string) error {
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
				runwrappers.HandleLicenceError(client),
			)(c.runWithCtx)

			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	cobraCommand.Flags().BoolVar(&c.useCaching, "cache", true, "caches volume data")
	cobraCommand.Flags().BoolVar(&c.useCompression, "compress", true, "compress data stored by the volume at rest and during transit")
	cobraCommand.Flags().StringVarP(&c.description, "description", "d", "", "a human-friendly description to give the volume")
	cobraCommand.Flags().StringVarP(&c.fsType, "fs-type", "f", "ext4", "the filesystem to format the new volume with once provisioned")
	cobraCommand.Flags().StringSliceVarP(&c.labelPairs, "labels", "l", []string{}, "an optional set of labels to assign to the new volume, provided as a comma-separated list of key=value pairs")
	cobraCommand.Flags().Uint64VarP(&c.withReplicas, "replicas", "r", 0, "the number of replicated copies of the volume to maintain")
	cobraCommand.Flags().StringVarP(&c.sizeStr, "size", "s", "5GiB", "the capacity to provision the volume with")
	cobraCommand.Flags().BoolVar(&c.useThrottle, "throttle", false, "deprioritises the volume's traffic by reducing the rate of disk I/O")

	flagutil.SupportAsync(cobraCommand.Flags(), &c.useAsync)

	return cobraCommand
}
