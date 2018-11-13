package licence

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

func newRemoveCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm",
		Aliases: []string{"remove"},
		Short:   "Remove the current licence",
		Args:    cli.RequiresMinArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(storageosCli, args)
		},
	}

	return cmd
}

func runRemove(storageosCli *command.StorageOSCli, policies []string) error {
	client := storageosCli.Client()
	status := 0

	if err := client.LicenceDelete(); err != nil {
		fmt.Fprintf(storageosCli.Err(), "%s\n", err)
		status = 1
	}

	if status != 0 {
		return cli.StatusError{StatusCode: status}
	}
	return nil
}
