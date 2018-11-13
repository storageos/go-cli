package node

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

type drainOptions struct {
	nodes []string
}

func newDrainCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt drainOptions

	cmd := &cobra.Command{
		Use:   "drain NODE [NODE...]",
		Short: "Migrate volumes from one or more nodes.",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.nodes = args
			return runDrain(storageosCli, opt, true)
		},
	}

	return cmd
}

func newUndrainCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt drainOptions

	cmd := &cobra.Command{
		Use:   "undrain NODE [NODE...]",
		Short: "Stop drain on one or more nodes.",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.nodes = args
			return runDrain(storageosCli, opt, false)
		},
	}

	return cmd
}

func runDrain(storageosCli *command.StorageOSCli, opt drainOptions, drain bool) error {
	client := storageosCli.Client()
	failed := make([]string, 0, len(opt.nodes))

	for _, nodeID := range opt.nodes {
		n, err := client.Node(nodeID)
		if err != nil {
			failed = append(failed, nodeID)
			continue
		}

		_, err = client.NodeUpdate(types.NodeUpdateOptions{
			ID:          n.ID,
			Name:        n.Name,
			Description: n.Description,
			Labels:      n.Labels,
			Cordon:      n.Cordon,
			Drain:       drain,
		})
		if err != nil {
			failed = append(failed, nodeID)
			continue
		}

		fmt.Fprintln(storageosCli.Out(), nodeID)
	}

	if len(failed) > 0 {
		return fmt.Errorf("Failed to drain: %s", strings.Join(failed, ", "))
	}
	return nil
}
