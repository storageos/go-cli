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

// ParseHostPort returns a host:port string if the endpoint input is valid.
func ParseHostPort(endpoint string, defaultPort string) (string, error) {

	// plain IP address
	ip := net.ParseIP(endpoint)
	if ip != nil {
		return hostport(ip.String(), defaultPort)
	}

	// http, https or tcp endpoint
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") || strings.HasPrefix(endpoint, "tcp://") {
		u, err := url.Parse(endpoint)
		if err != nil {
			return "", err
		}
		parts := strings.Split(u.Host, ":")
		if len(parts) == 2 {
			return hostport(parts[0], u.Port())
		}
		return hostport(parts[0], defaultPort)
	}

	// hostname:port
	host, port, err := net.SplitHostPort(endpoint)
	if err == nil {
		return hostport(host, port)
	}

	// hostname or invalid input
	return hostport(endpoint, defaultPort)
}

// hostport validates host and port input and returns host:port
func hostport(host, port string) (string, error) {
	if host == "" || port == "" {
		return "", fmt.Errorf("invalid endpoint")
	}
	if !hostnameRegexp.MatchString(host) {
		return "", fmt.Errorf("invalid hostname format")
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return "", err
	}
	if p < 1 || p > 65535 {
		return "", fmt.Errorf("invalid port")
	}
	return fmt.Sprintf("%s:%s", host, port), nil
}
