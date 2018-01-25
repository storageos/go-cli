package rule

import (
	"fmt"

	"context"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/pkg/validation"
)

type removeOptions struct {
	force bool
	rules []string
}

func newRemoveCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt removeOptions

	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] RULE [RULE...]",
		Aliases: []string{"remove"},
		Short:   "Remove one or more rules",
		Long:    removeDescription,
		Example: removeExample,
		Args:    cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.rules = args
			return runRemove(storageosCli, &opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.force, "force", "f", false, "Force the removal of one or more rules")
	return cmd
}

func runRemove(storageosCli *command.StorageOSCli, opt *removeOptions) error {
	client := storageosCli.Client()
	status := 0

	for _, ref := range opt.rules {
		namespace, name, err := validation.ParseRefWithDefault(ref)
		if err != nil {
			fmt.Fprintf(storageosCli.Err(), "%s\n", err)
			status = 1
			continue
		}
		params := types.DeleteOptions{
			Name:      name,
			Namespace: namespace,
			Force:     opt.force,
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
