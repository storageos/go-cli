package rule

import (
	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

// NewRuleCommand returns a cobra command for `rule` subcommands
func NewRuleCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "Manage rules",
		Args:  cli.NoArgs,
		RunE:  storageosCli.ShowHelp,
	}
	cmd.AddCommand(
		newCreateCommand(storageosCli),
		newInspectCommand(storageosCli),
		newListCommand(storageosCli),
		newUpdateCommand(storageosCli),
		newRemoveCommand(storageosCli),
	)
	return cmd
}
