package validation

import (
	"errors"
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

var depricatedLabels = map[string]string{
	"storageos.feature.replication":   "storageos.com/replication",
	"storageos.feature.deduplication": "storageos.com/deduplication",
	"storageos.feature.replicas":      "storageos.com/replicas",
	"storageos.feature.nocompress":    "storageos.com/compression",
	"storageos.feature.nocache":       "storageos.com/nocache", // TODO: is this one correct?
	"storageos.feature.throttle":      "storageos.com/throttle",
	"storageos.hint.master":           "storageos.com/hint.master",
	"storageos.hint.docker":           "storageos.com/hint.docker",
	"storageos.driver":                "storageos.com/driver",
}

func labeldeprecationNotice(old, new string) string {
	depNotice := fmt.Sprintf("the label '%s' has been deprecated in favour of '%s'", old, new)
	return depNotice + ", refer to https://docs.storageos.com for usage details"
}

func ValidateLabelSet(labels map[string]string) (warnings []string, err error) {
	errs := make([]string, 0, len(labels))

	for k, v := range labels {
		w, e := ValidateLabel(k, v)
		warnings = append(warnings, w...)
		if err != nil {
			errs = append(errs, e.Error())
		}
	}

	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, ","))
	}

	return warnings, err
}

func ValidateLabel(k, v string) (warnings []string, err error) {
	if updated, ok := depricatedLabels[k]; ok {
		warnings = append(warnings, labeldeprecationNotice(k, updated))

		// TODO: validate value, with extra context?
	}

	// TODO: validate value

	return warnings, nil
}
