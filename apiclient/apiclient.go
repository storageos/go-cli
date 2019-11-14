// Package apiclient provides a type which implements an abstraction layer
// for consuming the storageos API programmatically.
package apiclient

import (
	"context"
	"errors"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/c2-cli/volume"
)

// ConfigProvider describes the access to configuration values required by
// the apiclient package.
type ConfigProvider interface {
	Username() (string, error)
	Password() (string, error)
}

// Transport describes the set of methods required by an API client to use a
// given transport implementation provider.
//
// A Transport implementation only needs to provide a direct mapping to the
// StorageOS API - it is the responsibility of the client to compose
// functionality for complex operations.
type Transport interface {
	Authenticate(ctx context.Context, username, password string) (*user.Resource, error)

	GetCluster(context.Context) (*cluster.Resource, error)
	GetNode(context.Context, id.Node) (*node.Resource, error)
	GetVolume(context.Context, id.Namespace, id.Volume) (*volume.Resource, error)

	ListNodes(context.Context) ([]*node.Resource, error)
	ListVolumes(context.Context, id.Namespace) ([]*volume.Resource, error)
	ListNamespaces(context.Context) ([]*namespace.Resource, error)

	CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error)
}

// Client provides a collection of methods for consumers to interact with the
// StorageOS API.
type Client struct {
	config    ConfigProvider
	transport Transport
}

// TODO: I think maybe this authenticate boiler plate should be moved down into
// the OpenAPI layer? That way we can be smart and avoid re-authing etcd without
// breaking abstraction layers.
func (c *Client) authenticate(ctx context.Context) (*user.Resource, error) {
	username, err := c.config.Username()
	if err != nil {
		return nil, err
	}
	password, err := c.config.Password()
	if err != nil {
		return nil, err
	}

	return c.transport.Authenticate(ctx, username, password)
}

// fetchAllVolumes requests the list of all namespaces from the StorageOS API,
// then requests the list of volumes within each namespace, returning an
// aggregate list of the volumes returned.
//
// If access is not granted when listing volumes for a retrieved namespace it
// is noted but will not return an error. Only if access is denied for all
// attempts will this return a permissions error.
func (c *Client) fetchAllVolumes(ctx context.Context) ([]*volume.Resource, error) {
	namespaces, err := c.transport.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	var volumes []*volume.Resource

	for _, ns := range namespaces {
		nsvols, err := c.transport.ListVolumes(ctx, ns.ID)
		switch {
		case err == nil, errors.Is(err, UnauthorisedError{}):
			// For these two errors, ignore - they're not fatal to the operation.
		default:
			return nil, err
		}
		volumes = append(volumes, nsvols...)
	}

	return volumes, nil
}

// New initialises a new Client using config for configuration settings,
// with transport providing the underlying implementation for encoding
// requests and decoding responses.
func New(transport Transport, config ConfigProvider) *Client {
	return &Client{
		transport: transport,
		config:    config,
	}
}
