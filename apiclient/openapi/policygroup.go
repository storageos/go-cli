package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/policygroup"
)

// ListPolicyGroups returns a list of all policy groups from the StorageOS API,
// translating each one to a *policygroup.Resource.
func (o *OpenAPI) ListPolicyGroups(ctx context.Context) ([]*policygroup.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	models, resp, err := o.client.DefaultApi.ListPolicyGroups(ctx)
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	policyGroups := make([]*policygroup.Resource, 0, len(models))
	for _, m := range models {
		g, err := o.codec.decodePolicyGroup(m)
		if err != nil {
			return nil, err
		}

		policyGroups = append(policyGroups, g)
	}

	return policyGroups, nil
}
