package volume

import (
	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

// NewVolumeCommand returns a cobra command for `volume` subcommands
func NewVolumeCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "Manage volumes",
		Args:  cli.NoArgs,
		RunE:  storageosCli.ShowHelp,
	}
	cmd.AddCommand(
		command.WithAlias(newCreateCommand(storageosCli), command.CreateAliases...),
		command.WithAlias(newInspectCommand(storageosCli), command.InspectAliases...),
		command.WithAlias(newListCommand(storageosCli), command.ListAliases...),
		command.WithAlias(newUpdateCommand(storageosCli), command.UpdateAliases...),
		command.WithAlias(newRemoveCommand(storageosCli), command.RemoveAliases...),
		command.WithAlias(newMountCommand(storageosCli), "m", "mn", "mnt", "mo"),
		command.WithAlias(newUnmountCommand(storageosCli), "u", "un", "um", "umount"),
	)
	return cmd
}
