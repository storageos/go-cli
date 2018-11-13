// +build linux

package volume

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/docker/docker/pkg/system"

	api "github.com/storageos/go-api"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	cliconfig "github.com/storageos/go-cli/cli/config"
	"github.com/storageos/go-cli/pkg/mount"
	"github.com/storageos/go-cli/pkg/validation"

	"github.com/storageos/go-api/types"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrVolumeMounted should be returned if it the volume has a mount lock.
	ErrVolumeMounted = errors.New("volume is mounted, unmount it before mounting it again")
)

// Failed mount retry back-off upper bound.
const mountBackoffUB = 4 * time.Second

type mountOptions struct {
	ref        string
	timeout    time.Duration
	mountpoint string // mountpoint
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
	flags.DurationVarP(&opt.timeout, "timeout", "t", 20*time.Second, "Mount action timeout period")

	return cmd
}

func runMount(storageosCli *command.StorageOSCli, opt mountOptions) error {
	// Get current hostname.
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	client := storageosCli.Client()

	// Check whether we are on storageos node.
	node, err := client.Node(hostname)
	if err != nil {
		if err == api.ErrNoSuchNode {
			return fmt.Errorf("cannot mount volume: current host is not a registered storageos cluster node")
		}
		return fmt.Errorf("failed to check if this host is a storageos cluster node: %#v", err)
	}

	// Check whether device dir exists.
	_, err = system.Stat(node.DeviceDir)
	if err != nil {
		return fmt.Errorf("device root path %q not found, check whether StorageOS is running", cliconfig.DeviceRootPath)
	}

	// must be root
	if euid := syscall.Geteuid(); euid != 0 {
		return fmt.Errorf("volume mount requires root permission - try prefixing command with `sudo -E`")
	}

	namespace, name, err := validation.ParseRefWithDefault(opt.ref)
	if err != nil {
		return err
	}

	vol, err := client.Volume(namespace, name)
	if err != nil {
		return err
	}

	// checking readiness
	if err := isVolumeReady(vol); err != nil {
		return fmt.Errorf("cannot mount volume: %v", err)
	}

	// Create a context which times out the mount
	ctx, cancel := context.WithTimeout(context.Background(), opt.timeout)
	defer cancel()

	err = client.VolumeMount(types.VolumeMountOptions{
		Context:    ctx,
		ID:         vol.ID,
		Namespace:  namespace,
		Client:     hostname,
		Mountpoint: opt.mountpoint,
		FsType:     vol.FSType,
	})
	if err != nil {
		select {
		case <-ctx.Done():
			// Server should handle what to do if we close the connection whilst waiting for a response
			return fmt.Errorf("timed out waiting for volume mount lock")
		default:
			return err
		}
	}

	fst, err := mount.ParseFSType(vol.FSType)
	if err != nil {
		return err
	}

	err = retryableMount(ctx, vol, node.DeviceDir, opt, fst)
	if err != nil {
		select {
		case <-ctx.Done():
			client.VolumeUnmount(types.VolumeUnmountOptions{ID: vol.ID, Namespace: namespace})
			return fmt.Errorf("timed out performing volume mount")
		default:
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
			return fmt.Errorf("failed to mount volume: %v", err)
		}
	}

	fmt.Printf("volume %s mounted: %s\n", vol.Name, opt.mountpoint)

	return nil
}

func retryableMount(ctx context.Context, volume *types.Volume, deviceRootDir string, opts mountOptions, fsType mount.FSType) error {
	driver := mount.New(deviceRootDir)
	// Limit the time which can be spent retrying
	retries := 0
	backoff := 250 * time.Millisecond

RETRY:
	err := driver.MountVolume(ctx, volume.ID, opts.mountpoint, fsType, volume.MkfsDoneAt.IsZero() && !volume.MkfsDone)
	if err == nil {
		return nil
	}

	// If this is a permanent error, stop retrying
	if mountErr, ok := err.(*mount.Error); ok && mountErr.Fatal {
		log.WithFields(log.Fields{
			"volume_id":  volume.ID,
			"mountpoint": opts.mountpoint,
			"err":        err.Error(),
		}).Error("failed to mount volume")
		return err
	}

	select {
	case <-ctx.Done():
		return err
	case <-time.After(backoff):
		// Increase backoff period
		retries++
		if backoff < mountBackoffUB {
			backoff *= 2
		}
		fmt.Printf("failed to mount volume, beginning retry %d\n", retries)
		log.WithFields(log.Fields{
			"volume_id":  volume.ID,
			"mountpoint": opts.mountpoint,
			"err":        err.Error(),
			"retry":      retries,
		}).Error("failed to mount volume, beginning retry")
		goto RETRY
	}
}

// isVolumeReady - mount only unmounted and active volume
func isVolumeReady(vol *types.Volume) error {

	if vol.Status != "active" {
		return fmt.Errorf("can only mount active volumes, current status: '%s'", vol.Status)
	}

	if vol.Mounted {
		return ErrVolumeMounted
	}

	return nil
}
