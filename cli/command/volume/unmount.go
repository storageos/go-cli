package volume

import (
	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

type unmountOptions struct {
	ref   string
	force bool
}

func newUnmountCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt unmountOptions

	cmd := &cobra.Command{
		Use:   "unmount [OPTIONS] VOLUME",
		Short: "Unmount specified volume",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.ref = args[0]
			return runUnmount(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.force, "force", "f", false, `Force unmount`)

	return cmd
}
