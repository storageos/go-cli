package mount

import (
	"fmt"
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

func (f FSType) UnmountCommand(mountpoint string) (command string, args []string) {
	const umount = "/bin/umount"
	return umount, []string{mountpoint}
}

func (f FSType) MountCommand(path, mountpoint string) (command string, args []string) {
	const mount = "/bin/mount"
	return mount, []string{"-t", f.String(), mountpoint}
}

func (f FSType) MkfsCommand(path string) (command string, args []string, err error) {
	const (
		mkfsExt2  = "/sbin/mkfs.ext2"
		mkfsExt3  = "/sbin/mkfs.ext3"
		mkfsExt4  = "/sbin/mkfs.ext4"
		mkfsXfs   = "/sbin/mkfs.xfs"
		mkfsBtrfs = "/bin/mkfs.btrfs"
	)

	switch f {
	case ext2:
		return mkfsExt2, []string{path}, nil

	case ext3:
		return mkfsExt3, []string{path}, nil

	case ext4:
		_, volID := filepath.Split(path)
		return mkfsExt4, []string{"-F", "-U", volID, "-b", "4096", "-E", "lazy_itable_init=0,lazy_journal_init=0", path}, nil

	case xfs:
		return mkfsXfs, []string{path}, nil

	case btrfs:
		return mkfsBtrfs, []string{path}, nil

	default:
		return "", nil, fmt.Errorf("unsuported filesystem (%s)", f)
	}
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
