package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/policygroup"
)

// GetPolicyGroup requests the policy group with uid from the StorageOS API,
// translating it into a *policygroup.Resource.
func (o *OpenAPI) GetPolicyGroup(ctx context.Context, uid id.PolicyGroup) (*policygroup.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, resp, err := o.client.DefaultApi.GetPolicyGroup(ctx, uid.String())
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case notFoundError:
			return nil, apiclient.NewPolicyGroupIDNotFoundError(uid)
		default:
			return nil, v
		}
	}

	return o.codec.decodePolicyGroup(model)
}

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
