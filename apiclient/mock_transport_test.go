package apiclient

import (
	"context"
	"sync"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/diagnostics"
	"code.storageos.net/storageos/c2-cli/licence"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
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

	GetDiagnosticsReadCloser *diagnostics.BundleReadCloser
	GetDiagnosticsError      error

	GetSingleNodeDiagnosticsReadCloser *diagnostics.BundleReadCloser
	GetSingleNodeDiagnosticsError      error

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
	UpdateClusterGotParams   *UpdateClusterRequestParams

	UpdateLicenceGotLicence []byte
	UpdateLicenceGotParams  *UpdateLicenceRequestParams
	UpdateLicenceResource   *licence.Resource
	UpdateLicenceError      error

	SetReplicasGotNamespaceID id.Namespace
	SetReplicasGotVolumeID    id.Volume
	SetReplicasGotNumReplicas uint64
	SetReplicasGotParams      *SetReplicasRequestParams
	SetReplicasError          error

	UpdateVolumeGotNamespaceID id.Namespace
	UpdateVolumeGotVolumeID    id.Volume
	UpdateVolumeGotDescription string
	UpdateVolumeGotLabels      labels.Set
	UpdateVolumeGotParams      *UpdateVolumeRequestParams
	UpdateVolumeResource       *volume.Resource
	UpdateVolumeError          error

	ResizeVolumeGotNamespaceID id.Namespace
	ResizeVolumeGotVolumeID    id.Volume
	ResizeVolumeGotSizeBytes   uint64
	ResizeVolumeGotParams      *ResizeVolumeRequestParams
	ResizeVolumeResource       *volume.Resource
	ResizeVolumeError          error

	UpdateNFSVolumeExportsGotNamespaceID id.Namespace
	UpdateNFSVolumeExportsGotVolumeID    id.Volume
	UpdateNFSVolumeExportsGotExports     []volume.NFSExportConfig
	UpdateNFSVolumeExportsGotParams      *UpdateNFSVolumeExportsRequestParams
	UpdateNFSVolumeExportsError          error

	UpdateNFSVolumeMountEndpointGotNamespaceID id.Namespace
	UpdateNFSVolumeMountEndpointGotVolumeID    id.Volume
	UpdateNFSVolumeMountEndpointGotEndpoint    string
	UpdateNFSVolumeMountEndpointGotParams      *UpdateNFSVolumeMountEndpointRequestParams
	UpdateNFSVolumeMountEndpointError          error

	SetFailureModeIntentGotNamespaceID id.Namespace
	SetFailureModeIntentGotVolumeID    id.Volume
	SetFailureModeIntentGotIntent      string
	SetFailureModeIntentGotParams      *SetFailureModeRequestParams
	SetFailureModeIntentResource       *volume.Resource
	SetFailureModeIntentError          error

	SetFailureThresholdGotNamespaceID id.Namespace
	SetFailureThresholdGotVolumeID    id.Volume
	SetFailureThresholdGotThreshold   uint64
	SetFailureThresholdGotParams      *SetFailureModeRequestParams
	SetFailureThresholdResource       *volume.Resource
	SetFailureThresholdError          error

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

	AttachNFSGotNamespace id.Namespace
	AttachNFSGotVolume    id.Volume
	AttachNFSGotParams    *AttachNFSVolumeRequestParams
	AttachNFSError        error

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

func (m *mockTransport) GetDiagnostics(ctx context.Context) (*diagnostics.BundleReadCloser, error) {
	return m.GetDiagnosticsReadCloser, m.GetDiagnosticsError
}

func (m *mockTransport) GetSingleNodeDiagnostics(ctx context.Context, nodeID id.Node) (*diagnostics.BundleReadCloser, error) {
	return m.GetSingleNodeDiagnosticsReadCloser, m.GetSingleNodeDiagnosticsError
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

func (m *mockTransport) UpdateCluster(ctx context.Context, resource *cluster.Resource, params *UpdateClusterRequestParams) (*cluster.Resource, error) {
	m.UpdateClusterGotResource = resource
	m.UpdateClusterGotParams = params
	return m.UpdateClusterResource, m.UpdateClusterError
}

func (m *mockTransport) SetReplicas(ctx context.Context, nsID id.Namespace, volID id.Volume, numReplicas uint64, params *SetReplicasRequestParams) error {
	m.SetReplicasGotNamespaceID = nsID
	m.SetReplicasGotVolumeID = volID
	m.SetReplicasGotNumReplicas = numReplicas
	m.SetReplicasGotParams = params
	return m.SetReplicasError
}

func (m *mockTransport) UpdateVolume(ctx context.Context, nsID id.Namespace, volID id.Volume, description string, labels labels.Set, params *UpdateVolumeRequestParams) (*volume.Resource, error) {
	m.UpdateVolumeGotNamespaceID = nsID
	m.UpdateVolumeGotVolumeID = volID
	m.UpdateVolumeGotDescription = description
	m.UpdateVolumeGotLabels = labels
	m.UpdateVolumeGotParams = params
	return m.UpdateVolumeResource, m.UpdateVolumeError
}

func (m *mockTransport) ResizeVolume(ctx context.Context, nsID id.Namespace, volID id.Volume, sizeBytes uint64, params *ResizeVolumeRequestParams) (*volume.Resource, error) {
	m.ResizeVolumeGotNamespaceID = nsID
	m.ResizeVolumeGotVolumeID = volID
	m.ResizeVolumeGotSizeBytes = sizeBytes
	m.ResizeVolumeGotParams = params
	return m.ResizeVolumeResource, m.ResizeVolumeError
}

func (m *mockTransport) UpdateNFSVolumeExports(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, exports []volume.NFSExportConfig, params *UpdateNFSVolumeExportsRequestParams) error {
	m.UpdateNFSVolumeExportsGotNamespaceID = namespaceID
	m.UpdateNFSVolumeExportsGotVolumeID = volumeID
	m.UpdateNFSVolumeExportsGotExports = exports
	m.UpdateNFSVolumeExportsGotParams = params
	return m.UpdateNFSVolumeExportsError
}

func (m *mockTransport) UpdateNFSVolumeMountEndpoint(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, endpoint string, params *UpdateNFSVolumeMountEndpointRequestParams) error {
	m.UpdateNFSVolumeMountEndpointGotNamespaceID = namespaceID
	m.UpdateNFSVolumeMountEndpointGotVolumeID = volumeID
	m.UpdateNFSVolumeMountEndpointGotEndpoint = endpoint
	m.UpdateNFSVolumeMountEndpointGotParams = params
	return m.UpdateNFSVolumeMountEndpointError
}

func (m *mockTransport) SetFailureModeIntent(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, intent string, params *SetFailureModeRequestParams) (*volume.Resource, error) {
	m.SetFailureModeIntentGotNamespaceID = namespaceID
	m.SetFailureModeIntentGotVolumeID = volumeID
	m.SetFailureModeIntentGotIntent = intent
	m.SetFailureModeIntentGotParams = params
	return m.SetFailureModeIntentResource, m.SetFailureModeIntentError
}

func (m *mockTransport) SetFailureThreshold(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, threshold uint64, params *SetFailureModeRequestParams) (*volume.Resource, error) {
	m.SetFailureThresholdGotNamespaceID = namespaceID
	m.SetFailureThresholdGotVolumeID = volumeID
	m.SetFailureThresholdGotThreshold = threshold
	m.SetFailureThresholdGotParams = params
	return m.SetFailureThresholdResource, m.SetFailureThresholdError
}

func (m *mockTransport) UpdateLicence(ctx context.Context, licence []byte, params *UpdateLicenceRequestParams) (*licence.Resource, error) {
	m.UpdateLicenceGotLicence = licence
	m.UpdateLicenceGotParams = params
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

func (m *mockTransport) AttachNFSVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *AttachNFSVolumeRequestParams) error {
	m.AttachNFSGotNamespace = namespaceID
	m.AttachNFSGotVolume = volumeID
	m.AttachNFSGotParams = params
	return m.AttachNFSError
}

func (m *mockTransport) DetachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *DetachVolumeRequestParams) error {
	m.DetachGotNamespace = namespaceID
	m.DetachGotVolume = volumeID
	m.DetachGotParams = params
	return m.DetachError
}
