package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/cluster"
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
