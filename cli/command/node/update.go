package node

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"
)

const (
	flagDescription = "description"
	flagLabelAdd    = "label-add"
	flagLabelRemove = "label-rm"
)

type updateOptions struct {
	description string
	addLabel    string
	rmLabel     string
}

func newUpdateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := updateOptions{}

	cmd := &cobra.Command{
		Use:   "update [OPTIONS] NODE",
		Short: "Update a node",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(storageosCli, opt, args[0])
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.description, flagDescription, "d", "", `Node description`)
	flags.StringVar(&opt.addLabel, flagLabelAdd, "", "Add or update a node label (key=value)")
	flags.StringVar(&opt.rmLabel, flagLabelRemove, "", "Remove a node label if exists")
	return cmd
}

func runUpdate(storageosCli *command.StorageOSCli, opt updateOptions, nodeID string) error {
	client := storageosCli.Client()

	n, err := client.Node(nodeID)
	if err != nil {
		return fmt.Errorf("Failed to find node (%s): %v", nodeID, err)
	}

	if opt.description != "" {
		n.Description = opt.description
	}

	if opt.rmLabel != "" {
		delete(n.Labels, opt.rmLabel)
	}

	if opt.addLabel != "" {
		if err := updateLabel(n, opt.addLabel); err != nil {
			return err
		}
	}

	if _, err = client.NodeUpdate(types.NodeUpdateOptions{
		ID:          n.ID,
		Name:        n.Name,
		Description: n.Description,
		Labels:      n.Labels,
		Cordon:      n.Cordon,
	}); err != nil {
		return fmt.Errorf("Failed to update node (%s): %v", nodeID, err)
	}

	fmt.Fprintln(storageosCli.Out(), nodeID)
	return nil
}

func updateLabel(n *types.Node, labels string) error {
	if _, err := opts.ValidateLabel(labels); err != nil {
		return err
	}

	arr := strings.Split(labels, "=")

	if len(arr) != 2 || arr[0] == "" || arr[1] == "" {
		return errors.New("Bad label format: " + labels)
	}

	// If labels map is uninitialized, initialize it.
	if n.Labels == nil {
		n.Labels = make(map[string]string)
	}
	n.Labels[arr[0]] = arr[1]

	return nil
}
