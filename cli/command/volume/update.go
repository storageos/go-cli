package volume

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"
	"github.com/storageos/go-cli/pkg/validation"

	"context"
)

const (
	flagDescription = "description"
	flagSize        = "size"
	flagLabelAdd    = "label-add"
	flagLabelRemove = "label-rm"
)

type updateOptions struct {
	description string
	size        int
	labels      opts.ListOpts
}

var (
	errNoSizeChange = errors.New("size was already set to the requested value")
)

func newUpdateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := updateOptions{
		labels: opts.NewListOpts(opts.ValidationPipeline(opts.ValidateEnv, opts.ValidateLabel)),
	}

	cmd := &cobra.Command{
		Use:   "update [OPTIONS] VOLUME",
		Short: "Update a volume",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(storageosCli, cmd.Flags(), args[0])
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.description, flagDescription, "d", "", `Volume description`)
	flags.IntVarP(&opt.size, flagSize, "s", 5, "Volume size in GB")
	flags.Var(&opt.labels, flagLabelAdd, "Add or update a volume label (key=value)")
	labelKeys := opts.NewListOpts(nil)
	flags.Var(&labelKeys, flagLabelRemove, "Remove a volume label if exists")
	return cmd
}

func runUpdate(storageosCli *command.StorageOSCli, flags *pflag.FlagSet, ref string) error {
	success := func(_ string) {
		fmt.Fprintln(storageosCli.Out(), ref)
	}
	return updateVolumes(storageosCli, []string{ref}, mergeVolumeUpdate(flags), success)
}

func updateVolumes(storageosCli *command.StorageOSCli, refs []string, mergeVolume func(volume *types.Volume) error, success func(name string)) error {
	client := storageosCli.Client()
	ctx := context.Background()

	for _, ref := range refs {

		namespace, name, err := validation.ParseRefWithDefault(ref)
		if err != nil {
			return err
		}

		volume, err := client.Volume(namespace, name)
		if err != nil {
			return err
		}

		err = mergeVolume(volume)
		if err != nil {
			return err
		}
		params := types.VolumeUpdateOptions{
			Name:        volume.Name,
			Namespace:   volume.Namespace,
			Description: volume.Description,
			Size:        volume.Size,
			Labels:      volume.Labels,
			Context:     ctx,
		}
		_, err = client.VolumeUpdate(params)
		if err != nil {
			return err
		}
		success(name)
	}
	return nil
}

func mergeVolumeUpdate(flags *pflag.FlagSet) func(*types.Volume) error {
	return func(volume *types.Volume) error {
		if flags.Changed(flagDescription) {
			str, err := flags.GetString(flagDescription)
			if err != nil {
				return err
			}
			volume.Description = str
		}
		if flags.Changed(flagSize) {
			gb, err := flags.GetInt(flagSize)
			if err != nil {
				return err
			}
			volume.Size = gb
		}
		if volume.Labels == nil {
			volume.Labels = make(map[string]string)
		}
		if flags.Changed(flagLabelAdd) {
			labels := flags.Lookup(flagLabelAdd).Value.(*opts.ListOpts).GetAll()
			for k, v := range opts.ConvertKVStringsToMap(labels) {
				volume.Labels[k] = v
			}
		}
		if flags.Changed(flagLabelRemove) {
			keys := flags.Lookup(flagLabelRemove).Value.(*opts.ListOpts).GetAll()
			for _, k := range keys {
				// if a key doesn't exist, fail the command explicitly
				if _, exists := volume.Labels[k]; !exists {
					return fmt.Errorf("key %s doesn't exist in volume's labels", k)
				}
				delete(volume.Labels, k)
			}
		}
		return nil
	}
}
