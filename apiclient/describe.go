package apiclient

import (
	"context"
	"fmt"

	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

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
				return nil, NewNotFoundError(fmt.Sprintf("node %v not found", id))
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
