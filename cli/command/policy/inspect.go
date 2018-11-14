package policy

import (
	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/inspect"
)

type inspectOptions struct {
	format   string
	policies []string
}

func newInspectCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] POLICY [POLICY...]",
		Short: "Display detailed information on one or more polic(y|ies)",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.policies = args
			return runInspect(storageosCli, opt)
		},
	}

	cmd.Flags().StringVarP(&opt.format, "format", "f", "", "Format the output using the given Go template (type --format -h or --help for a detail usage)")

	return cmd
}

func runInspect(storageosCli *command.StorageOSCli, opt inspectOptions) error {
	client := storageosCli.Client()

	getFunc := func(ref string) (interface{}, []byte, error) {
		i, err := client.Policy(ref)
		return i, nil, err
	}

	return inspect.Inspect(storageosCli.Out(), opt.policies, opt.format, getFunc)
}
