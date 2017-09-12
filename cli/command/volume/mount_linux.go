// +build linux

package volume

import (
	"context"
	"fmt"
	"strings"
	"syscall"
	"time"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	cliconfig "github.com/storageos/go-cli/cli/config"
	"github.com/storageos/go-cli/pkg/host"
	"github.com/storageos/go-cli/pkg/mount"
	"github.com/storageos/go-cli/pkg/system"
	"github.com/storageos/go-cli/pkg/validation"

	"github.com/storageos/go-api/types"

	log "github.com/sirupsen/logrus"
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

func runMount(storageosCli *command.StorageOSCli, opt mountOptions) error {

	// checking whether we are on storageos node
	_, err := system.Stat(cliconfig.DeviceRootPath)
	if err != nil {
		return fmt.Errorf("device root path %q not found, check whether StorageOS is running", cliconfig.DeviceRootPath)
	}

	// must be root
	if euid := syscall.Geteuid(); euid != 0 {
		return fmt.Errorf("volume mount requires root permission - try prefixing command with `sudo`")
	}

	// validating fsType
	err = validation.IsValidFSType(opt.fsType)
	if err != nil {
		return err
	}

	client := storageosCli.Client()
	namespace, name, err := validation.ParseRefWithDefault(opt.ref)
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
