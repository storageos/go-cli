package cluster

import (
	"fmt"

	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"

	"github.com/storageos/go-cli/discovery"
)

type removeOptions struct {
	force bool
	rules []string
}

func newRemoveCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt removeOptions

	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] CLUSTER [CLUSTER...]",
		Aliases: []string{"remove"},
		Short:   "Remove one or more clusters",
		Long:    removeDescription,
		Example: removeExample,
		Args:    cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.rules = args
			return runRemove(storageosCli, &opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.force, "force", "f", false, "Force the removal of one or more clusters")
	return cmd
}

func runRemove(storageosCli *command.StorageOSCli, opt *removeOptions) error {

	client, err := discovery.NewClient("", "", "")
	if err != nil {
		return err
	}

	status := 0

	for _, ref := range opt.rules {

		if err := client.ClusterDelete(ref); err != nil {
			fmt.Fprintf(storageosCli.Err(), "%s\n", err)
			status = 1
			continue
		}
		fmt.Fprintf(storageosCli.Out(), "%s\n", ref)
	}

	if status != 0 {
		return cli.StatusError{StatusCode: status}
	}
	return nil
}

var removeDescription = `
Remove one or more clusters..
`

var removeExample = `
$ storageos cluster rm bd89e0ee-39a6-4790-9073-823181dbd69c
bd89e0ee-39a6-4790-9073-823181dbd69c
`
