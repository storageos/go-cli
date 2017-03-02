package namespace

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

type removeOptions struct {
	force      bool
	namespaces []string
}

func newRemoveCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opts removeOptions

	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] NAMESPACE [NAMESPACE...]",
		Aliases: []string{"remove"},
		Short:   "Remove one or more namespaces",
		Long:    removeDescription,
		Example: removeExample,
		Args:    cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.namespaces = args
			return runRemove(storageosCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.force, "force", "f", false, "Force the removal of one or more namespaces")
	return cmd
}

func runRemove(storageosCli *command.StorageOSCli, opts *removeOptions) error {
	client := storageosCli.Client()
	status := 0

	for _, name := range opts.namespaces {
		params := types.DeleteOptions{
			Name:    name,
			Force:   opts.force,
			Context: context.Background(),
		}

		if err := client.NamespaceDelete(params); err != nil {
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
Remove one or more namespaces. You cannot remove a namespace that is active.
`

var removeExample = `
$ storageos namespace rm testnamespace
testnamespace
`
