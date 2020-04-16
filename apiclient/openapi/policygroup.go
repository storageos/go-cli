package openapi

import (
	"context"

	"github.com/antihax/optional"

	"code.storageos.net/storageos/openapi"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/policygroup"
)

// CreatePolicyGroup requests the creation of a new policy group through the
// StorageOS API using the provided parameters.
func (o *OpenAPI) CreatePolicyGroup(ctx context.Context, name string, specs []*policygroup.Spec) (*policygroup.Resource, error) {
	slice := make([]openapi.PoliciesSpecs, 0, len(specs))
	for _, s := range specs {
		slice = append(slice, openapi.PoliciesSpecs{
			NamespaceID:  s.NamespaceID.String(),
			ResourceType: s.ResourceType,
			ReadOnly:     s.ReadOnly,
		})
	}

	createData := openapi.CreatePolicyGroupData{
		Name:  name,
		Specs: &slice,
	}

	model, resp, err := o.client.DefaultApi.CreatePolicyGroup(ctx, createData)
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case badRequestError:
			return nil, apiclient.NewInvalidPolicyGroupCreationError(v.msg)
		case conflictError:
			return nil, apiclient.NewPolicyGroupExistsError(name)
		default:
			return nil, v
		}
	}

	return o.codec.decodePolicyGroup(model)
}

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

// DeletePolicyGroup makes a delete request for a policy group given its ID.
//
// The behaviour of the operation is dictated by params:
//
// 	Version constraints:
//  - If params is nil or params.CASVersion is empty then the delete request is
//    unconditional
//  - If params.CASVersion is set, the request is conditional upon it matching
//    the volume entity's version as seen by the server.
func (o *OpenAPI) DeletePolicyGroup(ctx context.Context, uid id.PolicyGroup, params *apiclient.DeletePolicyGroupRequestParams) error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	var casVersion string
	var ignoreVersion optional.Bool = optional.NewBool(true)

	if params != nil {
		if params.CASVersion.String() != "" {
			ignoreVersion = optional.NewBool(false)
			casVersion = params.CASVersion.String()
		}
	}

	resp, err := o.client.DefaultApi.DeletePolicyGroup(
		ctx,
		uid.String(),
		casVersion,
		&openapi.DeletePolicyGroupOpts{
			IgnoreVersion: ignoreVersion,
		},
	)
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case notFoundError:
			return apiclient.NewPolicyGroupIDNotFoundError(uid)
		default:
			return v
		}
	}

	return nil
}
