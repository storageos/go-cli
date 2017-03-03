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
		newCreateCommand(storageosCli),
		newInspectCommand(storageosCli),
		newListCommand(storageosCli),
		newUpdateCommand(storageosCli),
		newRemoveCommand(storageosCli),
	)
	return cmd
}
