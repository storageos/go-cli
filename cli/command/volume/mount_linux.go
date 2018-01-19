package volume

import (
	"context"
	"fmt"
	"strings"
	"time"

	storageos "github.com/storageos/go-api"
	"github.com/storageos/go-cli/cli/command"
	cliconfig "github.com/storageos/go-cli/cli/config"
	"github.com/storageos/go-cli/pkg/host"
	"github.com/storageos/go-cli/pkg/mount"
	"github.com/storageos/go-cli/pkg/system"
	"github.com/storageos/go-cli/pkg/validation"

	"github.com/storageos/go-api/types"

	log "github.com/Sirupsen/logrus"
)

func runMount(storageosCli *command.StorageOSCli, opt mountOptions) error {
	// checking whether we are on storageos node
	_, err := system.Stat(cliconfig.DeviceRootPath)
	if err != nil {
		return fmt.Errorf("device root path '%s' not found, check whether StorageOS is running", cliconfig.DeviceRootPath)
	}

	// validating fsType
	err = validation.IsValidFSType(opt.fsType)
	if err != nil {
		return err
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

	// checking readiness
	errs := isVolumeReady(vol, name)
	if len(errs) > 0 {
		return fmt.Errorf("cannot mount volume: %s", strings.Join(errs, ", "))
	}

	var hostname string

	// getting current hostname
	hostname, err = host.Get()
	if err != nil {
		hostname = "unknown"
	}

	err = client.VolumeMount(types.VolumeMountOptions{
		ID: vol.ID, Namespace: namespace,
		Client:     hostname,
		Mountpoint: opt.mountpoint,
		FsType:     opt.fsType,
	})
	if err != nil {
		return err
	}

	err = retryableMount(vol, opt.mountpoint, opt.fsType)
	if err != nil {
		log.WithFields(log.Fields{
			"namespace":  namespace,
			"volumeName": name,
			"error":      err,
		}).Error("error while mounting volume")
		// should unmount volume in the CP if we failed here
		newErr := client.VolumeUnmount(types.VolumeUnmountOptions{ID: vol.ID, Namespace: namespace})
		if newErr != nil {
			log.WithFields(log.Fields{
				"volumeId": vol.ID,
				"err":      newErr,
			}).Error("failed to unmount volume")
		}

		return err
	}

	fmt.Printf("volume %s mounted: %s\n", vol.Name, opt.mountpoint)

	return nil
}

func retryableMount(volume *types.Volume, mountpoint, fsType string) error {

	driver := mount.New(cliconfig.DeviceRootPath)

	// unmount can take some time
	maxRetries := 10
	retries := 0

RETRY:
	// Perform the mount
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err := driver.MountVolume(ctx, volume.ID, mountpoint, fsType)
	if err != nil {
		log.WithFields(log.Fields{
			"volume_id":  volume.ID,
			"mountpoint": mountpoint,
			"err":        err.Error(),
		}).Error(" failed to mount volume")

		if retries < maxRetries {
			time.Sleep(250 * time.Millisecond)
			retries++
			goto RETRY
		}

		return err
	}

	return nil
}

// isVolumeReady - mount only unmounted and active volume
func isVolumeReady(vol *types.Volume, ref string) (errs []string) {

	if vol.Status != "active" {
		errs = append(errs, fmt.Sprintf("can only mount active volumes, current status: '%s'", vol.Status))
	}

	if vol.Mounted {
		errs = append(errs, "volume is mounted, unmount it before mounting it again")
	}

	return errs
}
