package mount

import (
	"fmt"
	"os"
	"path/filepath"
)

//go:generate stringer -type=FSType fs_types.go
type FSType int

const (
	ext2 FSType = iota
	ext3
	ext4
	xfs
	btrfs
)

func (f FSType) mountBin() string {
	return "/bin/mount"
}

func (f FSType) umountBin() string {
	return "/bin/umount"
}

func (f FSType) mkfsBin() string {
	switch f {
	case ext2:
		return "/sbin/mkfs.ext2"
	case ext3:
		return "/sbin/mkfs.ext3"
	case ext4:
		return "/sbin/mkfs.ext4"
	case xfs:
		return "/sbin/mkfs.xfs"
	case btrfs:
		return "/bin/mkfs.btrfs"
	default:
		return ""
	}
}

func (f FSType) checkPlatform() error {
	platformError := &MountError{fmt.Sprintf("filesystem (%v) not supported by host", f), true}

	// Check the mkfs binary exists and is executable
	fi, err := os.Lstat(f.mkfsBin())
	if err != nil || (fi.Mode()&0111) == 0 {
		return platformError
	}

	// Check the mount binary exists and is executable
	fi, err = os.Lstat(f.mountBin())
	if err != nil || (fi.Mode()&0111) == 0 {
		return platformError
	}

	// Check the umount binary exists and is executable
	fi, err = os.Lstat(f.umountBin())
	if err != nil || (fi.Mode()&0111) == 0 {
		return platformError
	}

	// If they all succeeded, we support the fs
	return nil

}

func (f FSType) UnmountCommand(mountpoint string) (command string, args []string, err error) {
	if err := f.checkPlatform(); err != nil {
		return "", nil, err
	}

	return f.umountBin(), []string{mountpoint}, nil
}

func (f FSType) MountCommand(path, mountpoint string) (command string, args []string, err error) {
	if err := f.checkPlatform(); err != nil {
		return "", nil, err
	}

	return f.mountBin(), []string{"-t", f.String(), path, mountpoint}, nil
}

func (f FSType) MkfsCommand(path string) (command string, args []string, err error) {
	if err := f.checkPlatform(); err != nil {
		return "", nil, err
	}

	args = []string{path}

	// We set some extra options on ext4 for additional performance and safety
	if f == ext4 {
		_, volID := filepath.Split(path)
		args = append([]string{"-F", "-U", volID, "-b", "4096", "-E", "lazy_itable_init=0,lazy_journal_init=0"}, args...)
	}

	return f.mkfsBin(), args, nil
}

// ParseFSType take the human readable fstype string and returns an FSType object related to the selected filesystem
func ParseFSType(s string) (FSType, error) {
	switch s {
	case ext2.String():
		return ext2, nil
	case ext3.String():
		return ext3, nil
	case ext4.String():
		return ext4, nil
	case xfs.String():
		return xfs, nil
	case btrfs.String():
		return btrfs, nil
	default:
		return 0, fmt.Errorf("unknown fs type (%v)", s)
	}
}
