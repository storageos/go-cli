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

type mockTransport struct {
	AuthenticateUserResource *user.Resource
	AuthenticateError        error

	GetClusterResource *cluster.Resource
	GetClusterError    error

	GetNodeResource *node.Resource
	GetNodeError    error

	GetVolumeResource *volume.Resource
	GetVolumeError    error

	GetNamespaceResource *namespace.Resource
	GetNamespaceError    error

	ListNodesResource []*node.Resource
	ListNodesError    error

	ListVolumesResource []*volume.Resource
	ListVolumesError    error

	ListNamespacesResource []*namespace.Resource
	ListNamespacesError    error

	ListPolicyGroupsResource []*policygroup.Resource
	ListPolicyGroupsError    error

	CreateUserResource *user.Resource
	CreateUserError    error

	CreateVolumeResource *volume.Resource
	CreateVolumeError    error

	UpdateClusterResource      *cluster.Resource
	UpdateClusterError         error
	UpdateClusterGotResource   *cluster.Resource
	UpdateClusterGotLicenceKey []byte

	AttachGotNamespace id.Namespace
	AttachGotVolume    id.Volume
	AttachGotNode      id.Node
	AttachError        error
}

var _ Transport = (*mockTransport)(nil)

func (m *mockTransport) Authenticate(ctx context.Context, username, password string) (*user.Resource, error) {
	return m.AuthenticateUserResource, m.AuthenticateError
}

func (m *mockTransport) GetCluster(ctx context.Context) (*cluster.Resource, error) {
	return m.GetClusterResource, m.GetClusterError
}

func (m *mockTransport) GetNode(ctx context.Context, nodeID id.Node) (*node.Resource, error) {
	return m.GetNodeResource, m.GetNodeError
}

func (m *mockTransport) GetVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume) (*volume.Resource, error) {
	return m.GetVolumeResource, m.GetVolumeError
}

func (m *mockTransport) GetNamespace(ctx context.Context, namespaceID id.Namespace) (*namespace.Resource, error) {
	return m.GetNamespaceResource, m.GetNamespaceError
}

func (m *mockTransport) ListNodes(ctx context.Context) ([]*node.Resource, error) {
	return m.ListNodesResource, m.ListNodesError
}

func (m *mockTransport) ListVolumes(ctx context.Context, namespaceID id.Namespace) ([]*volume.Resource, error) {
	return m.ListVolumesResource, m.ListVolumesError
}

func (m *mockTransport) ListNamespaces(ctx context.Context) ([]*namespace.Resource, error) {
	return m.ListNamespacesResource, m.ListNamespacesError
}

func (m *mockTransport) ListPolicyGroups(ctx context.Context) ([]*policygroup.Resource, error) {
	return m.ListPolicyGroupsResource, m.ListPolicyGroupsError
}

func (m *mockTransport) CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error) {
	return m.CreateUserResource, m.CreateUserError
}

func (m *mockTransport) CreateVolume(ctx context.Context, namespaceID id.Namespace, name, description string, fs volume.FsType, sizeBytes uint64, labels map[string]string) (*volume.Resource, error) {
	return m.CreateVolumeResource, m.CreateVolumeError
}

func (m *mockTransport) UpdateCluster(ctx context.Context, resource *cluster.Resource, licenceKey []byte) (*cluster.Resource, error) {
	m.UpdateClusterGotResource = resource
	m.UpdateClusterGotLicenceKey = licenceKey
	return m.UpdateClusterResource, m.UpdateClusterError
}

func (m *mockTransport) AttachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, nodeID id.Node) error {
	m.AttachGotNamespace = namespaceID
	m.AttachGotVolume = volumeID
	m.AttachGotNode = nodeID
	return m.AttachError
}
