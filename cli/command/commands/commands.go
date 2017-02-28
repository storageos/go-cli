package commands

import (
	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/pool"
	"github.com/storageos/go-cli/cli/command/volume"
)

// AddCommands adds all the commands from cli/command to the root command
func AddCommands(cmd *cobra.Command, storageosCli *command.StorageOSCli) {
	cmd.AddCommand(
		volume.NewVolumeCommand(storageosCli),
		pool.NewPoolCommand(storageosCli),
	)
}
