package apiclient

import (
	"context"
	"fmt"

	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

// GetNamespace requests basic information for the namespace resource which
// corresponds to uid from the StorageOS API.
func (c *Client) GetNamespace(ctx context.Context, uid id.Namespace) (*namespace.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.GetNamespace(ctx, uid)
}

// GetNamespaceByName requests basic information for the namespace resource
// which has the given name.
//
// The resource model for the API is build around using unique identifiers,
// so this operation is inherently more expensive than the corresponding
// GetNamespace() operation.
//
// Retrieving a namespace resource by name involves requesting a list of all
// namespaces from the StorageOS API and returning the first one where the
// name matches.
func (c *Client) GetNamespaceByName(ctx context.Context, name string) (*namespace.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	namespaces, err := c.transport.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	for _, ns := range namespaces {
		if ns.Name == name {
			return ns, nil
		}
	}

	return nil, NewNotFoundError(fmt.Sprintf("namespace with name %v not found", name))
}

// GetListNamespaces requests a list containing basic information for every namespace
// in the StorageOS cluster.
func (c *Client) GetAllNamespaces(ctx context.Context) ([]*namespace.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.ListNamespaces(ctx)
}
