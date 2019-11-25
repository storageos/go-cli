package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/namespace"
)

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
