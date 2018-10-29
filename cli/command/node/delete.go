package node

import (
	"time"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

const minDeleteTimeout = 60 * time.Second

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
	if storageosCli.GetTimeout() < 60*time.Second {
		client.SetTimeout(minDeleteTimeout)
	}

	opts := types.DeleteOptions{
		Name: nodeID,
	}

	return client.NodeDelete(opts)
}
