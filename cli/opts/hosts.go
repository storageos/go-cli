package opts

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

var (

	// DefaultHTTPHost Default HTTP Host used if only port is provided to -H flag e.g. dockerd -H tcp://:8080
	DefaultHTTPHost = "localhost"
	// DefaultHTTPPort Default HTTP Port used if only the protocol is provided to -H flag e.g. dockerd -H tcp://
	DefaultHTTPPort = 5705 // Default HTTP Port
	// DefaultUnixSocket Path for the unix socket.
	DefaultUnixSocket = "/var/run/storageos.sock"
	// DefaultTCPHost constant defines the default host string used by docker on Windows
	DefaultTCPHost = fmt.Sprintf("tcp://%s:%d", DefaultHTTPHost, DefaultHTTPPort)
	// DefaultHost constant defines the default host string used by docker on other hosts than Windows
	DefaultHost = DefaultTCPHost
)

// parseSimpleProtoAddr parses and validates that the specified address is a valid
// socket address for simple protocols like unix and npipe. It returns a formatted
// socket address, either using the address parsed from addr, or the contents of
// defaultAddr if addr is a blank string.
func parseSimpleProtoAddr(proto, addr, defaultAddr string) (string, error) {
	addr = strings.TrimPrefix(addr, proto+"://")
	if strings.Contains(addr, "://") {
		return "", fmt.Errorf("Invalid proto, expected %s: %s", proto, addr)
	}
	if addr == "" {
		addr = defaultAddr
	}
	return fmt.Sprintf("%s://%s", proto, addr), nil
}

// ParseTCPAddr parses and validates that the specified address is a valid TCP
// address. It returns a formatted TCP address, either using the address parsed
// from tryAddr, or the contents of defaultAddr if tryAddr is a blank string.
// tryAddr is expected to have already been Trim()'d
// defaultAddr must be in the full `tcp://host:port` form
func ParseTCPAddr(tryAddr string, defaultAddr string) (string, error) {
	if tryAddr == "" || tryAddr == "tcp://" {
		return defaultAddr, nil
	}
	addr := strings.TrimPrefix(tryAddr, "tcp://")
	if strings.Contains(addr, "://") || addr == "" {
		return "", fmt.Errorf("Invalid proto, expected tcp: %s", tryAddr)
	}

	defaultAddr = strings.TrimPrefix(defaultAddr, "tcp://")
	defaultHost, defaultPort, err := net.SplitHostPort(defaultAddr)
	if err != nil {
		return "", err
	}
	// url.Parse fails for trailing colon on IPv6 brackets on Go 1.5, but
	// not 1.4. See https://github.com/golang/go/issues/12200 and
	// https://github.com/golang/go/issues/6530.
	if strings.HasSuffix(addr, "]:") {
		addr += defaultPort
	}

	u, err := url.Parse("tcp://" + addr)
	if err != nil {
		return "", err
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		// try port addition once
		host, port, err = net.SplitHostPort(net.JoinHostPort(u.Host, defaultPort))
	}
	if err != nil {
		return "", fmt.Errorf("Invalid bind address format: %s", tryAddr)
	}

	if host == "" {
		host = defaultHost
	}
	if port == "" {
		port = defaultPort
	}
	p, err := strconv.Atoi(port)
	if err != nil && p == 0 {
		return "", fmt.Errorf("Invalid bind address format: %s", tryAddr)
	}

	return fmt.Sprintf("tcp://%s%s", net.JoinHostPort(host, port), u.Path), nil
}

// ValidateExtraHost validates that the specified string is a valid extrahost and returns it.
// ExtraHost is in the form of name:ip where the ip has to be a valid ip (IPv4 or IPv6).
func ValidateExtraHost(val string) (string, error) {
	// allow for IPv6 addresses in extra hosts by only splitting on first ":"
	arr := strings.SplitN(val, ":", 2)
	if len(arr) != 2 || len(arr[0]) == 0 {
		return "", fmt.Errorf("bad format for add-host: %q", val)
	}
	if _, err := ValidateIPAddress(arr[1]); err != nil {
		return "", fmt.Errorf("invalid IP address in add-host: %q", arr[1])
	}
	return val, nil
}
