package create

import (
	"context"
	"errors"
	"fmt"
	"io"

	units "github.com/alecthomas/units"
	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/volume"
)

var (
	errNoVolumeNameSpecified = errors.New("must specify a single name for the volume")
	errNoNamespaceSpecified  = errors.New("must specify a namespace to create the volume in")
)

type volumeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	namespace string

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

	sizeBytes, err := units.ParseStrictBytes(c.sizeStr)
	if err != nil {
		return err
	}
	if sizeBytes <= 0 {
		return fmt.Errorf("provided invalid volume size of %d bytes", sizeBytes)
	}

	name := args[0]
	nsID := id.Namespace(c.namespace)

	if useIDs, err := c.config.UseIDs(); !useIDs || err != nil {
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
`,

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errNoVolumeNameSpecified
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

	cobraCommand.Flags().StringVarP(&c.description, "description", "d", "", "a human-friendly description to give the StorageOS volume")
	cobraCommand.Flags().StringVarP(&c.fsType, "fs-type", "f", "ext4", "the filesystem to format the new volume with once provisioned")
	cobraCommand.Flags().StringSliceVarP(&c.labelPairs, "labels", "l", []string{}, "an optional set of labels to assign to the new volume, provided as a comma-separated list of key=value pairs")
	cobraCommand.Flags().StringVarP(&c.sizeStr, "size", "s", "5GiB", "the capacity to provision the volume with")

	return cobraCommand
}
