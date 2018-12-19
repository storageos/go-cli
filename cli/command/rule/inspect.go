package rule

import (
	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/inspect"
	"github.com/storageos/go-cli/pkg/validation"
)

type inspectOptions struct {
	format string
	names  []string
}

func newInspectCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] RULE [RULE...]",
		Short: "Display detailed information on one or more rules",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.names = args
			return runInspect(storageosCli, opt)
		},
	}

	cmd.Flags().StringVarP(&opt.format, "format", "f", "", "Format the output using a custom template (try \"help\" for more info)")

	return cmd
}

func runInspect(storageosCli *command.StorageOSCli, opt inspectOptions) error {
	client := storageosCli.Client()

	getFunc := func(ref string) (interface{}, []byte, error) {
		namespace, name, err := validation.ParseRefWithDefault(ref)
		if err != nil {
			return nil, nil, err
		}
		i, err := client.Rule(namespace, name)
		return i, nil, err
	}

	return inspect.Inspect(storageosCli.Out(), opt.names, opt.format, getFunc)
}
