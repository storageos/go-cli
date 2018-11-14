package licence

import (
	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

// NewLicenceCommand returns a cobra command for `policy` subcommands
func NewLicenceCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licence",
		Short: "Manage the licence",
		Args:  cli.NoArgs,
		RunE:  storageosCli.ShowHelp,
	}
	cmd.AddCommand(
		command.WithAlias(newApplyCommand(storageosCli), command.ApplyAliases...),
		command.WithAlias(newInspectCommand(storageosCli), command.InspectAliases...),
		command.WithAlias(newRemoveCommand(storageosCli), command.RemoveAliases...),
	)
	return cmd
}
