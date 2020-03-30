// Package apiclient provides a type which implements an abstraction layer
// for consuming the storageos API programmatically.
package apiclient

import (
	"context"
	"errors"
	"io"
	"sync"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/policygroup"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/c2-cli/volume"
)

var (
	// ErrNoTransportConfigured indicates that the API client has not be
	// configured with an underlying transport implementation, which is required
	// for operation.
	ErrNoTransportConfigured = errors.New("the client has not been configured with a transport")
	// ErrTransportAlreadyConfigured indicates an attempt was made to configure
	// an API client with a new transport implementation when the client already
	// has one.
	//
	// To use a new transport implementation a consumer of the package must
	// instantiate a new Client.
	ErrTransportAlreadyConfigured = errors.New("the client's transport has already been configured")
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
	GetNode(ctx context.Context, nodeID id.Node) (*node.Resource, error)
	GetVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume) (*volume.Resource, error)
	GetNamespace(ctx context.Context, namespaceID id.Namespace) (*namespace.Resource, error)
	GetDiagnostics(ctx context.Context) (io.ReadCloser, error)

	ListNodes(ctx context.Context) ([]*node.Resource, error)
	ListVolumes(ctx context.Context, namespaceID id.Namespace) ([]*volume.Resource, error)
	ListNamespaces(ctx context.Context) ([]*namespace.Resource, error)
	ListPolicyGroups(ctx context.Context) ([]*policygroup.Resource, error)
	ListUsers(ctx context.Context) ([]*user.Resource, error)

	GetUser(ctx context.Context, username id.User) (*user.Resource, error)
	CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error)
	CreateVolume(ctx context.Context, namespaceID id.Namespace, name, description string, fs volume.FsType, sizeBytes uint64, labels map[string]string) (*volume.Resource, error)

	UpdateCluster(ctx context.Context, resource *cluster.Resource, licenceKey []byte) (*cluster.Resource, error)

	DeleteVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *DeleteVolumeRequestParams) error

	AttachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, nodeID id.Node) error
	DetachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *DetachVolumeRequestParams) error
}

// Client provides a collection of methods for consumers to interact with the
// StorageOS API.
type Client struct {
	config    ConfigProvider
	transport Transport

	configureOnce *sync.Once
}

// TODO(CP-4038): I think maybe this authenticate boiler plate should be moved
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

// ConfigureTransport must be called with a valid transport before any methods
// are called on c.
//
// The client's transport may only be set once through this method.
func (c *Client) ConfigureTransport(transport Transport) error {
	err := ErrTransportAlreadyConfigured

	c.configureOnce.Do(func() {
		c.transport = transport
		err = nil
	})

	return err
}

// New initialises a new Client which will source configuration settings using
// config. The returned client must be configured with a Transport before it is
// used.
func New(config ConfigProvider) *Client {
	return &Client{
		config:    config,
		transport: &noTransport{},

		configureOnce: &sync.Once{},
	}
}
