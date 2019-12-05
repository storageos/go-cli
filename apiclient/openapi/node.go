package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

func (o *OpenAPI) GetNode(ctx context.Context, uid id.Node) (*node.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, resp, err := o.client.DefaultApi.GetNode(ctx, uid.String())
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	return o.codec.decodeNode(model)
}

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
