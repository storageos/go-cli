package cluster

import (
	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

// NewClusterCommand returns a cobra command for `cluster` subcommands
func NewClusterCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Manage clusters",
		Args:  cli.NoArgs,
		RunE:  storageosCli.ShowHelp,
	}
	cmd.AddCommand(
		command.WithAlias(newCreateCommand(storageosCli), command.CreateAliases...),
		command.WithAlias(newInspectCommand(storageosCli), command.InspectAliases...),
		command.WithAlias(newRemoveCommand(storageosCli), command.RemoveAliases...),
		command.WithAlias(newHealthCommand(storageosCli), command.HealthAliases...),
		newConnectivityCommand(storageosCli),
	)
	return cmd
}
