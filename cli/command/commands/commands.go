package commands

import (
	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/cluster"
	"github.com/storageos/go-cli/cli/command/namespace"
	"github.com/storageos/go-cli/cli/command/node"
	"github.com/storageos/go-cli/cli/command/pool"
	"github.com/storageos/go-cli/cli/command/rule"
	"github.com/storageos/go-cli/cli/command/system"
	"github.com/storageos/go-cli/cli/command/volume"
)

// AddCommands adds all the commands from cli/command to the root command
func AddCommands(cmd *cobra.Command, storageosCli *command.StorageOSCli) {
	cmd.AddCommand(
		namespace.NewNamespaceCommand(storageosCli),
		pool.NewPoolCommand(storageosCli),
		rule.NewRuleCommand(storageosCli),
		volume.NewVolumeCommand(storageosCli),
		node.NewNodeCommand(storageosCli),

		// system
		// system.NewSystemCommand(storageosCli),
		system.NewVersionCommand(storageosCli),

		// clustering
		cluster.NewClusterCommand(storageosCli),
	)
}
