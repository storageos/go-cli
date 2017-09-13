package mount

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

// Driver - generic mount driver interface
type Driver interface {
	MountVolume(ctx context.Context, id, mountpoint, fsType string) error
	UnmountVolume(ctx context.Context, mountpoint string) error
}

// DefaultDriver - default mount driver
type DefaultDriver struct {
	deviceRootPath string
}

// New - creates new instance of default driver
func New(deviceRootPath string) *DefaultDriver {
	return &DefaultDriver{
		deviceRootPath: deviceRootPath,
	}
}

// MountVolume - mounts specified volume
func (d *DefaultDriver) MountVolume(ctx context.Context, id, mountpoint, fsType string) error {
	return mountVolume(ctx, d.deviceRootPath, id, mountpoint, fsType)
}

// UnmountVolume - unmounts specified mountpoint
func (d *DefaultDriver) UnmountVolume(ctx context.Context, mountpoint string) error {
	return unmountVolume(ctx, mountpoint)
}

// deviceRootPath is the location of the StorageOS raw volumes.
// const deviceRootPath = constants.DeviceRootPath

// mountpointPerms will be used to set the filesystem permissions on the
// mountpoint.  Only Docker (running as root) needs to see the directory.
const mountpointPerms os.FileMode = 0700

// MountVolume mounts a StorageOS-based filesystem for use by Docker.
//
// It checks the volume first, waiting 30 seconds for it to be created, and
// creates an ext4 filesystem on it if there isn't already a filesystem.  The
// mount will fail if the mount command can't determine the fstype.
func mountVolume(ctx context.Context, deviceRootPath string, id string, mp string, fsType string) error {
	if err := initRawVolume(ctx, deviceRootPath+"/"+id, fsType); err != nil {
		log.WithFields(log.Fields{
			"id":      id,
			"fs_type": fsType,
			"err":     err.Error(),
		}).Error("volume init error")
		return err
	}
	log.Debugf("StorageOS volume ready: %s ", mp)

	if err := createMountPoint(mp); err != nil {
		return err
	}
	log.Debugf("Mountpoint created: %s ", mp)

	_, err := runMount(ctx, deviceRootPath+"/"+id, mp)
	if err != nil {
		log.WithFields(log.Fields{
			"path":        deviceRootPath + "/" + id,
			"mount_point": mp,
			"fs_type":     fsType,
			"error":       err,
		}).Error("Mount failed")
		return err
	}
	log.Debugf("Mounted volume: %s %s", deviceRootPath+"/"+id, mp)

	return nil
}

// unmountVolume unmounts a StorageOS-based filesystem and removes the
// mountpoint.
func unmountVolume(ctx context.Context, mp string) error {

	_, err := runUmount(ctx, mp)
	if err != nil {
		log.Errorf("Unmount failed: %s (%s)", mp, err)
		return err
	}
	log.Debugf("Unmounted volume: %s", mp)

	return nil
}

// createMountPoint creates a the target mountpoint on the filesystem given the
// path.
func createMountPoint(path string) error {
	return os.MkdirAll(path, mountpointPerms)
}
