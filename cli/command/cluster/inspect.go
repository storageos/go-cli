package cluster

import (
	"github.com/dnephin/cobra"
	// storageos "github.com/storageos/go-api"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/inspect"

	"github.com/storageos/go-cli/discovery"
)

type inspectOptions struct {
	format string
	names  []string
}

func newInspectCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] CLUSTER [CLUSTER...]",
		Short: "Display detailed information on one or more cluster",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.names = args
			return runInspect(storageosCli, opt)
		},
	}

	cmd.Flags().StringVarP(&opt.format, "format", "f", "", "Format the output using the given Go template")

	return cmd
}

func runInspect(storageosCli *command.StorageOSCli, opt inspectOptions) error {
	// client := storageosCli.Client()

	client, err := discovery.NewClient("", "", "")
	if err != nil {
		return err
	}

	getFunc := func(ref string) (interface{}, []byte, error) {
		i, err := client.ClusterStatus(ref)
		return i, nil, err
	}

	return inspect.Inspect(storageosCli.Out(), opt.names, opt.format, getFunc)
}
