package node

import (
	"fmt"
	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"strings"
)

type uncordonOptions struct {
	nodes []string
}

func newUncordonCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt uncordonOptions

	cmd := &cobra.Command{
		Use:   "uncordon NODE [NODE...]",
		Short: "Restore one or more nodes from an unschedulable state",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.nodes = args
			return runUncordon(storageosCli, opt)
		},
	}

	return cmd
}

func runUncordon(storageosCli *command.StorageOSCli, opt uncordonOptions) error {
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
			Cordon:      false,
		})
		if err != nil {
			failed = append(failed, nodeID)
			continue
		}

		fmt.Fprintln(storageosCli.Out(), nodeID)
	}

	if len(failed) > 0 {
		return fmt.Errorf("Failed to uncordon: %s", strings.Join(failed, ", "))
	}
	return nil
}
