package pool

import (
	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

// NewPoolCommand returns a cobra command for `pool` subcommands
func NewPoolCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool",
		Short: "Manage capacity pools",
		Args:  cli.NoArgs,
		RunE:  storageosCli.ShowHelp,
	}
	cmd.AddCommand(
		// newCreateCommand(dockerCli),
		// newInspectCommand(dockerCli),
		// newPsCommand(dockerCli),
		newListCommand(storageosCli),
	// newRemoveCommand(dockerCli),
	// newScaleCommand(dockerCli),
	// newUpdateCommand(dockerCli),
	// newLogsCommand(dockerCli),
	)
	return cmd
}
