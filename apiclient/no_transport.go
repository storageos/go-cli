package apiclient

import (
	"context"
	"io"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/licence"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
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

func (t *noTransport) Authenticate(ctx context.Context, username, password string) (AuthSession, error) {
	return AuthSession{}, ErrNoTransportConfigured
}

func (t *noTransport) UseAuthSession(ctx context.Context, session AuthSession) error {
	return ErrNoTransportConfigured
}

func (t *noTransport) GetCluster(ctx context.Context) (*cluster.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) GetLicence(ctx context.Context) (*licence.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) GetDiagnostics(ctx context.Context) (io.ReadCloser, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) GetUser(ctx context.Context, username id.User) (*user.Resource, error) {
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

func (t *noTransport) GetPolicyGroup(ctx context.Context, uid id.PolicyGroup) (*policygroup.Resource, error) {
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

func (t *noTransport) ListUsers(ctx context.Context) ([]*user.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) CreateVolume(ctx context.Context, namespaceID id.Namespace, name, description string, fs volume.FsType, sizeBytes uint64, labels labels.Set, params *CreateVolumeRequestParams) (*volume.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) CreateNamespace(ctx context.Context, name string, labels labels.Set) (*namespace.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) CreatePolicyGroup(ctx context.Context, name string, specs []*policygroup.Spec) (*policygroup.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) UpdateCluster(ctx context.Context, resource *cluster.Resource) (*cluster.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) UpdateLicence(ctx context.Context, licence []byte, casVersion version.Version) (*licence.Resource, error) {
	return nil, ErrNoTransportConfigured
}

func (t *noTransport) DeleteUser(ctx context.Context, uid id.User, params *DeleteUserRequestParams) error {
	return ErrNoTransportConfigured
}

func (t *noTransport) DeleteVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *DeleteVolumeRequestParams) error {
	return ErrNoTransportConfigured
}

func (t *noTransport) DeleteNamespace(ctx context.Context, uid id.Namespace, params *DeleteNamespaceRequestParams) error {
	return ErrNoTransportConfigured
}

func (t *noTransport) DeletePolicyGroup(ctx context.Context, uid id.PolicyGroup, params *DeletePolicyGroupRequestParams) error {
	return ErrNoTransportConfigured
}

func (t *noTransport) AttachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, nodeID id.Node) error {
	return ErrNoTransportConfigured
}

func (t *noTransport) DetachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *DetachVolumeRequestParams) error {
	return ErrNoTransportConfigured
}
