package mount

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
)

// FilePermissions is an interface which allows a struct to set
// ownership and permissions easily on a file it describes.
type FilePermissions interface {
	// User returns a user ID or user name
	User() string

	// Group returns a group ID. Group names are not supported.
	Group() string

	// Mode returns a string of file mode bits e.g. "0644"
	Mode() string
}

// SetFilePermissions handles configuring ownership and permissions settings
// on a given file. It takes a path and any struct implementing the
// FilePermissions interface. All permission/ownership settings are optional.
// If no user or group is specified, the current user/group will be used. Mode
// is optional, and has no default (the operation is not performed if absent).
// User may be specified by name or ID, but group may only be specified by ID.
func SetFilePermissions(path string, p FilePermissions) error {
	var err error
	uid, gid := os.Getuid(), os.Getgid()

	if p.User() != "" {
		if uid, err = strconv.Atoi(p.User()); err == nil {
			goto GROUP
		}

		// Try looking up the user by name
		if u, err := user.Lookup(p.User()); err == nil {
			uid, _ = strconv.Atoi(u.Uid)
			goto GROUP
		}

		return fmt.Errorf("invalid user specified: %v", p.User())
	}

GROUP:
	if p.Group() != "" {
		if gid, err = strconv.Atoi(p.Group()); err != nil {
			return fmt.Errorf("invalid group specified: %v", p.Group())
		}
	}
	if err := os.Chown(path, uid, gid); err != nil {
		return fmt.Errorf("failed setting ownership to %d:%d on %q: %s",
			uid, gid, path, err)
	}

	if p.Mode() != "" {
		mode, err := strconv.ParseUint(p.Mode(), 8, 32)
		if err != nil {
			return fmt.Errorf("invalid mode specified: %v", p.Mode())
		}
		if err := os.Chmod(path, os.FileMode(mode)); err != nil {
			return fmt.Errorf("failed setting permissions to %d on %q: %s",
				mode, path, err)
		}
	}

	return nil
}
