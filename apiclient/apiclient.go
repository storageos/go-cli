// Package apiclient provides a type which implements an abstraction layer
// for consuming the storageos API programmatically.
package apiclient

import (
	"context"
	"fmt"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/c2-cli/volume"
)

// ConfigProvider describes the access to configuration values required by
// the apiclient package.
type ConfigProvider interface {
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
	Authenticate(ctx context.Context, username, password string) (*user.Resource, error)

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
	config    ConfigProvider
	transport Transport
}

// TODO: I think maybe this authenticate boiler plate should be moved down into
// the OpenAPI layer? That way we can be smart and avoid re-authing etcd without
// breaking abstraction layers.
func (c *Client) authenticate(ctx context.Context) (*user.Resource, error) {
	username, err := c.config.Username()
	if err != nil {
		return nil, err
	}
	password, err := c.config.Password()
	if err != nil {
		return nil, err
	}

	return c.transport.Authenticate(ctx, username, password)
}

// GetCluster requests basic information for the cluster resource from the
// StorageOS API.
func (c *Client) GetCluster(ctx context.Context) (*cluster.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.GetCluster(ctx)
}

// GetNode requests basic information for the node resource which
// corresponds to uid from the StorageOS API.
func (c *Client) GetNode(ctx context.Context, uid id.Node) (*node.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.GetNode(ctx, uid)
}

// GetListNodes requests a list containing basic information on each
// node resource in the cluster.
//
// The returned list is filtered using uids so that it contains only those
// resources which have a matching ID. Omitting uids will skip the filtering.
func (c *Client) GetListNodes(ctx context.Context, uids ...id.Node) ([]*node.Resource, error) {
	_, err := c.authenticate(ctx)
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
			return nil, fmt.Errorf("node %w: %v", ErrNotFound, id)
		}
	}

	return nodes, nil
}

// GetVolume requests basic information for the volume resource which
// corresponds to uid in namespace from the StorageOS API.
func (c *Client) GetVolume(ctx context.Context, namespace id.Namespace, uid id.Volume) (*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.GetVolume(ctx, namespace, uid)
}

// GetNamespaceVolumes requests basic information for each volume resource in
// namespace from the StorageOS API.
//
// The returned list is filtered using uids so that it contains only those
// resources which have a matching ID. Omitting uids will skip the filtering.
func (c *Client) GetNamespaceVolumes(ctx context.Context, namespace id.Namespace, uids ...id.Volume) ([]*volume.Resource, error) {
	_, err := c.authenticate(ctx)
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
			return nil, fmt.Errorf("volume %w: %v", ErrNotFound, id)
		}
	}

	return filtered, nil
}

// GetAllVolumes requests basic information for each volume resource in every
// namespace exposed by the StorageOS API to the authenticated user.
func (c *Client) GetAllVolumes(ctx context.Context) ([]*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.fetchAllVolumes(ctx)
}

// fetchAllVolumes requests the list of all namespaces from the StorageOS API,
// then requests the list of volumes within each namespace, returning an
// aggregate list of the volumes returned.
//
// If access is not granted when listing volumes for a retrieved namespace it
// is noted but will not return an error. Only if access is denied for all
// attempts will this return a permissions error.
func (c *Client) fetchAllVolumes(ctx context.Context) ([]*volume.Resource, error) {
	namespaces, err := c.transport.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	var volumes []*volume.Resource

	for _, ns := range namespaces {
		nsvols, err := c.transport.ListVolumes(ctx, ns.ID)
		switch err {
		case nil, ErrUnauthorised:
			// For these two errors, ignore - they're not fatal to the operation.
		default:
			return nil, err
		}
		volumes = append(volumes, nsvols...)
	}

	return volumes, nil
}

// DescribeNode requests detailed information for the node resource which
// corresponds to uid from the StorageOS API.
func (c *Client) DescribeNode(ctx context.Context, uid id.Node) (*node.State, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	resource, err := c.transport.GetNode(ctx, uid)
	if err != nil {
		return nil, err
	}

	volumes, err := c.fetchAllVolumes(ctx)
	if err != nil {
		return nil, err
	}

	deployments := deploymentsForNode(resource.ID, volumes)

	return &node.State{
		Resource:    resource,
		Deployments: deployments,
	}, nil
}

// DescribeListNodes requests a list containing detailed information on each
// node resource in the cluster.
//
// The returned list is filtered using uids so that it contains only those
// resources which have a matching ID. If none are specified, all are returned.
func (c *Client) DescribeListNodes(ctx context.Context, uids ...id.Node) ([]*node.State, error) {
	_, err := c.authenticate(ctx)
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
				return nil, fmt.Errorf("node %w: %v", ErrNotFound, id)
			}
		}

		resources = filtered
	}

	nodes := make([]*node.State, len(resources))

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

// New initialises a new Client using config for configuration settings,
// with transport providing the underlying implementation for encoding
// requests and decoding responses.
func New(transport Transport, config ConfigProvider) *Client {
	return &Client{
		transport: transport,
		config:    config,
	}
}
