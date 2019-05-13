package mount

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/storageos/blockcheck"
)

// defaultTimeOut is the maximum time to wait for async volume operations.
const defaultTimeOut = 45

// initRawVolume makes sure there is a raw volume at the path provided, and that
// it's ready for use by Docker.  This includes creating a filesystem if there's
// not one already present.
func initRawVolume(ctx context.Context, path string, fsType FSType) error {

	if err := waitForVolume(ctx, path); err != nil {
		return err
	}
	log.Debugf("volume found: %s", path)

	ft, err := getVolumeFSType(ctx, path)
	log.Debugf("volume %s has fs type: %s", path, ft)
	if err != nil {
		return err
	}

	if ft == "raw" {
		log.Debugf("creating %s filesystem on volume %s", fsType, path)
		if err := createFilesystem(ctx, fsType, path); err != nil {
			return err
		}
		log.Infof("%s filesystem created on volume %s", fsType, path)
	}

	return nil
}

// waitForVolume waits for the StorageOS raw volume to be created.  The retry
// timeout doubles on every invocation.
func waitForVolume(ctx context.Context, path string) error {

	var retries int
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("deadline exceeded while waiting for volume")
		default:
			timeOff := backoff(retries)
			if volumeExists(path) {
				if volumeWritable(path) {
					// volume ready, exit
					return nil
				}
				log.Debugf("waiting for volume to come online: %s, retrying in %v", path, timeOff)
			} else {
				log.Debugf("waiting for volume: %s, retrying in %v", path, timeOff)
			}

			if abort(start, timeOff) {
				return fmt.Errorf("timed out waiting for volume to be ready")
			}

			retries++
			time.Sleep(timeOff)
		}
	}
}

// getVolumeFSType returns the volumes filesystem type if supported (currently
// only ext4), "raw" if none, or an error if unsupported.  This may need
// alteration for different OS flavours.
func getVolumeFSType(ctx context.Context, path string) (string, error) {

	// Run the file command against the path
	// NOTE: -s returns same info for file and block devices. *DO NOT REMOVE*
	out, err := runFile(ctx, "-s", path)
	if err != nil {
		return "", err
	}

	log.Debugf("checking volume for existing filesystem: %s: output: %s", path, out)
	return parseFileOutput(path, out)
}

// parseFileOutput ttries to determine the volume's filesystem type or if it is
// "raw" and can be formatted safely.
func parseFileOutput(path string, out string) (string, error) {
	switch {
	case out == path+": data", out == path+": empty":
		// At least on Ubuntu, file will report "/path/to/file: data" for a raw
		// volume without a filesystem.  If this matches, we expect to be able to
		// reformat. There seems to be a similar case for empty, we make the same
		// assumptions for this too.
		return "raw", nil
	case strings.HasPrefix(out, path+": block special"):
		// block special devices are reported on block devices when `file -s` is not
		// used.  Leaving this check in "just in case", as we should never reformat
		// a based on this info as it hides whether there is a filesystem.
		return "", fmt.Errorf("detected block special device, aborting")
	case strings.HasPrefix(out, path+": Linux rev 1.0 ext2 filesystem data"):
		return "ext2", nil
	case strings.HasPrefix(out, path+": Linux rev 1.0 ext3 filesystem data"):
		return "ext3", nil
	case strings.HasPrefix(out, path+": Linux rev 1.0 ext4 filesystem data"):
		return "ext4", nil
	case strings.HasPrefix(out, path+": SGI XFS filesystem data"):
		return "xfs", nil
	case strings.HasPrefix(out, path+": BTRFS Filesystem"):
		return "btrfs", nil
	case strings.HasPrefix(out, path+": DOS/MBR boot sector, code offset 0x58+2, OEM-ID \"mkfs.fat\""):
		return "fat", nil
	case strings.HasPrefix(out, path+": DOS/MBR boot sector, code offset 0x52+2, OEM-ID \"NTFS"):
		return "ntfs", nil
	}
	return "", fmt.Errorf("unknown fs type: %s", out)
}

func createFilesystem(ctx context.Context, fstype FSType, path string) error {
	// At this point, the block device is about to be formatted after running
	// "/usr/bin/file" on the device which has indicated there is no filesystem.
	//
	// As a failsafe, check that path is a path to a valid block device that
	// contains no data.
	empty, err := blockcheck.IsBlockDeviceEmpty(path)
	if err != nil {
		// An error occurred when trying to open and read the block device.
		//
		// As the state of the device cannot be determined, take the safest path
		// and stop.
		log.WithError(err).Error("unable to read device")
		return err
	}

	// Ensure the device is empty before trying to do anything to it.
	if !empty {
		// The device contains data. Do not attempt to format the volume.
		//
		// Stop, returning the original mount error to the caller.
		log.Warn("device contains data, aborting mkfs")
		return errors.New("device contains data, aborting mkfs")
	}

	var retries int
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("deadline exceeded while trying to create filesystem")
		default:
			timeOff := backoff(retries)
			err := runMkfs(ctx, fstype, path)
			if err == nil {
				// filesystem created, exit
				return nil
			}

			// Bail early if this is a fatal error
			if mountErr, ok := err.(*Error); ok && mountErr.Fatal {
				return err
			}

			log.WithFields(log.Fields{
				"fstype": fstype,
				"path":   path,
				"err":    err.Error(),
			}).Warnf("create filesystem failed, retrying in %v", timeOff)

			if abort(start, timeOff) {
				return fmt.Errorf("failed to create filesystem")
			}

			retries++
			time.Sleep(timeOff)
		}
	}
}

func runMkfs(ctx context.Context, fstype FSType, path string) error {
	bin, args, err := fstype.MkfsCommand(path)
	if err != nil {
		return err
	}

	out, err := runCmd(ctx, bin, args...)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"fstype": fstype,
		"path":   path,
		"output": out,
	}).Debug("created filesystem")

	return nil
}

// getVolumeIDFromPath returns the volume id from the path name.
func getVolumeIDFromPath(path string) string {
	_, file := filepath.Split(path)
	return file
}

// volumeExists returns true if the volume exists.
func volumeExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// volumeWritable returns true if the volume is writable (ready).
func volumeWritable(path string) bool {
	file, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return false
	}
	return true
}

// backoff is used to calculate the retry delay. It doubles wait time on each
// retry.
func backoff(retries int) time.Duration {
	b, max := 1, defaultTimeOut
	for b < max && retries > 0 {
		b *= 2
		retries--
	}
	if b > max {
		b = max
	}
	return time.Duration(b) * time.Second
}

// abort returns true if the timeout has been exceeded
func abort(start time.Time, timeOff time.Duration) bool {
	return timeOff+time.Since(start) >= time.Duration(defaultTimeOut)*time.Second
}
