package apiclient

import (
	"context"
	"io"
	"sync"

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

type mockTransport struct {
	mu sync.RWMutex

	AuthenticateGotUsername string
	AuthenticateGotPassword string
	AuthenticateAuthSession AuthSession
	AuthenticateError       error

	UseAuthSessionGotAuthSession AuthSession
	UseAuthSessionError          error

	GetClusterResource *cluster.Resource
	GetClusterError    error

	GetLicenceResource *licence.Resource
	GetLicenceError    error

	GetUserGotID    id.User
	GetUserResource *user.Resource
	GetUserError    error

	GetDiagnosticsReadCloser io.ReadCloser
	GetDiagnosticsError      error

	GetNodeGotID    id.Node
	GetNodeResource *node.Resource
	GetNodeError    error

	GetVolumeGotNamespaceID id.Namespace
	GetVolumeGotVolumeID    id.Volume
	GetVolumeResource       *volume.Resource
	GetVolumeError          error

	GetNamespaceGotID    id.Namespace
	GetNamespaceResource *namespace.Resource
	GetNamespaceError    error

	GetPolicyGroupGotID    id.PolicyGroup
	GetPolicyGroupResource *policygroup.Resource
	GetPolicyGroupError    error

	ListNodesResource []*node.Resource
	ListNodesError    error

	ListVolumesGotNamespaceIDs []id.Namespace
	ListVolumesResource        map[id.Namespace][]*volume.Resource
	ListVolumesError           error

	ListNamespacesResource []*namespace.Resource
	ListNamespacesError    error

	ListPolicyGroupsResource []*policygroup.Resource
	ListPolicyGroupsError    error

	ListUsersResource []*user.Resource
	ListUserError     error

	CreateUserGotName     string
	CreateUserGotPassword string
	CreateUserGotAdmin    bool
	CreateUserGotGroups   []id.PolicyGroup
	CreateUserResource    *user.Resource
	CreateUserError       error

	CreateVolumeGotNamespace   id.Namespace
	CreateVolumeGotName        string
	CreateVolumeGotDescription string
	CreateVolumeGotFs          volume.FsType
	CreateVolumeGotSizeBytes   uint64
	CreateVolumeGotLabels      labels.Set
	CreateVolumeGotParams      *CreateVolumeRequestParams
	CreateVolumeResource       *volume.Resource
	CreateVolumeError          error

	CreateNamespaceGotName   string
	CreateNamespaceGotLabels labels.Set
	CreateNamespaceResource  *namespace.Resource
	CreateNamespaceError     error

	CreatePolicyGroupGotName  string
	CreatePolicyGroupGotSpecs []*policygroup.Spec
	CreatePolicyGroupResource *policygroup.Resource
	CreatePolicyGroupError    error

	UpdateClusterResource    *cluster.Resource
	UpdateClusterError       error
	UpdateClusterGotResource *cluster.Resource

	UpdateLicenceGotLicence []byte
	UpdateLicenceGotVersion version.Version
	UpdateLicenceResource   *licence.Resource
	UpdateLicenceError      error

	SetReplicasGotNamespaceID id.Namespace
	SetReplicasGotVolumeID    id.Volume
	SetReplicasGotNumReplicas uint64
	SetReplicasGotVersion     version.Version
	SetReplicasError          error

	UpdateVolumeGotNamespaceID id.Namespace
	UpdateVolumeGotVolumeID    id.Volume
	UpdateVolumeGotDescription string
	UpdateVolumeGotLabels      labels.Set
	UpdateVolumeGotVersion     version.Version
	UpdateVolumeResource       *volume.Resource
	UpdateVolumeError          error

	DeleteUserGotID     id.User
	DeleteUserGotParams *DeleteUserRequestParams
	DeleteUserError     error

	DeleteVolumeGotNamespace id.Namespace
	DeleteVolumeGotVolume    id.Volume
	DeleteVolumeGotParams    *DeleteVolumeRequestParams
	DeleteVolumeError        error

	DeleteNodeGotNode   id.Node
	DeleteNodeGotParams *DeleteNodeRequestParams
	DeleteNodeError     error

	DeleteNamespaceGotID     id.Namespace
	DeleteNamespaceGotParams *DeleteNamespaceRequestParams
	DeleteNamespaceError     error

	DeletePolicyGroupGotID     id.PolicyGroup
	DeletePolicyGroupGotParams *DeletePolicyGroupRequestParams
	DeletePolicyGroupError     error

	AttachGotNamespace id.Namespace
	AttachGotVolume    id.Volume
	AttachGotNode      id.Node
	AttachError        error

	DetachGotNamespace id.Namespace
	DetachGotVolume    id.Volume
	DetachGotParams    *DetachVolumeRequestParams
	DetachError        error
}

var _ Transport = (*mockTransport)(nil)

func (m *mockTransport) Authenticate(ctx context.Context, username, password string) (AuthSession, error) {
	m.AuthenticateGotUsername = username
	m.AuthenticateGotPassword = password
	return m.AuthenticateAuthSession, m.AuthenticateError
}

func (m *mockTransport) UseAuthSession(ctx context.Context, session AuthSession) error {
	m.UseAuthSessionGotAuthSession = session
	return m.UseAuthSessionError
}

func (m *mockTransport) GetCluster(ctx context.Context) (*cluster.Resource, error) {
	return m.GetClusterResource, m.GetClusterError
}

func (m *mockTransport) GetLicence(ctx context.Context) (*licence.Resource, error) {
	return m.GetLicenceResource, m.GetLicenceError
}

func (m *mockTransport) GetUser(ctx context.Context, uid id.User) (*user.Resource, error) {
	m.GetUserGotID = uid
	return m.GetUserResource, m.GetUserError
}

func (m *mockTransport) ListUsers(ctx context.Context) ([]*user.Resource, error) {
	return m.ListUsersResource, m.ListUserError
}

func (m *mockTransport) GetDiagnostics(ctx context.Context) (io.ReadCloser, error) {
	return m.GetDiagnosticsReadCloser, m.GetDiagnosticsError
}

func (m *mockTransport) GetNode(ctx context.Context, nodeID id.Node) (*node.Resource, error) {
	m.GetNodeGotID = nodeID
	return m.GetNodeResource, m.GetNodeError
}

func (m *mockTransport) GetVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume) (*volume.Resource, error) {
	m.GetVolumeGotNamespaceID = namespaceID
	m.GetVolumeGotVolumeID = volumeID
	return m.GetVolumeResource, m.GetVolumeError
}

func (m *mockTransport) GetNamespace(ctx context.Context, namespaceID id.Namespace) (*namespace.Resource, error) {
	m.GetNamespaceGotID = namespaceID
	return m.GetNamespaceResource, m.GetNamespaceError
}

func (m *mockTransport) GetPolicyGroup(ctx context.Context, uid id.PolicyGroup) (*policygroup.Resource, error) {
	m.GetPolicyGroupGotID = uid
	return m.GetPolicyGroupResource, m.GetPolicyGroupError
}

func (m *mockTransport) ListNodes(ctx context.Context) ([]*node.Resource, error) {
	return m.ListNodesResource, m.ListNodesError
}

func (m *mockTransport) ListVolumes(ctx context.Context, namespaceID id.Namespace) ([]*volume.Resource, error) {
	m.mu.Lock()
	m.ListVolumesGotNamespaceIDs = append(m.ListVolumesGotNamespaceIDs, namespaceID)
	m.mu.Unlock()
	return m.ListVolumesResource[namespaceID], m.ListVolumesError
}

func (m *mockTransport) ListNamespaces(ctx context.Context) ([]*namespace.Resource, error) {
	return m.ListNamespacesResource, m.ListNamespacesError
}

func (m *mockTransport) ListPolicyGroups(ctx context.Context) ([]*policygroup.Resource, error) {
	return m.ListPolicyGroupsResource, m.ListPolicyGroupsError
}

func (m *mockTransport) CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error) {
	m.CreateUserGotName = username
	m.CreateUserGotPassword = password
	m.CreateUserGotAdmin = withAdmin
	m.CreateUserGotGroups = groups
	return m.CreateUserResource, m.CreateUserError
}

func (m *mockTransport) CreateVolume(ctx context.Context, namespaceID id.Namespace, name, description string, fs volume.FsType, sizeBytes uint64, labels labels.Set, params *CreateVolumeRequestParams) (*volume.Resource, error) {
	m.CreateVolumeGotNamespace = namespaceID
	m.CreateVolumeGotName = name
	m.CreateVolumeGotDescription = description
	m.CreateVolumeGotFs = fs
	m.CreateVolumeGotSizeBytes = sizeBytes
	m.CreateVolumeGotLabels = labels
	m.CreateVolumeGotParams = params
	return m.CreateVolumeResource, m.CreateVolumeError
}

func (m *mockTransport) CreateNamespace(ctx context.Context, name string, labels labels.Set) (*namespace.Resource, error) {
	m.CreateNamespaceGotName = name
	m.CreateNamespaceGotLabels = labels
	return m.CreateNamespaceResource, m.CreateNamespaceError
}

func (m *mockTransport) CreatePolicyGroup(ctx context.Context, name string, specs []*policygroup.Spec) (*policygroup.Resource, error) {
	m.CreatePolicyGroupGotName = name
	m.CreatePolicyGroupGotSpecs = specs
	return m.CreatePolicyGroupResource, m.CreatePolicyGroupError
}

func (m *mockTransport) UpdateCluster(ctx context.Context, resource *cluster.Resource) (*cluster.Resource, error) {
	m.UpdateClusterGotResource = resource
	return m.UpdateClusterResource, m.UpdateClusterError
}

func (m *mockTransport) SetReplicas(ctx context.Context, nsID id.Namespace, volID id.Volume, numReplicas uint64, version version.Version) error {
	m.SetReplicasGotNamespaceID = nsID
	m.SetReplicasGotVolumeID = volID
	m.SetReplicasGotNumReplicas = numReplicas
	m.SetReplicasGotVersion = version
	return m.SetReplicasError
}

func (m *mockTransport) UpdateVolume(ctx context.Context, nsID id.Namespace, volID id.Volume, description string, labels labels.Set, version version.Version) (*volume.Resource, error) {
	m.UpdateVolumeGotNamespaceID = nsID
	m.UpdateVolumeGotVolumeID = volID
	m.UpdateVolumeGotDescription = description
	m.UpdateVolumeGotLabels = labels
	m.UpdateVolumeGotVersion = version
	return m.UpdateVolumeResource, m.UpdateVolumeError
}

func (m *mockTransport) UpdateLicence(ctx context.Context, licence []byte, casVersion version.Version) (*licence.Resource, error) {
	m.UpdateLicenceGotLicence = licence
	m.UpdateLicenceGotVersion = casVersion
	return m.UpdateLicenceResource, m.UpdateLicenceError
}

func (m *mockTransport) DeleteUser(ctx context.Context, uid id.User, params *DeleteUserRequestParams) error {
	m.DeleteUserGotID = uid
	m.DeleteUserGotParams = params
	return m.DeleteUserError
}

func (m *mockTransport) DeleteVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *DeleteVolumeRequestParams) error {
	m.DeleteVolumeGotNamespace = namespaceID
	m.DeleteVolumeGotVolume = volumeID
	m.DeleteVolumeGotParams = params
	return m.DeleteVolumeError
}

func (m *mockTransport) DeleteNode(ctx context.Context, volumeID id.Node, params *DeleteNodeRequestParams) error {
	m.DeleteNodeGotNode = volumeID
	m.DeleteNodeGotParams = params
	return m.DeleteNodeError
}

func (m *mockTransport) DeleteNamespace(ctx context.Context, uid id.Namespace, params *DeleteNamespaceRequestParams) error {
	m.DeleteNamespaceGotID = uid
	m.DeleteNamespaceGotParams = params
	return m.DeleteNamespaceError
}

func (m *mockTransport) DeletePolicyGroup(ctx context.Context, uid id.PolicyGroup, params *DeletePolicyGroupRequestParams) error {
	m.DeletePolicyGroupGotID = uid
	m.DeletePolicyGroupGotParams = params
	return m.DeletePolicyGroupError
}

func (m *mockTransport) AttachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, nodeID id.Node) error {
	m.AttachGotNamespace = namespaceID
	m.AttachGotVolume = volumeID
	m.AttachGotNode = nodeID
	return m.AttachError
}

func (m *mockTransport) DetachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *DetachVolumeRequestParams) error {
	m.DetachGotNamespace = namespaceID
	m.DetachGotVolume = volumeID
	m.DetachGotParams = params
	return m.DetachError
}
