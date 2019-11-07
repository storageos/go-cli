// Package apiclient provides a type which implements an abstraction layer
// for consuming the storageos API programmatically.
package apiclient

import (
	"context"
	"fmt"

	"code.storageos.net/storageos/c2-cli/pkg/cluster"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/namespace"
	"code.storageos.net/storageos/c2-cli/pkg/node"
	"code.storageos.net/storageos/c2-cli/pkg/volume"
)

// ConfigProvider describes the access to configuration values required by
// the apiclient package.
type ConfigProvider interface {
	APIEndpoint() (string, error)
	Username() (string, error)
	Password() (string, error)
}

// Transport describes the set of methods required by an API client to use a
// given transport implementation provider.
//
// A Transport implementation only needs to provide a direct mapping to the
// StorageOS API - it is the responsibility of the client to compose
// functionality for complex operations.
type Transport interface {
	Authenticate(ctx context.Context, username, password string) error

	GetCluster(context.Context) (*cluster.Resource, error)
	GetNode(context.Context, id.Node) (*node.Resource, error)
	GetVolume(context.Context, id.Namespace, id.Volume) (*volume.Resource, error)

	ListNodes(context.Context) ([]*node.Resource, error)
	ListVolumes(context.Context, id.Namespace) ([]*volume.Resource, error)
	ListNamespaces(context.Context) ([]*namespace.Resource, error)
}

// Client provides a collection of methods for consumers to interact with the
// StorageOS API.
type Client struct {
	transport Transport
	// TODO: Config options?
	username string
	password string
}

// GetCluster requests basic information for the cluster resource from the
// StorageOS API.
func (c *Client) GetCluster(ctx context.Context) (*cluster.Resource, error) {
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	cluster, err := c.transport.GetCluster(ctx)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

// GetNode requests basic information for the node resource which
// corresponds to uid from the StorageOS API.
func (c *Client) GetNode(ctx context.Context, uid id.Node) (*node.Resource, error) {
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	n, err := c.transport.GetNode(ctx, uid)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// GetListNodes requests a list containing basic information on each
// node resource in the cluster.
//
// The returned list is filtered using uids so that it contains only those
// resources which have a matching ID. Omitting uids will skip the filtering.
func (c *Client) GetListNodes(ctx context.Context, uids ...id.Node) ([]*node.Resource, error) {
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	nodes, err := c.transport.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	if len(uids) == 0 {
		return nodes, nil
	}

	// Filter uids have been provided:
	retrieved := map[id.Node]*node.Resource{}

	for _, n := range nodes {
		retrieved[n.ID] = n
	}

	filtered := make([]*node.Resource, len(uids))

	i := 0
	for _, id := range uids {
		n, ok := retrieved[id]
		if ok {
			filtered[i] = n
			i++
		} else {
			return nil, fmt.Errorf("node not found: %v", id)
		}
	}

	return nodes, nil
}

// GetVolume requests basic information for the volume resource which
// corresponds to uid in namespace from the StorageOS API.
func (c *Client) GetVolume(ctx context.Context, namespace id.Namespace, uid id.Volume) (*volume.Resource, error) {
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	v, err := c.transport.GetVolume(ctx, namespace, uid)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// GetNamespaceVolumes requests basic information for each volume resource in
// namespace from the StorageOS API.
//
// The returned list is filtered using uids so that it contains only those
// resources which have a matching ID. Omitting uids will skip the filtering.
func (c *Client) GetNamespaceVolumes(ctx context.Context, namespace id.Namespace, uids ...id.Volume) ([]*volume.Resource, error) {
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	volumes, err := c.transport.ListVolumes(ctx, namespace)
	if err != nil {
		return nil, err
	}

	if len(uids) == 0 {
		return volumes, nil
	}

	// Filter uids have been provided:
	retrieved := map[id.Volume]*volume.Resource{}

	for _, v := range volumes {
		retrieved[v.ID] = v
	}

	filtered := make([]*volume.Resource, len(uids))

	i := 0
	for _, id := range uids {
		v, ok := retrieved[id]
		if ok {
			filtered[i] = v
			i++
		} else {
			return nil, fmt.Errorf("volume not found: %v", id)
		}
	}

	return filtered, nil
}

// GetAllVolumes requests basic information for each volume resource in every
// namespace exposed by the StorageOS API.
//
// TODO:
// If access is not granted when listing volumes for a retrieved namespace it
// is noted but will not return an error. Only if access is denied for all
// attempts will this return a permissions error.
func (c *Client) GetAllVolumes(ctx context.Context) ([]*volume.Resource, error) {
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	return c.fetchAllVolumes(ctx)
}

func (c *Client) fetchAllVolumes(ctx context.Context) ([]*volume.Resource, error) {
	namespaces, err := c.transport.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	var volumes []*volume.Resource

	for _, ns := range namespaces {
		nsvols, err := c.transport.ListVolumes(ctx, ns.ID)
		if err != nil {
			return nil, err
		}
		volumes = append(volumes, nsvols...)
	}

	return volumes, nil
}

// DescribeNode requests detailed information for the node resource which
// corresponds to uid from the StorageOS API.
func (c *Client) DescribeNode(ctx context.Context, uid id.Node) (*node.State, error) {
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	resource, err := c.transport.GetNode(ctx, uid)
	if err != nil {
		return nil, err
	}

	// TODO: For the retrieved node we then need to build the detailed
	// information by performing other API requests.
	volumes, err := c.fetchAllVolumes(ctx)
	if err != nil {
		return nil, err
	}

	deployments := deploymentsForNode(resource.ID, volumes)

	n := &node.State{
		Resource:    resource,
		Deployments: deployments,
	}

	return n, nil
}

// DescribeListNodes requests a list containing detailed information on each
// node resource in the cluster.
//
// The returned list is filtered using uids so that it contains only those
// resources which have a matching ID. If none are specified, all are returned.
func (c *Client) DescribeListNodes(ctx context.Context, uids ...id.Node) ([]*node.State, error) {
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	resources, err := c.transport.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	if len(uids) > 0 {
		retrieved := map[id.Node]*node.Resource{}

		for _, n := range resources {
			retrieved[n.ID] = n
		}

		filtered := make([]*node.Resource, len(uids))

		i := 0
		for _, id := range uids {
			n, ok := retrieved[id]
			if ok {
				filtered[i] = n
				i++
			} else {
				return nil, fmt.Errorf("node not found: %v", id)
			}
		}

		resources = filtered
	}

	nodes := make([]*node.State, len(resources))

	// TODO: For each node resource we then need to build the detailed
	// information by performing other API requests.
	volumes, err := c.fetchAllVolumes(ctx)
	if err != nil {
		return nil, err
	}

	deploymentMap := mapNodeDeployments(volumes)

	for i, r := range resources {
		nodes[i] = &node.State{
			Resource:    r,
			Deployments: deploymentMap[r.ID], // No need to check - zero value is ok.
		}
	}

	return nodes, nil
}

// New initialises a new Client configured from config, with transport
// providing the underlying implementation for encoding requests and
// decoding responses.
func New(transport Transport, config ConfigProvider) (*Client, error) {
	username, err := config.Username()
	if err != nil {
		return nil, err
	}

	password, err := config.Password()
	if err != nil {
		return nil, err
	}

	return &Client{
		transport: transport,
		username:  username,
		password:  password,
	}, nil
}
