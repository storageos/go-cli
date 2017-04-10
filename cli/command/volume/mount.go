package volume

import (
	"fmt"
	"strings"

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

	fmt.Println("checking volume ", opt.ref, opt.mountpoint)

	errs := isVolumeReady(storageosCli, opt.ref)
	if len(errs) > 0 {
		return fmt.Errorf("cannot mount volume: %s", strings.Join(errs, ", "))
	}

	fmt.Println("volume is ready for mount, mounting..")

	return nil
}

// isVolumeReady - mount only unmounted and active volume
func isVolumeReady(storageosCli *command.StorageOSCli, ref string) (errs []string) {
	client := storageosCli.Client()

	namespace, name, err := storageos.ParseRef(ref)
	if err != nil {
		return []string{err.Error()}
	}
	vol, err := client.Volume(namespace, name)
	if err != nil {
		return []string{err.Error()}
	}

	if vol.Status != "active" {
		errs = append(errs, fmt.Sprintf("can only mount active volumes, current status: '%s'", vol.Status))
	}

	if vol.Mounted {
		errs = append(errs, "volume is mounted, unmount it before mounting it again")
	}

	return errs
}
