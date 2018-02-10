package jointools

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	api "github.com/storageos/go-api"
	"github.com/storageos/go-cli/discovery"
)

func VerifyJOIN(join string) (errs []error) {
	// Split into fragments by ',' and verify the fragments
	for _, frag := range strings.Split(strings.TrimSpace(join), ",") {
		errs = append(errs, VerifyJOINFragment(strings.TrimSpace(frag))...)
	}

	return errs
}

func VerifyJOINFragment(joinfragment string) (errs []error) {
	// Check to see if this is a discovery token
	if api.IsUUID(joinfragment) {
		return VerifyClusterToken(joinfragment)
	}

	split := strings.Split(joinfragment, "://")
	switch strings.Count(joinfragment, "://") {
	case 1:
		if !IsValidScheme(split[0]) {
			errs = append(errs, fmt.Errorf("invalid scheme '%s'", split[0]))
		}

		split = split[1:len(split)] // cut off the scheme and continue
		fallthrough

	case 0:
		switch hostPort := strings.Split(split[0], ":"); len(hostPort) {
		case 2:
			if !IsValidPort(hostPort[1]) {
				errs = append(errs, fmt.Errorf("invalid port '%s'", hostPort[1]))
			}
			fallthrough

		case 1:
			errs = append(errs, VerifyHost(hostPort[0])...)

		default:
			errs = append(errs, errors.New("empty or invalid join fragment"))
		}

	// Multiple '://'?
	default:
		errs = append(errs, fmt.Errorf("invalid join fragment '%s'", joinfragment))
	}

	return errs
}

func VerifyHost(host string) []error {
	if IsIPAddr(host) {
		return nil // valid IP address
	}

	// Probably a host or fqdn
	_, err := net.LookupHost(host)
	if err != nil {
		return []error{fmt.Errorf("failed to lookup host '%s', %s", host, err)}
	}

	return nil
}

func VerifyClusterToken(token string) (errs []error) {
	c, err := discovery.NewClient("", "", "")
	if err != nil {
		return []error{fmt.Errorf("failed to query discovery service, %s", err)}
	}

	cluster, err := c.ClusterStatus(token)
	if err != nil {
		return []error{fmt.Errorf("failed to query discovery service, %s", err)}
	}

	// Check all the host information known about this cluster ID
	for _, node := range cluster.Nodes {
		// Nodes could not have joined yet, this is not per-se an error
		if node.AdvertiseAddress != "" {
			errs = append(errs, VerifyJOINFragment(node.AdvertiseAddress)...)
		}
	}

	return errs
}

func IsValidScheme(s string) bool {
	switch s {
	case "http", "tcp":
		return true

	default:
		return false
	}
}

func IsValidPort(p string) bool {
	port, err := strconv.Atoi(p)

	return (err == nil) && (port > 0) && (port <= 65535)
}

func IsSchemeHost(s string) bool {
	if strings.Count(s, ":") != 1 {
		return false
	}

	if _, err := strconv.Atoi(strings.Split(s, ":")[1]); err == nil {
		return false
	}

	return true
}

func IsHostPort(s string) bool {
	if strings.Count(s, ":") != 1 {
		return false
	}

	if _, err := strconv.Atoi(strings.Split(s, ":")[1]); err != nil {
		return false
	}

	return true
}

func IsIPAddr(s string) bool {
	return net.ParseIP(s) != nil
}
