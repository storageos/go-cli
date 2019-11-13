package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
)

func (o *OpenAPI) GetCluster(ctx context.Context) (*cluster.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, _, err := o.client.DefaultApi.GetCluster(ctx)
	if err != nil {
		return nil, err
	}

	return o.codec.decodeCluster(model)
}

func (o *OpenAPI) GetNode(ctx context.Context, uid id.Node) (*node.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, resp, err := o.client.DefaultApi.GetNode(ctx, uid.String())
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	return o.codec.decodeNode(model)
}

func (o *OpenAPI) GetVolume(ctx context.Context, namespace id.Namespace, uid id.Volume) (*volume.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, resp, err := o.client.DefaultApi.GetVolume(ctx, namespace.String(), uid.String())
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	return o.codec.decodeVolume(model)
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

func (o *OpenAPI) ListVolumes(ctx context.Context, namespace id.Namespace) ([]*volume.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	models, resp, err := o.client.DefaultApi.ListVolumes(ctx, namespace.String())
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	volumes := make([]*volume.Resource, len(models))
	for i, m := range models {
		v, err := o.codec.decodeVolume(m)
		if err != nil {
			return nil, err
		}

		volumes[i] = v
	}

	return volumes, nil
}

func (o *OpenAPI) ListNamespaces(ctx context.Context) ([]*namespace.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	models, resp, err := o.client.DefaultApi.ListNamespaces(ctx)
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	namespaces := make([]*namespace.Resource, len(models))
	for i, m := range models {
		ns, err := o.codec.decodeNamespace(m)
		if err != nil {
			return nil, err
		}

		namespaces[i] = ns
	}

	return namespaces, nil
}
