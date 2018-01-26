package pool

import (
	"fmt"

	"context"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

type removeOptions struct {
	force bool
	pools []string
}

func newRemoveCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt removeOptions

	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] POOL [POOL...]",
		Aliases: []string{"remove"},
		Short:   "Remove one or more capacity pools",
		Long:    removeDescription,
		Example: removeExample,
		Args:    cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.pools = args
			return runRemove(storageosCli, &opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.force, "force", "f", false, "Force the removal of one or more pools")
	return cmd
}

func runRemove(storageosCli *command.StorageOSCli, opt *removeOptions) error {
	client := storageosCli.Client()
	status := 0

	for _, name := range opt.pools {
		params := types.DeleteOptions{
			Name:    name,
			Force:   opt.force,
			Context: context.Background(),
		}

		if err := client.PoolDelete(params); err != nil {
			fmt.Fprintf(storageosCli.Err(), "%s\n", err)
			status = 1
			continue
		}
		fmt.Fprintf(storageosCli.Out(), "%s\n", name)
	}

	if status != 0 {
		return cli.StatusError{StatusCode: status}
	}
	return nil
}

var removeDescription = `
Remove one or more capacity pools. You cannot remove a pool that is active.
`

var removeExample = `
$ storageos pool rm testpool
testpool
`
