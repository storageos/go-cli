package user

import (
	"github.com/spf13/cobra"
	//storageos "github.com/storageos/go-api"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/inspect"
	"github.com/storageos/go-cli/pkg/constants"
)

type inspectOptions struct {
	format string
	users  []string
}

func newInspectCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] USER [USER...]",
		Short: "Display detailed information on one or more user(s)",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.users = args
			return runInspect(storageosCli, opt)
		},
	}

	cmd.Flags().StringVarP(&opt.format, "format", "f", "", "Format the output using the given Go template"+constants.UserTemplate)

	return cmd
}

func runInspect(storageosCli *command.StorageOSCli, opt inspectOptions) error {
	client := storageosCli.Client()

	getFunc := func(ref string) (interface{}, []byte, error) {
		i, err := client.User(ref)
		return i, nil, err
	}

	return inspect.Inspect(storageosCli.Out(), opt.users, opt.format, getFunc)
}
