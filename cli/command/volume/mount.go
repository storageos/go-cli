package volume

import (
	"fmt"

	"github.com/dnephin/cobra"
	storageos "github.com/storageos/go-api"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	cliconfig "github.com/storageos/go-cli/cli/config"
	"github.com/storageos/go-cli/pkg/system"
)

type mountOptions struct {
	ref        string
	mountpoint string // mountpoint
}

func newMountCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt mountOptions

	cmd := &cobra.Command{
		Use:   "mount [OPTIONS] VOLUME",
		Short: "Mount specified volume",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.ref = args[0]
			return runMount(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.mountpoint, "mountpoint", "m", "", `Volume mountpoint`)

	return cmd
}

func runMount(storageosCli *command.StorageOSCli, opt mountOptions) error {

	// checking whether we are on storageos node
	_, err := system.Stat(cliconfig.DeviceRootPath)
	if err != nil {
		return fmt.Errorf("device root path '%s' not found, check whether StorageOS is running", cliconfig.DeviceRootPath)
	}

	fmt.Println("mounting volume ", opt.ref, opt.mountpoint)

	client := storageosCli.Client()

	_ = func(ref string) (interface{}, []byte, error) {
		namespace, name, err := storageos.ParseRef(ref)
		if err != nil {
			return nil, nil, err
		}
		i, err := client.Volume(namespace, name)
		return i, nil, err
	}

	return nil
}
