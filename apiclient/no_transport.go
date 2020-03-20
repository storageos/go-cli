package apiclient

import (
	"context"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/policygroup"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/c2-cli/volume"
)

// noTransport implements the Transport interface for the API client, but
// returns a known error from all method invocations. It is a placeholder
// to allow for configuration of the API client to take place after it
// has been constructed.
type noTransport struct{}

var _ Transport = (*noTransport)(nil)

func (t *noTransport) Authenticate(ctx context.Context, username, password string) (*user.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) GetCluster(ctx context.Context) (*cluster.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) GetNode(ctx context.Context, nodeID id.Node) (*node.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) GetVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume) (*volume.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) GetNamespace(ctx context.Context, namespaceID id.Namespace) (*namespace.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) ListNodes(ctx context.Context) ([]*node.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) ListVolumes(ctx context.Context, namespaceID id.Namespace) ([]*volume.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) ListNamespaces(ctx context.Context) ([]*namespace.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) ListPolicyGroups(ctx context.Context) ([]*policygroup.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) CreateVolume(ctx context.Context, namespaceID id.Namespace, name, description string, fs volume.FsType, sizeBytes uint64, labels map[string]string) (*volume.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) UpdateCluster(ctx context.Context, resource *cluster.Resource, licenceKey []byte) (*cluster.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) AttachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, nodeID id.Node) error {
	return nil
}
