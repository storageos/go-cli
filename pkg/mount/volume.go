package mount

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// defaultTimeOut is the maximum time to wait for async volume operations.
const defaultTimeOut = 45

// initRawVolume makes sure there is a raw volume at the path provided, and that
// it's ready for use by Docker.  This includes creating a filesystem if there's
// not one already present.
func initRawVolume(ctx context.Context, path string, fsType string) error {

	if err := waitForVolume(path); err != nil {
		return err
	}
	log.Debugf("volume found: %s", path)

	ft, err := getVolumeFSType(ctx, path)
	log.Debugf("volume %s has fs type: %s", path, ft)
	if err != nil {
		return err
	}

	if ft == "raw" {
		_, err := createFilesystem(ctx, fsType, path, "")
		if err != nil {
			return err
		}
		log.Infof("%s filesystem created on volume %s", fsType, path)
	}

	return nil
}

// waitForVolume waits for the StorageOS raw volume to be created.  The retry
// timeout doubles on every invocation.
func waitForVolume(path string) error {

	var retries int
	start := time.Now()

	for {
		timeOff := backoff(retries)
		if volumeExists(path) {
			if volumeWritable(path) {
				// volume ready, exit
				return nil
			}
			log.Infof("waiting for volume to come online: %s, retrying in %v", path, timeOff)
		} else {
			log.Infof("waiting for volume: %s, retrying in %v", path, timeOff)
		}

		if abort(start, timeOff) {
			return fmt.Errorf("timed out waiting for volume to be ready")
		}

		retries++
		time.Sleep(timeOff)
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
	case out == path+": data":
		// At least on Ubuntu, file will report "/path/to/file: data" for a raw
		// volume without a filesystem.  If this matches, we expect to be able to
		// reformat.
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

func createFilesystem(ctx context.Context, fstype string, path string, options string) (string, error) {

	var out string
	var err error

	// Run mkfs
	switch fstype {
	case "ext2":
		out, err = runCmd(ctx, mkfsExt2, path)
		if err != nil {
			log.Warnf("mkfs output: %s", err.Error())
		}
	case "ext3":
		out, err = runCmd(ctx, mkfsExt3, path)
		if err != nil {
			log.Warnf("mkfs output: %s", err.Error())
		}
	case "ext4":
		// Get the volume id from the path
		id := getVolumeIDFromPath(path)
		out, err = runCmd(ctx, mkfsExt4, "-F", "-U", id, "-b", "4096", "-E", "lazy_itable_init=1,lazy_journal_init=1", path)
		if err != nil {
			log.Warnf("mkfs output: %s", err.Error())
		}
	case "xfs":
		out, err = runCmd(ctx, mkfsXfs, path)
		if err != nil {
			log.Warnf("mkfs output: %s", err.Error())
		}
	case "btrfs":
		out, err = runCmd(ctx, mkfsBtrfs, path)
		if err != nil {
			log.Warnf("mkfs output: %s", err.Error())
		}
	case "":
		return "", fmt.Errorf("filesystem not specified")
	default:
		return "", fmt.Errorf("unsupported filesystem: %s", fstype)
	}
	return out, nil
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
