package pool

import (
	"fmt"

	"github.com/dnephin/cobra"
	"github.com/spf13/pflag"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"

	"context"
)

const (
	flagDescription    = "description"
	flagDefault        = "default"
	flagLabelAdd       = "label-add"
	flagLabelRemove    = "label-rm"
	flagNodeSelector   = "node-selector"
	flagDeviceSelector = "device-selector"
)

type updateOptions struct {
	name           string
	description    string
	isDefault      bool
	nodeSelector   string
	deviceSelector string
	labels         opts.ListOpts
}

func newUpdateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := updateOptions{
		labels: opts.NewListOpts(opts.ValidationPipeline(opts.ValidateEnv, opts.ValidateLabel)),
	}

	cmd := &cobra.Command{
		Use:   "update [OPTIONS] POOL",
		Short: "Update a pool",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(storageosCli, cmd.Flags(), args[0])
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&opt.nodeSelector, flagNodeSelector, "", "Node selector")
	flags.StringVar(&opt.deviceSelector, flagDeviceSelector, "", "Device selector (on filtered nodes)")
	flags.BoolVar(&opt.isDefault, flagDefault, false, "Set as default pool")

	flags.StringVarP(&opt.description, flagDescription, "d", "", `Volume description`)
	flags.Var(&opt.labels, flagLabelAdd, "Add or update a volume label (key=value)")
	labelKeys := opts.NewListOpts(nil)
	flags.Var(&labelKeys, flagLabelRemove, "Remove a volume label if exists")
	return cmd
}

func runUpdate(storageosCli *command.StorageOSCli, flags *pflag.FlagSet, ref string) error {
	success := func(_ string) {
		fmt.Fprintln(storageosCli.Out(), ref)
	}
	return updatePools(storageosCli, []string{ref}, mergePoolUpdate(flags), success)
}

func updatePools(storageosCli *command.StorageOSCli, refs []string, mergePool func(pool *types.Pool) error, success func(name string)) error {
	client := storageosCli.Client()
	ctx := context.Background()

	for _, ref := range refs {
		pool, err := client.Pool(ref)
		if err != nil {
			return err
		}

		err = mergePool(pool)
		if err != nil {
			return err
		}
		params := types.PoolOptions{
			ID:             pool.ID,
			Name:           pool.Name,
			Description:    pool.Description,
			Default:        pool.Default,
			NodeSelector:   pool.NodeSelector,
			DeviceSelector: pool.DeviceSelector,
			Labels:         pool.Labels,
			Context:        ctx,
		}
		_, err = client.PoolUpdate(params)
		if err != nil {
			return err
		}
		success(ref)
	}
	return nil
}

func mergePoolUpdate(flags *pflag.FlagSet) func(*types.Pool) error {
	return func(pool *types.Pool) error {
		if flags.Changed(flagDescription) {
			str, err := flags.GetString(flagDescription)
			if err != nil {
				return err
			}
			pool.Description = str
		}

		if flags.Changed(flagNodeSelector) {
			str, err := flags.GetString(flagNodeSelector)
			if err != nil {
				return err
			}
			pool.NodeSelector = str
		}

		if flags.Changed(flagDeviceSelector) {
			str, err := flags.GetString(flagDeviceSelector)
			if err != nil {
				return err
			}
			pool.DeviceSelector = str
		}

		if flags.Changed(flagDefault) {
			b, err := flags.GetBool(flagDefault)
			if err != nil {
				return err
			}
			pool.Default = b
		}

		if pool.Labels == nil {
			pool.Labels = make(map[string]string)
		}
		if flags.Changed(flagLabelAdd) {
			labels := flags.Lookup(flagLabelAdd).Value.(*opts.ListOpts).GetAll()
			for k, v := range opts.ConvertKVStringsToMap(labels) {
				pool.Labels[k] = v
			}
		}
		if flags.Changed(flagLabelRemove) {
			keys := flags.Lookup(flagLabelRemove).Value.(*opts.ListOpts).GetAll()
			for _, k := range keys {
				// if a key doesn't exist, fail the command explicitly
				if _, exists := pool.Labels[k]; !exists {
					return fmt.Errorf("key %s doesn't exist in volume's labels", k)
				}
				delete(pool.Labels, k)
			}
		}
		return nil
	}
}
