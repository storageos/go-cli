package openapi

import (
	"context"

	"github.com/antihax/optional"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/openapi"
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

// DeleteNode makes a delete request for nodeID.
//
// The behaviour of the operation is dictated by params:
//
//   Version constraints:
//   - If params is nil or params.CASVersion is empty then the delete request is
//     unconditional
//   - If params.CASVersion is set, the request is conditional upon it matching
//     the node entity's version as seen by the server.
//
//   Asynchrony:
//   - If params is nil or params.AsyncMax is empty/zero valued then the delete
//     request is performed synchronously.
//   - If params.AsyncMax is set, the request is performed asynchronously using
//     the duration given as the maximum amount of time allowed for the request
//     before it times out.
func (o *OpenAPI) DeleteNode(ctx context.Context, nodeID id.Node, params *apiclient.DeleteNodeRequestParams) error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	var casVersion string
	var ignoreVersion optional.Bool = optional.NewBool(true)
	var asyncMax optional.String = optional.EmptyString()

	if params != nil {
		if params.CASVersion.String() != "" {
			ignoreVersion = optional.NewBool(false)
			casVersion = params.CASVersion.String()
		}

		if params.AsyncMax != 0 {
			asyncMax = optional.NewString(params.AsyncMax.String())
		}
	}

	resp, err := o.client.DefaultApi.DeleteNode(
		ctx,
		nodeID.String(),
		casVersion,
		&openapi.DeleteNodeOpts{
			IgnoreVersion: ignoreVersion,
			AsyncMax:      asyncMax,
		},
	)
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case notFoundError:
			return apiclient.NewNodeNotFoundError(nodeID)
		default:
			return v
		}
	}

	return nil
}
