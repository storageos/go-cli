package jointools

import (
	"fmt"
	"net"
	"strings"

	api "github.com/storageos/go-api"
	"github.com/storageos/go-cli/discovery"
)

// ExpandJOIN will expand each JOIN fragment out to the form <scheme>://<ip>:<port>.
func ExpandJOIN(join string) string {
	fragments := make([]string, 0)

	for _, frag := range strings.Split(strings.TrimSpace(join), ",") {
		fragments = append(fragments, ExpandJOINFragment(strings.TrimSpace(frag))...)
	}

	return strings.Join(fragments, ",")
}

// ExpandJOINFragment will, given an individal JOIN string fragment, expand it out
// to the form <scheme>://<ip>:<port>.
func ExpandJOINFragment(joinfragment string) []string {
	// Check to see if this is a discovery token
	if api.IsUUID(joinfragment) {
		return ExpandClusterToken(joinfragment)
	}

	var scheme string
	port := "5705"
	switch split := strings.Split(joinfragment, "://"); len(split) {
	case 2:
		scheme = split[0]
		split = split[1:len(split)]
		fallthrough
	case 1:
		// If there was no explicit scheme, default to http.
		if scheme == "" {
			scheme = "http"
		}
		switch hostPort := strings.Split(split[0], ":"); len(hostPort) {
		case 2:
			port = hostPort[1]
			fallthrough
		case 1:
			addrs := ExpandHost(hostPort[0])
			for i, addr := range addrs {
				addrs[i] = fmt.Sprintf("%s://%s:%s", scheme, addr, port)
			}
			return addrs
		}
	}

	return nil
}

// ExpandHost will perform a lookup on the host given to it,
// returning it's IPv4 address if successful.
func ExpandHost(host string) []string {
	if IsIPAddr(host) {
		return []string{host}
	}

	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil
	}

	// Only take the IPv4 addrs
	filtered := addrs[:0]
	for _, addr := range addrs {
		if ip := net.ParseIP(addr); ip != nil && ip.To4() != nil {
			filtered = append(filtered, addr)
		}
	}

	return filtered
}

// ExpandClusterToken will query the Discovery service for the cluster
// corresponding to the token given, performing the expansion on the
// JOIN fragments retrieved.
func ExpandClusterToken(token string) (nodes []string) {
	c, err := discovery.NewClient("", "", "")
	if err != nil {
		return nil
	}

	cluster, err := c.ClusterStatus(token)
	if err != nil {
		return nil
	}

	// Check all the host information known about this cluster ID
	for _, node := range cluster.Nodes {
		// Nodes could not have joined yet, this is not per-se an error
		if node.AdvertiseAddress != "" {
			nodes = append(nodes, ExpandJOINFragment(node.AdvertiseAddress)...)
		}
	}

	return nodes
}
