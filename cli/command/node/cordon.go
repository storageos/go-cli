package node

import (
	"fmt"
	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"strings"
)

type cordonOptions struct {
	nodes []string
}

func newCordonCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt cordonOptions

	cmd := &cobra.Command{
		Use:   "cordon NODE [NODE...]",
		Short: "Put one or more nodes into an unschedulable state",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.nodes = args
			return runCordon(storageosCli, opt)
		},
	}

	return cmd
}

func runCordon(storageosCli *command.StorageOSCli, opt cordonOptions) error {
	client := storageosCli.Client()
	failed := make([]string, 0, len(opt.nodes))

	for _, nodeID := range opt.nodes {
		n, err := client.Controller(nodeID)
		if err != nil {
			failed = append(failed, nodeID)
			continue
		}

		_, err = client.ControllerUpdate(types.ControllerUpdateOptions{
			ID:          n.ID,
			Name:        n.Name,
			Description: n.Description,
			Labels:      n.Labels,
			Cordon:      true,
		})
		if err != nil {
			failed = append(failed, nodeID)
			continue
		}

		fmt.Fprintln(storageosCli.Out(), nodeID)
	}

	if len(failed) > 0 {
		return fmt.Errorf("Failed to cordon: %s", strings.Join(failed, ", "))
	}
	return nil
}
