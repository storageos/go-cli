package validation

import (
	"fmt"
	//"net"
	//"net/url"
	"regexp"
	//"strconv"
	"strings"

	storageos "github.com/storageos/go-api"
)

const hostnameFmt string = `^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`

var hostnameRegexp = regexp.MustCompile(hostnameFmt)

// ValidFSTypes lists the filesystem types that may be supported.
var ValidFSTypes = []string{"ext2", "ext3", "ext4", "xfs", "btrfs"}

// IsValidFSType tests that the argument is a valid IP address.
func IsValidFSType(value string) error {
	for _, t := range ValidFSTypes {
		if value == t {
			return nil
		}
	}
	return fmt.Errorf("fs type not valid, available types: %s", strings.Join(ValidFSTypes, ", "))
}

// ParseRefWithDefault wraps a call to the go-api's ParseRef
// function, but adds default if the namespace is not defined.
func ParseRefWithDefault(ref string) (string, string, error) {
	namespace, name, err := storageos.ParseRef(ref)
	if err != nil {
		return storageos.ParseRef("default/" + ref)
	}
	return namespace, name, err
}
