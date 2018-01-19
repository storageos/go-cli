package volume

import (
	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	cliconfig "github.com/storageos/go-cli/cli/config"
)

type mountOptions struct {
	ref        string
	mountpoint string // mountpoint
	fsType     string
}

func newMountCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt mountOptions

	cmd := &cobra.Command{
		Use:   "mount [OPTIONS] VOLUME MOUNTPOINT",
		Short: "Mount specified volume",
		Args:  cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.ref = args[0]
			opt.mountpoint = args[1]
			return runMount(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.fsType, "fsType", "m", cliconfig.DefaultFSType, `Volume fs type`)

	return cmd
}
