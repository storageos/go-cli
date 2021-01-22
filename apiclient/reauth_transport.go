package apiclient

import (
	"context"
	"errors"

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

// CredentialsProvider defines a type which provides a configured username and
// password that can be used for authentication against the StorageOS API.
type CredentialsProvider interface {
	Username() (string, error)
	Password() (string, error)
}

// TransportWithReauth wraps calls to an inner transport implementation with a
// re-authenticate and retry mechanism when an authentication error is
// encountered.
type TransportWithReauth struct {
	inner Transport

	credentials CredentialsProvider
}

// Authenticate is passed through and does not try to reauth. An
// authentication error here cannot be due to a session timeout.
func (tr *TransportWithReauth) Authenticate(ctx context.Context, username, password string) (AuthSession, error) {
	return tr.inner.Authenticate(ctx, username, password)
}

// UseAuthSession is passed through and does not try to reauth.
func (tr *TransportWithReauth) UseAuthSession(ctx context.Context, session AuthSession) error {
	return tr.inner.UseAuthSession(ctx, session)
}

// GetUser wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) GetUser(ctx context.Context, uid id.User) (*user.Resource, error) {

	var resource *user.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resource, err = tr.inner.GetUser(ctx, uid)

		return err
	})

	return resource, err
}

// GetCluster wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) GetCluster(ctx context.Context) (*cluster.Resource, error) {

	var resource *cluster.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resource, err = tr.inner.GetCluster(ctx)

		return err
	})

	return resource, err
}

// GetLicence wraps the inner transport'call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) GetLicence(ctx context.Context) (*licence.Resource, error) {

	var resource *licence.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resource, err = tr.inner.GetLicence(ctx)

		return err
	})

	return resource, err
}

// GetNode wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) GetNode(ctx context.Context, nodeID id.Node) (*node.Resource, error) {

	var resource *node.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resource, err = tr.inner.GetNode(ctx, nodeID)

		return err
	})

	return resource, err
}

// GetVolume wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) GetVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume) (*volume.Resource, error) {

	var resource *volume.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resource, err = tr.inner.GetVolume(ctx, namespaceID, volumeID)

		return err
	})

	return resource, err
}

// GetNamespace wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) GetNamespace(ctx context.Context, namespaceID id.Namespace) (*namespace.Resource, error) {

	var resource *namespace.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resource, err = tr.inner.GetNamespace(ctx, namespaceID)

		return err
	})

	return resource, err
}

// GetDiagnostics wraps the inner transport's call with a reauthenticate and
// retry upon encountering an authentication error.
func (tr *TransportWithReauth) GetDiagnostics(ctx context.Context) (*diagnostics.BundleReadCloser, error) {

	var diagnostics *diagnostics.BundleReadCloser
	err := tr.doWithReauth(ctx, func() error {
		var err error
		diagnostics, err = tr.inner.GetDiagnostics(ctx)

		return err
	})

	return diagnostics, err
}

// GetPolicyGroup wraps the inner transport's call with a reauthenticate and
// retry upon encountering an authentication error.
func (tr *TransportWithReauth) GetPolicyGroup(ctx context.Context, uid id.PolicyGroup) (*policygroup.Resource, error) {

	var resource *policygroup.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resource, err = tr.inner.GetPolicyGroup(ctx, uid)

		return err
	})

	return resource, err
}

// ListNodes wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) ListNodes(ctx context.Context) ([]*node.Resource, error) {

	var resources []*node.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resources, err = tr.inner.ListNodes(ctx)

		return err
	})

	return resources, err
}

// ListVolumes wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) ListVolumes(ctx context.Context, namespaceID id.Namespace) ([]*volume.Resource, error) {

	var resources []*volume.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resources, err = tr.inner.ListVolumes(ctx, namespaceID)

		return err
	})

	return resources, err
}

// ListNamespaces wraps the inner transport's call with a reauthenticate and
// retry upon encountering an authentication error.
func (tr *TransportWithReauth) ListNamespaces(ctx context.Context) ([]*namespace.Resource, error) {

	var resources []*namespace.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resources, err = tr.inner.ListNamespaces(ctx)

		return err
	})

	return resources, err
}

// ListPolicyGroups wraps the inner transport's call with a reauthenticate and
// retry upon encountering an authentication error.
func (tr *TransportWithReauth) ListPolicyGroups(ctx context.Context) ([]*policygroup.Resource, error) {

	var resources []*policygroup.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resources, err = tr.inner.ListPolicyGroups(ctx)

		return err
	})

	return resources, err
}

// ListUsers wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) ListUsers(ctx context.Context) ([]*user.Resource, error) {

	var resources []*user.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resources, err = tr.inner.ListUsers(ctx)

		return err
	})

	return resources, err
}

// CreateUser wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error) {

	var resource *user.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resource, err = tr.inner.CreateUser(ctx, username, password, withAdmin, groups...)

		return err
	})

	return resource, err
}

// CreateVolume wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) CreateVolume(ctx context.Context, namespaceID id.Namespace, name, description string, fs volume.FsType, sizeBytes uint64, labels labels.Set, params *CreateVolumeRequestParams) (*volume.Resource, error) {

	var resource *volume.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resource, err = tr.inner.CreateVolume(
			ctx,
			namespaceID,
			name,
			description,
			fs,
			sizeBytes,
			labels,
			params,
		)

		return err
	})

	return resource, err
}

// CreatePolicyGroup wraps the inner transport's call with a reauthenticate
// and retry upon encountering an authentication error.
func (tr *TransportWithReauth) CreatePolicyGroup(ctx context.Context, name string, specs []*policygroup.Spec) (*policygroup.Resource, error) {

	var resource *policygroup.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resource, err = tr.inner.CreatePolicyGroup(ctx, name, specs)

		return err
	})

	return resource, err
}

// CreateNamespace wraps the inner transport's call with a reauthenticate and
// retry upon encountering an authentication error.
func (tr *TransportWithReauth) CreateNamespace(ctx context.Context, name string, labels labels.Set) (*namespace.Resource, error) {

	var resource *namespace.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		resource, err = tr.inner.CreateNamespace(ctx, name, labels)

		return err
	})

	return resource, err
}

// UpdateCluster wraps the inner transport's call with a reauthenticate and
// retry upon encountering an authentication error.
func (tr *TransportWithReauth) UpdateCluster(ctx context.Context, c *cluster.Resource, params *UpdateClusterRequestParams) (*cluster.Resource, error) {

	var updated *cluster.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		updated, err = tr.inner.UpdateCluster(ctx, c, params)

		return err
	})

	return updated, err
}

// UpdateLicence wraps the inner transport's call with a reauthenticate and
// retry upon encountering an authentication error.
func (tr *TransportWithReauth) UpdateLicence(ctx context.Context, lic []byte, params *UpdateLicenceRequestParams) (*licence.Resource, error) {

	var updated *licence.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		updated, err = tr.inner.UpdateLicence(ctx, lic, params)

		return err
	})

	return updated, err
}

// DeleteNode wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) DeleteNode(ctx context.Context, nodeID id.Node, params *DeleteNodeRequestParams) error {

	err := tr.doWithReauth(ctx, func() error {
		return tr.inner.DeleteNode(ctx, nodeID, params)
	})

	return err
}

// SetReplicas wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) SetReplicas(ctx context.Context, nsID id.Namespace, volID id.Volume, numReplicas uint64, params *SetReplicasRequestParams) error {

	err := tr.doWithReauth(ctx, func() error {
		return tr.inner.SetReplicas(ctx, nsID, volID, numReplicas, params)
	})

	return err
}

// UpdateVolume wraps the inner transport's call with a reauthenticate and
// retry upon encountering an authentication error.
func (tr *TransportWithReauth) UpdateVolume(ctx context.Context, nsID id.Namespace, volID id.Volume, description string, labels labels.Set, params *UpdateVolumeRequestParams) (*volume.Resource, error) {

	var updated *volume.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		updated, err = tr.inner.UpdateVolume(ctx, nsID, volID, description, labels, params)
		return err
	})

	return updated, err
}

// ResizeVolume wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) ResizeVolume(ctx context.Context, nsID id.Namespace, volID id.Volume, sizeBytes uint64, params *ResizeVolumeRequestParams) (*volume.Resource, error) {

	var updated *volume.Resource
	err := tr.doWithReauth(ctx, func() error {
		var err error
		updated, err = tr.inner.ResizeVolume(ctx, nsID, volID, sizeBytes, params)
		return err
	})

	return updated, err
}

// UpdateNFSVolumeExports wraps the inner transport's call with a reauthenticate
// and retry upon encountering an authentication error.
func (tr *TransportWithReauth) UpdateNFSVolumeExports(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, exports []volume.NFSExportConfig, params *UpdateNFSVolumeExportsRequestParams) error {

	err := tr.doWithReauth(ctx, func() error {
		return tr.inner.UpdateNFSVolumeExports(ctx, namespaceID, volumeID, exports, params)
	})

	return err
}

// UpdateNFSVolumeMountEndpoint wraps the inner transport's call with a
// reauthenticate and retry upon encountering an authentication error.
func (tr *TransportWithReauth) UpdateNFSVolumeMountEndpoint(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, endpoint string, params *UpdateNFSVolumeMountEndpointRequestParams) error {

	err := tr.doWithReauth(ctx, func() error {
		return tr.inner.UpdateNFSVolumeMountEndpoint(ctx, namespaceID, volumeID, endpoint, params)
	})

	return err
}

// DeleteVolume wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) DeleteVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *DeleteVolumeRequestParams) error {

	err := tr.doWithReauth(ctx, func() error {
		return tr.inner.DeleteVolume(ctx, namespaceID, volumeID, params)
	})

	return err
}

// DeleteNamespace wraps the inner transport's call with a reauthenticate and
// retry upon encountering an authentication error.
func (tr *TransportWithReauth) DeleteNamespace(ctx context.Context, uid id.Namespace, params *DeleteNamespaceRequestParams) error {

	err := tr.doWithReauth(ctx, func() error {
		return tr.inner.DeleteNamespace(ctx, uid, params)
	})

	return err
}

// DeleteUser wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) DeleteUser(ctx context.Context, uid id.User, params *DeleteUserRequestParams) error {

	err := tr.doWithReauth(ctx, func() error {
		return tr.inner.DeleteUser(ctx, uid, params)
	})

	return err
}

// DeletePolicyGroup wraps the inner transport's call with a reauthenticate and
// retry upon encountering an authentication error.
func (tr *TransportWithReauth) DeletePolicyGroup(ctx context.Context, uid id.PolicyGroup, params *DeletePolicyGroupRequestParams) error {

	err := tr.doWithReauth(ctx, func() error {
		return tr.inner.DeletePolicyGroup(ctx, uid, params)
	})

	return err
}

// AttachVolume wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) AttachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, nodeID id.Node) error {

	err := tr.doWithReauth(ctx, func() error {
		return tr.inner.AttachVolume(ctx, namespaceID, volumeID, nodeID)
	})

	return err
}

// AttachNFSVolume wraps the inner transport's call with a reauthenticate and
// retry upon encountering an authentication error.
func (tr *TransportWithReauth) AttachNFSVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *AttachNFSVolumeRequestParams) error {

	err := tr.doWithReauth(ctx, func() error {
		return tr.inner.AttachNFSVolume(ctx, namespaceID, volumeID, params)
	})

	return err
}

// DetachVolume wraps the inner transport's call with a reauthenticate and retry
// upon encountering an authentication error.
func (tr *TransportWithReauth) DetachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *DetachVolumeRequestParams) error {

	err := tr.doWithReauth(ctx, func() error {
		return tr.inner.DetachVolume(ctx, namespaceID, volumeID, params)
	})

	return err
}

// doWithReauth invokes fn, checking the resultant error.
//
//  - If the error is an *AuthenticationError then tr's credentials are
// 	used to reauthenticate before returning the result from re-invoking fn.
//  If any errors occur during reauthentication, they are returned.
//  - Otherwise, the original error is returned to the caller.
func (tr *TransportWithReauth) doWithReauth(ctx context.Context, fn func() error) error {
	originalErr := fn()

	// If the returned error from fn indicates authentication failure then
	// fetch credentials from the provider, reauthenticate and try the request
	// one more time.
	//
	// This will reliably catch a cached auth session expiring.
	if errors.As(originalErr, &AuthenticationError{}) {

		username, err := tr.credentials.Username()
		if err != nil {
			return err
		}
		password, err := tr.credentials.Password()
		if err != nil {
			return err
		}

		// Attempt to reauth with credentials from the provider.
		_, err = tr.Authenticate(ctx, username, password)
		if err != nil {
			return err
		}

		return fn()
	}

	return originalErr
}

// NewTransportWithReauth wraps calls to transport with a retry on
// authentication failure, sourcing username and password from credentials.
func NewTransportWithReauth(transport Transport, credentials CredentialsProvider) *TransportWithReauth {
	return &TransportWithReauth{
		inner:       transport,
		credentials: credentials,
	}
}
