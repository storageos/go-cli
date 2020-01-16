package create

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	units "github.com/alecthomas/units"
	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/volume"
)

type volumeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	description string
	fsType      string
	sizeStr     string
	labelPairs  []string

	usingID bool

	writer io.Writer
}

func (c *volumeCommand) run(cmd *cobra.Command, args []string) error {
	timeout, err := c.config.CommandTimeout()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Convert the flag values to the desired types/units
	labelSet, err := labels.SetFromPairs(c.labelPairs)
	if err != nil {
		return err
	}

	sizeBytes, err := units.ParseStrictBytes(c.sizeStr)
	if err != nil {
		return err
	}
	if sizeBytes <= 0 {
		return fmt.Errorf("provided invalid volume size of %d bytes", sizeBytes)
	}

	ref := args[0]

	if c.usingID {
		vol, err := c.createVolumeByNamespaceID(ctx, ref, uint64(sizeBytes), labelSet)
		if err != nil {
			return err
		}

		return c.display.CreateVolume(ctx, c.writer, vol)
	}

	vol, err := c.createVolume(ctx, ref, uint64(sizeBytes), labelSet)
	if err != nil {
		return err
	}

	return c.display.CreateVolume(ctx, c.writer, vol)
}

func (c *volumeCommand) createVolume(ctx context.Context, ref string, sizeBytes uint64, labelSet labels.Set) (*volume.Resource, error) {
	nsName, volName, err := volume.ParseReferenceName(ref)
	if err != nil {
		return nil, err
	}

	ns, err := c.client.GetNamespaceByName(ctx, nsName)
	if err != nil {
		return nil, err
	}

	return c.client.CreateVolume(ctx,
		ns.ID,
		volName,
		c.description,
		volume.FsTypeFromString(c.fsType),
		sizeBytes,
		labelSet,
	)
}

func (c *volumeCommand) createVolumeByNamespaceID(ctx context.Context, ref string, sizeBytes uint64, labelSet labels.Set) (*volume.Resource, error) {
	parts := strings.Split(ref, "/")
	switch len(parts) {
	case 2:
		// OK
	case 1:
		return nil, errors.New("requested to use namespace ID for volume creation but no namespace specified")
	default:
		return nil, errors.New("invalid volume reference string")
	}

	return c.client.CreateVolume(ctx,
		id.Namespace(parts[0]),
		parts[1],
		c.description,
		volume.FsTypeFromString(c.fsType),
		sizeBytes,
		labelSet,
	)
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
		$ storageos create volume --description "This volume contains the data for my app" --fs-type "ext4" --labels env=prod,rack=db-1 --size 10GiB my-namespace-name/my-app
		`,

		Args: cobra.ExactArgs(1),

		RunE: c.run,

		SilenceUsage: true,
	}

	cobraCommand.Flags().StringVarP(&c.description, "description", "d", "", "a human-friendly description to give the StorageOS volume")
	cobraCommand.Flags().StringVarP(&c.fsType, "fs-type", "f", "ext4", "the filesystem to format the new volume with once provisioned")
	cobraCommand.Flags().BoolVar(&c.usingID, "use-id", false, "specify the namespace to create the StorageOS volume in by its ID instead of its name")
	cobraCommand.Flags().StringSliceVarP(&c.labelPairs, "labels", "l", []string{}, "an optional set of labels to assign to the new volume, provided as a comma-separated list of key=value pairs")
	cobraCommand.Flags().StringVarP(&c.sizeStr, "size", "s", "5GiB", "the capacity to provision the volume with")

	return cobraCommand
}
