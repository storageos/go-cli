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
		newCreateCommand(storageosCli),
		newInspectCommand(storageosCli),
		newListCommand(storageosCli),
		//newUpdateCommand(storageosCli),
		newRemoveCommand(storageosCli),
	)
	return cmd
}
