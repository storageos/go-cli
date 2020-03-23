package openapi

import (
	"context"

	"github.com/antihax/optional"

	"code.storageos.net/storageos/openapi"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
)

// CreateVolume requests the creation of a new volume through the StorageOS API
// using the provided parameters.
func (o *OpenAPI) CreateVolume(ctx context.Context, namespace id.Namespace, name, description string, fs volume.FsType, sizeBytes uint64, labels map[string]string) (*volume.Resource, error) {

	fsType, err := o.codec.encodeFsType(fs)
	if err != nil {
		return nil, err
	}

	createData := openapi.CreateVolumeData{
		NamespaceID: namespace.String(),
		Labels:      labels,
		Name:        name,
		FsType:      fsType,
		Description: description,
		SizeBytes:   sizeBytes,
	}

	// TODO(CP-3928): Creation of volumes can be done asynchronously, this should be
	// supported via setting the CreateVolumeOpts values when adding --async
	model, resp, err := o.client.DefaultApi.CreateVolume(ctx, namespace.String(), createData, nil)
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case badRequestError:
			return nil, apiclient.NewInvalidVolumeCreationError(v.msg)
		case conflictError:
			return nil, apiclient.NewVolumeExistsError(name, namespace)
		default:
			return nil, v
		}
	}

	return o.codec.decodeVolume(model)
}

// GetVolume requests the volume with uid from the StorageOS API, translating
// it into a *volume.Resource.
func (o *OpenAPI) GetVolume(ctx context.Context, namespaceID id.Namespace, uid id.Volume) (*volume.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, resp, err := o.client.DefaultApi.GetVolume(ctx, namespaceID.String(), uid.String())
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case notFoundError:
			return nil, apiclient.NewVolumeNotFoundError(v.msg)
		default:
			return nil, v
		}
	}

	return o.codec.decodeVolume(model)
}

// ListVolumes requests a list of all volume in namespace from the StorageOS
// API, translating each one to a *volume.Resource.
func (o *OpenAPI) ListVolumes(ctx context.Context, namespaceID id.Namespace) ([]*volume.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	models, resp, err := o.client.DefaultApi.ListVolumes(ctx, namespaceID.String())
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case notFoundError:
			return nil, apiclient.NewNamespaceNotFoundError(namespaceID)
		default:
			return nil, v
		}
	}

	volumes := make([]*volume.Resource, len(models))
	for i, m := range models {
		v, err := o.codec.decodeVolume(m)
		if err != nil {
			return nil, err
		}

		volumes[i] = v
	}

	return volumes, nil
}

// AttachVolume request to attach the volume `volumeID` in the namespace
// `namespaceID` to the node `nodeID`. It can return an error or nil if it
// succeeds.
func (o *OpenAPI) AttachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, nodeID id.Node) error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	resp, err := o.client.DefaultApi.AttachVolume(
		ctx,
		namespaceID.String(),
		volumeID.String(),
		openapi.AttachVolumeData{
			NodeID: nodeID.String(),
		},
	)
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case notFoundError:
			return apiclient.NewVolumeNotFoundError(v.msg)
		default:
			return v
		}
	}

	return nil
}

// DetachVolume makes a detach request for volumeID in namespaceID.
//
// The behaviour of the operation is dictated by params:
// 	- If params is nil or params.CASVersion is empty then the detach request is
// 	unconditional
// 	- If params.CASVersion is set, the request is conditional upon it matching
// 	the volume entity's version as seen by the server.
func (o *OpenAPI) DetachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *apiclient.DetachVolumeRequestParams) error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	var casVersion string
	var ignoreVersion bool = true

	// Set the CAS version constraint if provided
	if params != nil && params.CASVersion.String() != "" {
		ignoreVersion = false
		casVersion = params.CASVersion.String()
	}

	resp, err := o.client.DefaultApi.DetachVolume(
		ctx,
		namespaceID.String(),
		volumeID.String(),
		casVersion,
		&openapi.DetachVolumeOpts{
			IgnoreVersion: optional.NewBool(ignoreVersion),
		},
	)
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case notFoundError:
			return apiclient.NewVolumeNotFoundError(v.msg)
		default:
			return v
		}
	}

	return nil
}
