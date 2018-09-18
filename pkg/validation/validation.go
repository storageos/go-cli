package validation

import (
	"fmt"
	"regexp"
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

var deprecatedLabels = map[string]string{
	"storageos.feature.replication":  "storageos.com/replication",
	"storageos.feature.replicas":     "storageos.com/replicas",
	"storageos.feature.nocompress":   "storageos.com/nocompress",
	"storageos.feature.nocache":      "storageos.com/nocache",
	"storageos.feature.nowritecache": "storageos.com/nowritecache",
	"storageos.feature.throttle":     "storageos.com/throttle",
	"storageos.hint.master":          "storageos.com/hint.master",
	"storageos.hint.docker":          "storageos.com/hint.docker",
	"storageos.driver":               "storageos.com/driver",
}

func labeldeprecationNotice(old, new string) string {
	depNotice := fmt.Sprintf("the label '%s' has been deprecated in favour of '%s'", old, new)
	return depNotice + ", refer to https://docs.storageos.com for usage details"
}

// GetDeprecations will check the provided labels for deprecated values.
// A list of deprecation notices are returned, if any.
func GetDeprecations(labels map[string]string) (notices []string) {

	for l := range labels {
		n, deprecated := IsDeprecated(l)
		if deprecated {
			notices = append(notices, n)
		}
	}

	return notices
}

// IsDeprecated will check if the label given has been deprecated,
// returning the appropriate deprecation notice if so.
func IsDeprecated(label string) (notice string, deprecated bool) {

	if updated, ok := deprecatedLabels[label]; ok {
		return labeldeprecationNotice(label, updated), true
	}

	return "", false
}
