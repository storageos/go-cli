package validation

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	storageos "github.com/storageos/go-api"
)

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

func ParseHostPort(host string, defaultPort string) (string, error) {
	host = strings.TrimSuffix(host, "/")

	validHostname := regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)

	switch strings.Count(host, ":") {
	// No port number found
	case 0:
		if defaultPort == "" {
			return "", fmt.Errorf("invalid value: '%v' dosn't have a port number", host)
		}

		host += ":" + defaultPort
		fallthrough

	case 1:
		s := strings.Split(host, ":")

		if strings.HasPrefix(s[1], "//") {
			h := strings.TrimPrefix(s[1], "//")

			if net.ParseIP(h) == nil && !validHostname.MatchString(h) {
				return "", fmt.Errorf("invalid value: '%v' is not a valid hostname or IP address", h)
			}

			if defaultPort == "" {
				return "", fmt.Errorf("invalid value: '%v' doesn't have a port number", host)
			}

			return h + ":" + defaultPort, nil
		}

		if i, err := strconv.Atoi(s[1]); err != nil || i > 0xFFFF {
			return "", fmt.Errorf("invalid value: '%v' is not a valid port number", s[1])
		}

		if net.ParseIP(s[0]) == nil && !validHostname.MatchString(s[0]) {
			return "", fmt.Errorf("invalid value: '%v' is not a valid hostname or IP address", s[0])
		}

		return host, nil

	case 2:
		u, err := url.Parse(host)
		if err != nil {
			return "", fmt.Errorf("invalid value: %v", err)
		}

		h, p := u.Hostname(), u.Port()

		if h == "" {
			return "", fmt.Errorf("invalid value: '%s' is not a valid hostname", h)
		}

		if p == "" {
			if defaultPort == "" {
				return "", fmt.Errorf("invalid value: '%s' is not a valid port", p)
			}

			p = defaultPort
		}

		return fmt.Sprintf("%s:%s", h, p), nil

	// Unrecognised format
	default:
		return "", fmt.Errorf("invalid value: '%s'", host)
	}
}
