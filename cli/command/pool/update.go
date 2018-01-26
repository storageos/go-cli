package pool

import (
	"context"
	"fmt"
	"strings"

	"github.com/dnephin/cobra"
	"github.com/spf13/pflag"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
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
	active           bool
	controllerAdd    string
	controllerRemove string
	defaultDriver    string
	description      string
	isDefault        bool
	labelAdd         string
	labelRemove      string
}

func newUpdateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := updateOptions{}

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
	flags.StringVar(&opt.controllerAdd, flagControllerAdd, "", "Add a controller to the capacity pool")
	flags.StringVar(&opt.controllerRemove, flagControllerRemove, "", "Remove a controller from the capacity pool")
	flags.BoolVar(&opt.isDefault, flagDefault, false, "Set as default capacity pool")
	flags.StringVar(&opt.defaultDriver, flagDefaultDriver, "", "Default driver for the capacity pool")
	flags.StringVarP(&opt.description, flagDescription, "d", "", "Pool description")
	flags.StringVar(&opt.labelAdd, flagLabelAdd, "", "Add or update a capacity pool label (key=value)")
	flags.StringVar(&opt.labelRemove, flagLabelRemove, "", "Remove a capacity pool label")

	return cmd
}

func runUpdate(storageosCli *command.StorageOSCli, flags *pflag.FlagSet, ref string) error {
	client := storageosCli.Client()
	ctx := context.Background()

	pool, err := client.Pool(ref)
	if err != nil {
		return fmt.Errorf("Failed to find pool (%s): %v", ref, err)
	}

	// Ensure that there is a slice before attempting to modify controllers
	if pool.ControllerNames == nil {
		pool.ControllerNames = make([]string, 0)
	}

	if flags.Changed(flagControllerAdd) {
		controller, err := flags.GetString(flagControllerAdd)
		if err != nil {
			return fmt.Errorf("Error retrieving name of controller to add: %v", err)
		}
		if controller != "" {
			pool.ControllerNames = append(pool.ControllerNames, controller)
		}
	}

	if flags.Changed(flagControllerRemove) {
		controller, err := flags.GetString(flagControllerRemove)
		if err != nil {
			return fmt.Errorf("Error retrieving name of controller to remove: %v", err)
		}
		if controller != "" {
			for i, c := range pool.ControllerNames {
				if controller == c {
					pool.ControllerNames = append(pool.ControllerNames[:i], pool.ControllerNames[i+1:]...)
					break
				}
			}
		}
	}

	if flags.Changed(flagActive) {
		active, err := flags.GetBool(flagActive)
		if err != nil {
			return fmt.Errorf("Error retrieving value of active flag: %v", err)
		}
		pool.Active = active
	}

	if flags.Changed(flagDefault) {
		def, err := flags.GetBool(flagDefault)
		if err != nil {
			return fmt.Errorf("Error retrieving value of default flag: %v", err)
		}
		pool.Default = def
	}

	if flags.Changed(flagDefaultDriver) {
		driver, err := flags.GetString(flagDefaultDriver)
		if err != nil {
			return fmt.Errorf("Error retrieving name of default driver to use: %v", err)
		}
		pool.DefaultDriver = driver
	}

	if flags.Changed(flagDescription) {
		desc, err := flags.GetString(flagDescription)
		if err != nil {
			return fmt.Errorf("Error retrieving description to use: %v", err)
		}
		pool.Description = desc
	}

	// Ensure there is a label map before attempting to edit
	if pool.Labels == nil {
		pool.Labels = make(map[string]string)
	}

	if flags.Changed(flagLabelAdd) {
		label, err := flags.GetString(flagLabelAdd)
		if err != nil {
			return fmt.Errorf("Error retrieving label pair to add: %v", err)
		}
		if label != "" {
			pair := strings.Split(label, "=")

			if len(pair) != 2 || pair[0] == "" || pair[1] == "" {
				return fmt.Errorf("Bad label format: %s", label)
			}

			pool.Labels[pair[0]] = pair[1]
		}
	}

	if flags.Changed(flagLabelRemove) {
		label, err := flags.GetString(flagLabelRemove)
		if err != nil {
			return fmt.Errorf("Error retrieving label to remove: %v", err)
		}
		if label != "" {
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
