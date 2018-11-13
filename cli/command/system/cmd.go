package system

import (
	"github.com/spf13/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

// NewSystemCommand returns a cobra command for `system` subcommands
func NewSystemCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "Manage StorageOS",
		Args:  cli.NoArgs,
		RunE:  storageosCli.ShowHelp,
	}
	cmd.AddCommand(
	// NewEventsCommand(storageosCli),
	// NewInfoCommand(storageosCli),
	// NewDiskUsageCommand(storageosCli),
	)

	return cmd
}
