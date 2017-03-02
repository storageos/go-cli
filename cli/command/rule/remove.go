package rule

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

type removeOptions struct {
	force   bool
	rules []string
}

func newRemoveCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opts removeOptions

	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] RULE [RULE...]",
		Aliases: []string{"remove"},
		Short:   "Remove one or more rules",
		Long:    removeDescription,
		Example: removeExample,
		Args:    cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.rules = args
			return runRemove(storageosCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.force, "force", "f", false, "Force the removal of one or more rules")
	return cmd
}

func runRemove(storageosCli *command.StorageOSCli, opts *removeOptions) error {
	client := storageosCli.Client()
	status := 0

	for _, ref := range opts.rules {
		namespace, name, err := storageos.ParseRef(ref)
		if err != nil {
			fmt.Fprintf(storageosCli.Err(), "%s\n", err)
			status = 1
			continue
		}
		params := types.DeleteOptions{
			Name:      name,
			Namespace: namespace,
			Force:     opts.force,
			Context:   context.Background(),
		}

		if err := client.RuleDelete(params); err != nil {
			fmt.Fprintf(storageosCli.Err(), "%s\n", err)
			status = 1
			continue
		}
		fmt.Fprintf(storageosCli.Out(), "%s/%s\n", namespace, name)
	}

	if status != 0 {
		return cli.StatusError{StatusCode: status}
	}
	return nil
}

var removeDescription = `
Remove one or more rules. You cannot remove a rule that is in use by a container.
`

var removeExample = `
$ storageos rule rm default/testvol
testvol
`
