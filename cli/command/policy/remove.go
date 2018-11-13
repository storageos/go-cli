package policy

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

func newRemoveCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] POLICY [POLICY...]",
		Aliases: []string{"remove"},
		Short:   "Remove one or more polic(y|ies)",
		Args:    cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(storageosCli, args)
		},
	}

	return cmd
}

func runRemove(storageosCli *command.StorageOSCli, policies []string) error {
	client := storageosCli.Client()
	status := 0

	for _, policy := range policies {
		params := types.DeleteOptions{
			Name:    policy,
			Context: context.Background(),
		}

		if err := client.PolicyDelete(params); err != nil {
			fmt.Fprintf(storageosCli.Err(), "%s\n", err)
			status = 1
			continue
		}
		fmt.Fprintf(storageosCli.Out(), "%s\n", policy)
	}

	if status != 0 {
		return cli.StatusError{StatusCode: status}
	}
	return nil
}
