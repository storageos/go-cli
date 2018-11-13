package node

import (
	"github.com/spf13/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

func newDeleteCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	return &cobra.Command{
		Use:   "delete NODE ",
		Short: "Remove an offline node from the cluster.",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(storageosCli, args[0])
		},
	}
}

func runDelete(storageosCli *command.StorageOSCli, nodeID string) error {
	client := storageosCli.Client()

	opts := types.DeleteOptions{
		Name: nodeID,
	}

	return client.NodeDelete(opts)
}
