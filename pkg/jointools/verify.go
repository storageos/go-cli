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

// VerifyJOIN will separate the members of a JOIN string into
// fragments and verify that each is valid.
func VerifyJOIN(discoveryHost, join string) (errs []error) {
	// Split into fragments by ',' and verify the fragments
	for _, frag := range strings.Split(strings.TrimSpace(join), ",") {
		errs = append(errs, VerifyJOINFragment(discoveryHost, strings.TrimSpace(frag))...)
	}

	return errs
}

// VerifyJOINFragment will, given a fragment from a JOIN string,
// verify that it is valid for use.
//
// In the instance that the value passed into it is a UUID,
// the function will have the side effect of verifying that
// value is a cluster token.
func VerifyJOINFragment(discoveryHost, joinfragment string) (errs []error) {
	// Check to see if this is a discovery token
	if api.IsUUID(joinfragment) {
		return VerifyClusterToken(discoveryHost, joinfragment)
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

// VerifyHost will, given the host part of a URI authority segment, verify that
// it is either a valid IP address or attempt to resolve it to a hostname.
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

// VerifyClusterToken will, given a token string, query the Discovery service
// to retrieve information about the cluster if it exists. Further verification
// is then performed on the node addresses retrieved.
func VerifyClusterToken(discoveryHost, token string) (errs []error) {
	c, err := discovery.NewClient(discoveryHost, "", "")
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
			errs = append(errs, VerifyJOINFragment(discoveryHost, node.AdvertiseAddress)...)
		}
	}

	return errs
}

// IsValidScheme verifies whether or not the string value given to
// it is a valid URI scheme for a JOIN member.
func IsValidScheme(s string) bool {
	switch s {
	case "http", "tcp", "https":
		return true
	default:
		return false
	}
}

// IsValidPort verifies whether or not the string value given to
// it is a valid port number.
func IsValidPort(p string) bool {
	port, err := strconv.Atoi(p)

	return (err == nil) && (port > 0) && (port <= 65535)
}

// IsHostPort checks that the input takes the form <hostname>:<port>.
func IsHostPort(s string) bool {
	if strings.Count(s, ":") != 1 {
		return false
	}

	port := strings.Split(s, ":")[1]
	return IsValidPort(port)
}

// IsIPAddr attempts to parse the string into an IP address
// returning whether or not it succeeded.
func IsIPAddr(s string) bool {
	return net.ParseIP(s) != nil
}
