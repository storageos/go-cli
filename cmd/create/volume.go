package create

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	units "github.com/alecthomas/units"
	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
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
	hintMaster     []string
	hintReplicas   []string

	// Core volume configuration settings
	description string
	fsType      string
	sizeStr     string
	labelPairs  []string

	writer io.Writer
}

func (c *volumeCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	// Convert the flag values to the desired types/units
	labelSet, err := labels.SetFromPairs(c.labelPairs)
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
	nsID := id.Namespace(c.namespace)

	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	if !useIDs {
		ns, err := c.client.GetNamespaceByName(ctx, c.namespace)
		if err != nil {
			return err
		}
		nsID = ns.ID
	}

	vol, err := c.client.CreateVolume(
		ctx,
		nsID,
		name,
		c.description,
		volume.FsTypeFromString(c.fsType),
		uint64(sizeBytes),
		labelSet,
	)
	if err != nil {
		return err
	}

	return c.display.CreateVolume(ctx, c.writer, vol)
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

	// If master or replica hints have been provided then rejoin the specified
	// hints and set the appropriate labels.
	if len(c.hintMaster) > 0 {
		original[volume.LabelHintMaster] = strings.Join(c.hintMaster, ",")
	}

	if len(c.hintReplicas) > 0 {
		original[volume.LabelHintReplicas] = strings.Join(c.hintReplicas, ",")
	}
}

func newVolume(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &volumeCommand{
		config: config,
		client: client,
		display: jsonformat.NewDisplayer(
			jsonformat.DefaultEncodingIndent,
		),

		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "volume",
		Short: "volume requests the creation of a new StorageOS volume",
		Example: `
$ storageos create volume --description "This volume contains the data for my app" --fs-type "ext4" --labels env=prod,rack=db-1 --size 10GiB --namespace my-namespace-name my-app

$ storageos create volume --replicas 1 --hint-master reliable-node-1,reliable-node-2 --namespace my-namespace-name my-replicated-app
`,

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errVolumeNameSpecifiedWrong
			}
			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			ns, err := c.config.Namespace()
			if err != nil {
				return err
			}

			if ns == "" {
				return errNoNamespaceSpecified
			}
			c.namespace = ns

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
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
	cobraCommand.Flags().StringSliceVar(&c.hintMaster, "hint-master", []string{}, "an optional list of preferred nodes for placement of the volume master")
	cobraCommand.Flags().StringArrayVar(&c.hintReplicas, "hint-replicas", []string{}, "an optional list of preferred nodes for placement of volume replicas")
	cobraCommand.Flags().StringSliceVarP(&c.labelPairs, "labels", "l", []string{}, "an optional set of labels to assign to the new volume, provided as a comma-separated list of key=value pairs")
	cobraCommand.Flags().Uint64VarP(&c.withReplicas, "replicas", "r", 0, "the number of replicated copies of the volume to maintain")
	cobraCommand.Flags().StringVarP(&c.sizeStr, "size", "s", "5GiB", "the capacity to provision the volume with")
	cobraCommand.Flags().BoolVar(&c.useThrottle, "throttle", false, "deprioritises the volume's traffic by reducing the rate of disk I/O")

	return cobraCommand
}
