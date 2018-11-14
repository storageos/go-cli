package namespace

import (
	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

// NewNamespaceCommand returns a cobra command for `namespace` subcommands
func NewNamespaceCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namespace",
		Short: "Manage namespaces",
		Args:  cli.NoArgs,
		RunE:  storageosCli.ShowHelp,
	}
	cmd.AddCommand(
		command.WithAlias(newCreateCommand(storageosCli), command.CreateAliases...),
		command.WithAlias(newInspectCommand(storageosCli), command.InspectAliases...),
		command.WithAlias(newListCommand(storageosCli), command.ListAliases...),
		command.WithAlias(newUpdateCommand(storageosCli), command.UpdateAliases...),
		command.WithAlias(newRemoveCommand(storageosCli), command.RemoveAliases...),
	)
	return cmd
}
