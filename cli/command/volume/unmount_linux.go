package volume

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	storageos "github.com/storageos/go-api"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli/command"
	cliconfig "github.com/storageos/go-cli/cli/config"
	"github.com/storageos/go-cli/pkg/host"
	"github.com/storageos/go-cli/pkg/mount"
	"github.com/storageos/go-cli/pkg/system"
)

func runUnmount(storageosCli *command.StorageOSCli, opt unmountOptions) error {
	mountDriver := mount.New(cliconfig.DeviceRootPath)

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
