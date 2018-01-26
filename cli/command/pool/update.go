package pool

import (
	"context"
	"fmt"

	"github.com/dnephin/cobra"
	"github.com/spf13/pflag"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"
)

const (
	flagActive           = "active"
	flagControllerAdd    = "controller-add"
	flagControllerRemove = "controller-rm"
	flagDefault          = "default"
	flagDefaultDriver    = "default-driver"
	flagDescription      = "description"
	flagLabelAdd         = "label-add"
	flagLabelRemove      = "label-rm"
)

type updateOptions struct {
	active            bool
	addControllers    opts.ListOpts
	removeControllers opts.ListOpts
	defaultDriver     string
	description       string
	isDefault         bool
	addLabels         opts.ListOpts
	removeLabels      opts.ListOpts
}

func newUpdateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := updateOptions{
		addControllers:    opts.NewListOpts(opts.ValidateEnv),
		removeControllers: opts.NewListOpts(nil),
		addLabels:         opts.NewListOpts(opts.ValidateEnv),
		removeLabels:      opts.NewListOpts(nil),
	}

	cmd := &cobra.Command{
		Use:   "update [OPTIONS] POOL",
		Short: "Update a capacity pool",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(storageosCli, cmd.Flags(), args[0])
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opt.active, flagActive, true, "Enable or disable the pool")
	flags.Var(&opt.addControllers, flagControllerAdd, "Add controllers to the pool")
	flags.Var(&opt.removeControllers, flagControllerRemove, "Remove controllers from the pool")
	flags.BoolVar(&opt.isDefault, flagDefault, false, "Set as default pool")
	flags.StringVar(&opt.defaultDriver, flagDefaultDriver, "", "Default driver for the pool")
	flags.StringVarP(&opt.description, flagDescription, "d", "", "Pool description")
	flags.Var(&opt.addLabels, flagLabelAdd, "Add or update pool labels (key=value)")
	flags.Var(&opt.removeLabels, flagLabelRemove, "Remove pool labels")

	return cmd
}

func runUpdate(storageosCli *command.StorageOSCli, flags *pflag.FlagSet, ref string) error {
	client := storageosCli.Client()
	ctx := context.Background()

	pool, err := client.Pool(ref)
	if err != nil {
		return fmt.Errorf("failed to find pool (%s): %v", ref, err)
	}

	// Ensure that there is a slice before attempting to modify controllers
	if pool.ControllerNames == nil {
		pool.ControllerNames = make([]string, 0)
	}

	if flags.Changed(flagControllerAdd) {
		controllers := flags.Lookup(flagControllerAdd).Value.(*opts.ListOpts).GetAll()
		for _, controller := range controllers {
			pool.ControllerNames = append(pool.ControllerNames, controller)
		}
	}

	if flags.Changed(flagControllerRemove) {
		controllers := flags.Lookup(flagControllerRemove).Value.(*opts.ListOpts).GetAll()
		for _, controller := range controllers {
			for i, c := range pool.ControllerNames {
				if controller == c {
					pool.ControllerNames = append(pool.ControllerNames[:i], pool.ControllerNames[i+1:]...)
					break
				}
				// Fail if any controller to be removed doesn't exist
				return fmt.Errorf("%s is not a member of the pool", controller)
			}
		}
	}

	if flags.Changed(flagActive) {
		active, err := flags.GetBool(flagActive)
		if err != nil {
			return fmt.Errorf("error retrieving value of active flag: %v", err)
		}
		pool.Active = active
	}

	if flags.Changed(flagDefault) {
		def, err := flags.GetBool(flagDefault)
		if err != nil {
			return fmt.Errorf("error retrieving value of default flag: %v", err)
		}
		pool.Default = def
	}

	if flags.Changed(flagDefaultDriver) {
		driver, err := flags.GetString(flagDefaultDriver)
		if err != nil {
			return fmt.Errorf("error retrieving name of default driver to use: %v", err)
		}
		pool.DefaultDriver = driver
	}

	if flags.Changed(flagDescription) {
		desc, err := flags.GetString(flagDescription)
		if err != nil {
			return fmt.Errorf("error retrieving description to use: %v", err)
		}
		pool.Description = desc
	}

	// Ensure there is a label map before attempting to edit
	if pool.Labels == nil {
		pool.Labels = make(map[string]string)
	}

	if flags.Changed(flagLabelAdd) {
		labels := flags.Lookup(flagLabelAdd).Value.(*opts.ListOpts).GetAll()

		for label, value := range opts.ConvertKVStringsToMap(labels) {
			pool.Labels[label] = value
		}
	}

	if flags.Changed(flagLabelRemove) {
		keys := flags.Lookup(flagLabelRemove).Value.(*opts.ListOpts).GetAll()
		// Fail if any label to be removed doesn't exist
		for _, label := range keys {
			if _, exists := pool.Labels[label]; !exists {
				return fmt.Errorf("key %s doesn't exist in the pool's labels", label)
			}
			delete(pool.Labels, label)
		}
	}

	if _, err = client.PoolUpdate(types.PoolUpdateOptions{
		ID:              pool.ID,
		Name:            pool.Name,
		Description:     pool.Description,
		Default:         pool.Default,
		DefaultDriver:   pool.DefaultDriver,
		ControllerNames: pool.ControllerNames,
		DriverNames:     pool.DriverNames,
		Active:          pool.Active,
		Labels:          pool.Labels,
		Context:         ctx,
	}); err != nil {
		return fmt.Errorf("Failed to update pool (%s): %v", ref, err)
	}

	fmt.Fprintln(storageosCli.Out(), ref)
	return nil
}
