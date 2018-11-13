package policy

import (
	"github.com/spf13/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

// NewPolicyCommand returns a cobra command for `policy` subcommands
func NewPolicyCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage policies",
		Args:  cli.NoArgs,
		RunE:  storageosCli.ShowHelp,
	}
	cmd.AddCommand(
		command.WithAlias(newCreateCommand(storageosCli), command.CreateAliases...),
		command.WithAlias(newInspectCommand(storageosCli), command.InspectAliases...),
		command.WithAlias(newListCommand(storageosCli), command.ListAliases...),
		command.WithAlias(newRemoveCommand(storageosCli), command.RemoveAliases...),
	)
	return cmd
}
