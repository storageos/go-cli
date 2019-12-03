// Package apiclient provides a type which implements an abstraction layer
// for consuming the storageos API programmatically.
package apiclient

import (
	"context"

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

	GetCluster(ctx context.Context) (*cluster.Resource, error)
	GetNode(ctx context.Context, uid id.Node) (*node.Resource, error)
	GetVolume(ctx context.Context, namespaceID id.Namespace, uid id.Volume) (*volume.Resource, error)
	GetNamespace(ctx context.Context, uid id.Namespace) (*namespace.Resource, error)

	ListNodes(ctx context.Context) ([]*node.Resource, error)
	ListVolumes(ctx context.Context, namespaceID id.Namespace) ([]*volume.Resource, error)
	ListNamespaces(ctx context.Context) ([]*namespace.Resource, error)

	CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error)
}

// Client provides a collection of methods for consumers to interact with the
// StorageOS API.
type Client struct {
	config    ConfigProvider
	transport Transport
}

// TODO(CP-3930): I think maybe this authenticate boiler plate should be moved
// down into the OpenAPI layer? That way we can be smart and avoid re-authing
// etc without breaking abstraction layers. Marking this as part of the JWT
// caching work because it's related and might lead to a nice solution - Fraser
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

// New initialises a new Client using config for configuration settings,
// with transport providing the underlying implementation for encoding
// requests and decoding responses.
func New(transport Transport, config ConfigProvider) *Client {
	return &Client{
		transport: transport,
		config:    config,
	}
}
