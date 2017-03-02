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
	var opts inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] POOL [POOL...]",
		Short: "Display detailed information on one or more capacity pools",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.names = args
			return runInspect(storageosCli, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.format, "format", "f", "", "Format the output using the given Go template")

	return cmd
}

func runInspect(storageosCli *command.StorageOSCli, opts inspectOptions) error {
	client := storageosCli.Client()

	getFunc := func(name string) (interface{}, []byte, error) {
		i, err := client.Pool(name)
		return i, nil, err
	}

	return inspect.Inspect(storageosCli.Out(), opts.names, opts.format, getFunc)
}
