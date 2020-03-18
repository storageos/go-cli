package apiclient

import (
	"context"
	"fmt"

	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
)

// NodeNotFoundError indicates that the API could not find the StorageOS node
// specified.
type NodeNotFoundError struct {
	uid  id.Node
	name string
}

// Error returns an error message indicating that the node with a given
// ID or name was not found, as configured.
func (e NodeNotFoundError) Error() string {
	switch {
	case e.uid != "":
		return fmt.Sprintf("node with ID %v not found", e.uid)
	case e.name != "":
		return fmt.Sprintf("node with name %v not found", e.name)
	}

	return "node not found"
}

// NewNodeNotFoundError returns a NodeNotFoundError for the node with uid.
func NewNodeNotFoundError(uid id.Node) NodeNotFoundError {
	return NodeNotFoundError{
		uid: uid,
	}
}

// NewNodeNameNotFoundError returns a NodeNotFoundError for the node with name.
func NewNodeNameNotFoundError(name string) NodeNotFoundError {
	return NodeNotFoundError{
		name: name,
	}
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

// GetNodeByName requests basic information for the node resource which has
// name.
//
// The resource model for the API is build around using unique identifiers,
// so this operation is inherently more expensive than the corresponding
// GetNode() operation.
//
// Retrieving a node resource by name involves requesting a list of all nodes
// in the cluster from the StorageOS API and returning the first node where the
// name matches.
func (c *Client) GetNodeByName(ctx context.Context, name string) (*node.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	nodes, err := c.transport.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		if n.Name == name {
			return n, nil
		}
	}

	return nil, NewNodeNameNotFoundError(name)
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

	return filterNodesForUIDs(nodes, uids...)
}

// GetListNodesByName requests a list containing basic information on each
// node resource in the cluster.
//
// The returned list is filtered using names so that it contains only those
// resources which have a matching name. Omitting names will skip the filtering.
func (c *Client) GetListNodesByName(ctx context.Context, names ...string) ([]*node.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	nodes, err := c.transport.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	return filterNodesForNames(nodes, names...)
}

// DescribeNode requests detailed information for the node resource which
// corresponds to uid from the StorageOS API.
func (c *Client) DescribeNode(ctx context.Context, uid id.Node) (*node.State, error) {
	resource, err := c.GetNode(ctx, uid)
	if err != nil {
		return nil, err
	}

	return c.describeNode(ctx, resource)
}

// DescribeNodeByName requests detailed information for the node resource which
// has name from the StorageOS API.
//
// The resource model for the API is build around using unique identifiers,
// so this operation is inherently more expensive than the corresponding
// DescribeNode() operation.
//
// Retrieving a node state by name involves requesting a list of all nodes
// in the cluster from the StorageOS API and gathering information about the
// the first node where the name matches using its unique identifier.
func (c *Client) DescribeNodeByName(ctx context.Context, name string) (*node.State, error) {
	resource, err := c.GetNodeByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return c.describeNode(ctx, resource)
}

// describeNode gathers extra information about the node resource, constructing
// a node state with it.
func (c *Client) describeNode(ctx context.Context, resource *node.Resource) (*node.State, error) {
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

	resources, err = filterNodesForUIDs(resources, uids...)
	if err != nil {
		return nil, err
	}

	return c.getNodeStates(ctx, resources)
}

// DescribeListNodesByName requests a list containing detailed information on each
// node resource in the cluster.
//
// The returned list is filtered using names so that it contains only those
// resources which have a matching name. If none are specified, all are returned.
func (c *Client) DescribeListNodesByName(ctx context.Context, names ...string) ([]*node.State, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	resources, err := c.transport.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	resources, err = filterNodesForNames(resources, names...)
	if err != nil {
		return nil, err
	}

	return c.getNodeStates(ctx, resources)
}

// getNodeStates is a helper function to build a slice of *node.State from
// resources. To accurately construct a *node.State from a *node.Resource
// requires knowledge of the complete volume set so this is expensive.
func (c *Client) getNodeStates(ctx context.Context, resources []*node.Resource) ([]*node.State, error) {
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

// filterNodesForNames will return a subset of nodes containing resources
// which have one of the provided names. If names is not provided, nodes is
// returned as is.
//
// If there is no resource for a given name then an error is returned, thus
// this is a strict helper.
func filterNodesForNames(nodes []*node.Resource, names ...string) ([]*node.Resource, error) {
	// return everything if no filter names given
	if len(names) == 0 {
		return nodes, nil
	}

	retrieved := map[string]*node.Resource{}

	for _, n := range nodes {
		retrieved[n.Name] = n
	}

	filtered := make([]*node.Resource, 0, len(names))

	for _, name := range names {
		n, ok := retrieved[name]
		if !ok {
			return nil, NewNodeNameNotFoundError(name)
		}
		filtered = append(filtered, n)
	}

	return filtered, nil
}

// filterNodesForUIDs will return a subset of nodes containing resources
// which have one of the provided uids. If uids is not provided, nodes is
// returned as is.
//
// If there is no resource for a given uid then an error is returned, thus
// this is a strict helper.
func filterNodesForUIDs(nodes []*node.Resource, uids ...id.Node) ([]*node.Resource, error) {
	// return everything if no filter uids given
	if len(uids) == 0 {
		return nodes, nil
	}

	retrieved := map[id.Node]*node.Resource{}

	for _, n := range nodes {
		retrieved[n.ID] = n
	}

	filtered := make([]*node.Resource, 0, len(uids))

	for _, id := range uids {
		n, ok := retrieved[id]
		if !ok {
			return nil, NewNodeNotFoundError(id)
		}
		filtered = append(filtered, n)
	}

	return filtered, nil
}

// deploymentsForNode is a utility function returning the list of all
// deployments located on node within volumes.
func deploymentsForNode(uid id.Node, volumes []*volume.Resource) []*node.Deployment {
	var deployments []*node.Deployment
	for _, v := range volumes {
		if v.Master.Node == uid {
			deployments = append(
				deployments,
				&node.Deployment{
					VolumeID:   v.ID,
					VolumeName: v.Name,
					Deployment: v.Master,
				},
			)
		}

		for _, r := range v.Replicas {
			if r.Node == uid {
				deployments = append(
					deployments,
					&node.Deployment{
						VolumeID:   v.ID,
						VolumeName: v.Name,
						Deployment: r,
					},
				)
			}
		}
	}
	return deployments
}

// mapNodeDeployments builds a mapping from node ID to hosted deployments
// for the list of volumes.
func mapNodeDeployments(volumes []*volume.Resource) map[id.Node][]*node.Deployment {
	deployMap := make(map[id.Node][]*node.Deployment)

	for _, v := range volumes {
		deployMap[v.Master.Node] = append(
			deployMap[v.Master.Node],
			&node.Deployment{
				VolumeID:   v.ID,
				VolumeName: v.Name,
				Deployment: v.Master,
			},
		)

		for _, r := range v.Replicas {
			deployMap[r.Node] = append(
				deployMap[r.Node],
				&node.Deployment{
					VolumeID:   v.ID,
					VolumeName: v.Name,
					Deployment: r,
				},
			)
		}
	}

	return deployMap
}
