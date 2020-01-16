package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

// GetNamespace requests the namespace with uid from the StorageOS API,
// translating it into a *namespace.Resource.
func (o *OpenAPI) GetNamespace(ctx context.Context, uid id.Namespace) (*namespace.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, resp, err := o.client.DefaultApi.GetNamespace(ctx, uid.String())
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case notFoundError:
			return nil, apiclient.NewNamespaceNotFoundError(uid)
		default:
			return nil, v
		}
	}

	return o.codec.decodeNamespace(model)
}

// ListNamespaces requests a list of all namespaces from the StorageOS API,
// translating each one to a *namespace.Resource.
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
