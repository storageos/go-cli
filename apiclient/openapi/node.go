package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

// GetNode requests the node with uid from the StorageOS API, translating it
// into a *node.Resource.
func (o *OpenAPI) GetNode(ctx context.Context, uid id.Node) (*node.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, resp, err := o.client.DefaultApi.GetNode(ctx, uid.String())
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case notFoundError:
			return nil, apiclient.NewNodeNotFoundError(uid)
		default:
			return nil, v
		}
	}

	return o.codec.decodeNode(model)
}

// ListNodes requests a list of all nodes from the StorageOS API, translating
// each one to a *node.Resource.
func (o *OpenAPI) ListNodes(ctx context.Context) ([]*node.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	models, resp, err := o.client.DefaultApi.ListNodes(ctx)
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	nodes := make([]*node.Resource, len(models))
	for i, m := range models {
		n, err := o.codec.decodeNode(m)
		if err != nil {
			return nil, err
		}

		nodes[i] = n
	}

	return nodes, nil
}
