package pool

import (
	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/inspect"
)

type inspectOptions struct {
	format string
	names  []string
}

func newInspectCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] POOL [POOL...]",
		Short: "Display detailed information on one or more capacity pools",
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

	getFunc := func(name string) (interface{}, []byte, error) {
		i, err := client.Pool(name)
		return i, nil, err
	}

	return inspect.Inspect(storageosCli.Out(), opt.names, opt.format, getFunc)
}
