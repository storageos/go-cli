package apiclient

import (
	"context"
	"fmt"

	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

// NamespaceNotFoundError indicates that the API could not find the StorageOS
// namespace specified.
type NamespaceNotFoundError struct {
	uid  id.Namespace
	name string
}

// Error returns an error message indicating that the namespace with a given
// ID or name was not found, as configured.
func (e NamespaceNotFoundError) Error() string {
	switch {
	case e.uid != "":
		return fmt.Sprintf("namespace with ID %v not found", e.uid)
	case e.name != "":
		return fmt.Sprintf("namespace with name %v not found", e.name)
	}

	return fmt.Sprintf("namespace not found")
}

// NewNamespaceNotFoundError returns a NamespaceNotFoundError for the namespace
// with uid.
func NewNamespaceNotFoundError(uid id.Namespace) NamespaceNotFoundError {
	return NamespaceNotFoundError{
		uid: uid,
	}
}

// NewNamespaceNameNotFoundError returns a NamespaceNotFoundError for the
// namespace with name.
func NewNamespaceNameNotFoundError(name string) NamespaceNotFoundError {
	return NamespaceNotFoundError{
		name: name,
	}
}

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

	return nil, NewNamespaceNameNotFoundError(name)
}

// GetListNamespaces requests a list of namespace resources present in the
// cluster.
//
// The returned list is filtered using uids so that it contains only those
// namespace resources which have a matching ID. If no uids are given then
// all namespaces are returned.
func (c *Client) GetListNamespaces(ctx context.Context, uids ...id.Namespace) ([]*namespace.Resource, error) {
	resources, err := c.GetAllNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	return filterNamespacesForUIDs(resources, uids...)
}

// GetListNamespacesByName requests a list of namespace resources present in
// the cluster.
//
// The returned list is filtered using names so that it contains only those
// namespaces resources which have a matching name. If no names are given then
// all namespaces are returned.
func (c *Client) GetListNamespacesByName(ctx context.Context, names ...string) ([]*namespace.Resource, error) {
	resources, err := c.GetAllNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	return filterNamespacesForNames(resources, names...)
}

// GetAllNamespaces requests a list containing basic information for every namespace
// in the StorageOS cluster.
func (c *Client) GetAllNamespaces(ctx context.Context) ([]*namespace.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.ListNamespaces(ctx)
}

// filterNamespacesForNames will return a subset of namespaces containing
// resources which have one of the provided names. If names is not provided,
// namespaces is returned as is.
//
// If there is no resource for a given name then an error is returned, thus
// this is a strict helper.
func filterNamespacesForNames(namespaces []*namespace.Resource, names ...string) ([]*namespace.Resource, error) {
	if len(names) == 0 {
		return namespaces, nil
	}

	retrieved := map[string]*namespace.Resource{}

	for _, ns := range namespaces {
		retrieved[ns.Name] = ns
	}

	filtered := make([]*namespace.Resource, 0, len(names))
	for _, name := range names {
		ns, ok := retrieved[name]
		if !ok {
			return nil, NewNamespaceNameNotFoundError(name)
		}
		filtered = append(filtered, ns)
	}

	return filtered, nil
}

// filterNamespacesForUIDS will return a subset of namespaces containing
// resources which have one of the provided uids. If uids is not provided,
// namespaces is returned as is.
//
// If there is no resource for a given uid then an error is returned, thus
// this is a strict helper.
func filterNamespacesForUIDs(namespaces []*namespace.Resource, uids ...id.Namespace) ([]*namespace.Resource, error) {
	if len(uids) == 0 {
		return namespaces, nil
	}

	retrieved := map[id.Namespace]*namespace.Resource{}

	for _, ns := range namespaces {
		retrieved[ns.ID] = ns
	}

	filtered := make([]*namespace.Resource, 0, len(uids))

	for _, id := range uids {
		ns, ok := retrieved[id]
		if !ok {
			return nil, NewNamespaceNotFoundError(id)
		}
		filtered = append(filtered, ns)
	}

	return filtered, nil
}
