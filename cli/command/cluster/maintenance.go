package cluster

import (
	"fmt"

	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/inspect"
)

type maintenanceOptions struct {
	format string
}

func newMaintenanceCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt maintenanceOptions

	cmd := &cobra.Command{
		Use:   "maintenance [OPTIONS] enable|disable|inspect",
		Short: "Enable|disable maintenance mode for the cluster",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "inspect":
				return getMaintenance(storageosCli, opt)
			case "enable":
				return enableMaintenance(storageosCli, opt)
			case "disable":
				return disableMaintenance(storageosCli, opt)
			default:
				return fmt.Errorf("wrong argument %s", args[0])
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.format, "format", "f", "", "Format the output using the given Go template.")
	return cmd
}

func getMaintenance(storageosCli *command.StorageOSCli, opt maintenanceOptions) error {
	client := storageosCli.Client()

	getFunc := func(string) (interface{}, []byte, error) {
		i, err := client.Maintenance()
		return i, nil, err
	}

	return inspect.Inspect(storageosCli.Out(), []string{""}, opt.format, getFunc)
}

func enableMaintenance(storageosCli *command.StorageOSCli, opt maintenanceOptions) error {
	return storageosCli.Client().EnableMaintenance()
}

func disableMaintenance(storageosCli *command.StorageOSCli, opt maintenanceOptions) error {
	return storageosCli.Client().DisableMaintenance()
}
