// +build linux

package volume

import (
	"context"
	"fmt"
	// "strings"
	"time"

	"github.com/dnephin/cobra"
	storageos "github.com/storageos/go-api"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	cliconfig "github.com/storageos/go-cli/cli/config"
	// "github.com/storageos/go-cli/pkg/host"
	"github.com/storageos/go-cli/pkg/mount"
	"github.com/storageos/go-cli/pkg/system"
	// "github.com/storageos/go-cli/pkg/validation"

	"github.com/storageos/go-api/types"

	log "github.com/sirupsen/logrus"
	"github.com/storageos/go-cli/pkg/host"
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
			return runUnmount(storageosCli, opt, mount.New(cliconfig.DeviceRootPath))
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.force, "force", "f", false, `Force unmount`)

	return cmd
}

func runUnmount(storageosCli *command.StorageOSCli, opt unmountOptions, mountDriver mount.Driver) error {

	// checking whether we are on storageos node
	_, err := system.Stat(cliconfig.DeviceRootPath)
	if err != nil {
		return fmt.Errorf("device root path '%s' not found, check whether StorageOS is running", cliconfig.DeviceRootPath)
	}

	client := storageosCli.Client()
	namespace, name, err := storageos.ParseRef(opt.ref)
	if err != nil {
		return err
	}

	vol, err := client.Volume(namespace, name)
	if err != nil {
		return err
	}

	// getting volume
	hostname, err := host.Get()
	if err != nil && !opt.force {
		return fmt.Errorf("failed to get current node hostname, unable to unmount volume (must be forced), error: %s", err)
	}

	if hostname != vol.MountedBy && !opt.force {
		return fmt.Errorf("current hostname '%s' doesn't match volume's hostname '%s', unable to unmount volume (must be forced)", hostname, vol.MountedBy)
	}

	// unmounting it
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	err = mountDriver.UnmountVolume(ctx, vol.Mountpoint)
	if err != nil && !opt.force {
		return fmt.Errorf("unable to unmount volume (must be forced), error: %s", err)
	}

	err = client.VolumeUnmount(types.VolumeUnmountOptions{ID: vol.ID, Namespace: namespace})
	if err != nil {
		log.WithFields(log.Fields{
			"volumeId":  vol.ID,
			"namespace": namespace,
			"err":       err,
		}).Error("failed to unmount volume")
		return fmt.Errorf("unable to unmount volume, error: %s", err)
	}

	fmt.Printf("volume %s unmounted: %s\n", vol.Name, vol.Mountpoint)
	return nil
}
