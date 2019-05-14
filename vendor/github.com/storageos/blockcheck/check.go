// Package blockcheck provides the ability to verify a block device has not been
// formatted with a filesystem.
package blockcheck

import (
	"errors"
	"io"
	"os"
	"syscall"
)

const (
	kilobyte = 1024

	// peekByteCount defines how far into the block device to read when
	// verifying the block device is empty.
	peekByteCount = 256 * kilobyte

	// The size of the read() syscalls.
	chunkSize = 4 * kilobyte
)

var (
	// ErrNotBlockDevice is returned when calling IsBlockDeviceEmpty with a path
	// to anything other than a block device.
	ErrNotBlockDevice = errors.New("not a block device")
)

// IsBlockDeviceEmpty returns true if path is the path to an empty block device.
//
// "Empty" is defined as "all zeros at the start of the device", which is
// suitable for detecting a valid filesystem.
//
// If the block device contains data, or an error occurs, IsBlockDeviceEmpty
// returns false.
func IsBlockDeviceEmpty(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	// Ensure the path is to a block device.
	if info.Mode()&os.ModeDevice == 0 {
		return false, ErrNotBlockDevice
	}

	var flags int
	flags |= os.O_RDONLY       // Read only
	flags |= syscall.O_CLOEXEC // Do not share the fd after a fork()

	// Open the block device
	f, err := os.OpenFile(path, flags, 0)
	if err != nil {
		return false, err
	}
	defer f.Close()

	return isEmpty(f)
}

// isEmpty returns true if the io.Reader passed contains no data in the first
// peekByteCount bytes.
//
// Reads are performed chunkSize bytes at a time.
func isEmpty(f io.Reader) (bool, error) {
	buf := make([]byte, chunkSize)
	for i := 0; i < peekByteCount; i += chunkSize {
		n, err := f.Read(buf)

		if n > 0 && !isZeros(buf[:n]) {
			return false, nil
		}

		switch err {
		case nil:
			continue
		case io.EOF:
			return true, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// isZeros returns true if all bytes in b are 0
//
// Unfortunately this isn't unrolled or replaced with SIMD by the compiler
// (1.11.4).
func isZeros(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if b[i] != 0 {
			return false
		}
	}
	return true
}
