package user

import (
	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

// NewUserCommand returns a cobra command for `user` subcommands
func NewUserCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
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
