package openapi

import (
	"context"

	"github.com/antihax/optional"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/openapi"
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

// CreateNamespace requests the creation of a new namespace through the
// StorageOS API using the provided parameters.
func (o *OpenAPI) CreateNamespace(ctx context.Context, name string, labels labels.Set) (*namespace.Resource, error) {
	createData := openapi.CreateNamespaceData{
		Name:   name,
		Labels: labels,
	}

	model, resp, err := o.client.DefaultApi.CreateNamespace(ctx, createData)
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case badRequestError:
			return nil, apiclient.NewInvalidNamespaceCreationError(v.msg)
		case conflictError:
			return nil, apiclient.NewNamespaceExistsError(name)
		default:
			return nil, v
		}
	}

	return o.codec.decodeNamespace(model)
}

// DeleteNamespace makes a delete request for a namespace given its ID.
//
// The behaviour of the operation is dictated by params:
//
//
// 	Version constraints:
// 	- If params is nil or params.CASVersion is empty then the delete request is
// 	unconditional
// 	- If params.CASVersion is set, the request is conditional upon it matching
// 	the volume entity's version as seen by the server.
func (o *OpenAPI) DeleteNamespace(ctx context.Context, uid id.Namespace, params *apiclient.DeleteNamespaceRequestParams) error {
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

	resp, err := o.client.DefaultApi.DeleteNamespace(
		ctx,
		uid.String(),
		casVersion,
		&openapi.DeleteNamespaceOpts{
			IgnoreVersion: ignoreVersion,
		},
	)
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case notFoundError:
			return apiclient.NewNamespaceNotFoundError(uid)
		default:
			return v
		}
	}

	return nil
}
