package openapi

import (
	"context"

	"github.com/antihax/optional"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/volume"
	"code.storageos.net/storageos/openapi"
)

// CreateVolume requests the creation of a new volume through the StorageOS API
// using the provided parameters.
//
// The behaviour of the operation is dictated by params:
//
//
//  Asynchrony:
//  - If params is nil or params.AsyncMax is empty/zero valued then the create
//  request is performed synchronously.
//  - If params.AsyncMax is set, the request is performed asynchronously using
//  the duration given as the maximum amount of time allowed for the request
//  before it times out.
func (o *OpenAPI) CreateVolume(ctx context.Context, namespace id.Namespace, name, description string, fs volume.FsType, sizeBytes uint64, labels labels.Set, params *apiclient.CreateVolumeRequestParams) (*volume.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	fsType, err := o.codec.encodeFsType(fs)
	if err != nil {
		return nil, err
	}

	var asyncMax optional.String = optional.EmptyString()

	if params != nil {
		if params.AsyncMax != 0 {
			asyncMax = optional.NewString(params.AsyncMax.String())
		}
	}

	createData := openapi.CreateVolumeData{
		NamespaceID: namespace.String(),
		Labels:      labels,
		Name:        name,
		FsType:      fsType,
		Description: description,
		SizeBytes:   sizeBytes,
	}

	model, resp, err := o.client.DefaultApi.CreateVolume(
		ctx,
		namespace.String(),
		createData,
		&openapi.CreateVolumeOpts{
			AsyncMax: asyncMax,
		},
	)
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

// DeleteVolume makes a delete request for volumeID in namespaceID.
//
// The behaviour of the operation is dictated by params:
//
//
// 	Version constraints:
// 	- If params is nil or params.CASVersion is empty then the delete request is
// 	unconditional
// 	- If params.CASVersion is set, the request is conditional upon it matching
// 	the volume entity's version as seen by the server.
//
//  Asynchrony:
//  - If params is nil or params.AsyncMax is empty/zero valued then the delete
//  request is performed synchronously.
//  - If params.AsyncMax is set, the request is performed asynchronously using
//  the duration given as the maximum amount of time allowed for the request
//  before it times out.
//
//  Offline deletion behaviour:
//  - If params is nil then offline deletion behaviour is not requested,
//  otherwise the value of params.OfflineDelete determines is used. The default
//  value of false reflects normal deletion behaviour, so does not need setting
//  unless offline deletion behaviour is desired.
func (o *OpenAPI) DeleteVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *apiclient.DeleteVolumeRequestParams) error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	var casVersion string
	var ignoreVersion optional.Bool = optional.NewBool(true)
	var asyncMax optional.String = optional.EmptyString()
	var offlineDelete optional.Bool = optional.NewBool(false)

	if params != nil {
		if params.CASVersion.String() != "" {
			ignoreVersion = optional.NewBool(false)
			casVersion = params.CASVersion.String()
		}

		if params.AsyncMax != 0 {
			asyncMax = optional.NewString(params.AsyncMax.String())
		}

		offlineDelete = optional.NewBool(params.OfflineDelete)
	}

	resp, err := o.client.DefaultApi.DeleteVolume(
		ctx,
		namespaceID.String(),
		volumeID.String(),
		casVersion,
		&openapi.DeleteVolumeOpts{
			IgnoreVersion: ignoreVersion,
			AsyncMax:      asyncMax,
			OfflineDelete: offlineDelete,
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

// AttachNFSVolume request to attach the volume `volumeID` in the namespace
// `namespaceID` for a NFS volume. It can return an error or nil if it
// succeeds.
func (o *OpenAPI) AttachNFSVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *apiclient.AttachNFSVolumeRequestParams) error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	var casVersion string
	var ignoreVersion optional.Bool = optional.NewBool(true)
	var asyncMax optional.String = optional.EmptyString()

	if params != nil {
		if params.CASVersion.String() != "" {
			ignoreVersion = optional.NewBool(false)
			casVersion = params.CASVersion.String()
		}

		if params.AsyncMax != 0 {
			asyncMax = optional.NewString(params.AsyncMax.String())
		}
	}

	resp, err := o.client.DefaultApi.AttachNFSVolume(
		ctx,
		namespaceID.String(),
		volumeID.String(),
		openapi.AttachNfsVolumeData{
			Version: casVersion,
		},
		&openapi.AttachNFSVolumeOpts{
			IgnoreVersion: ignoreVersion,
			AsyncMax:      asyncMax,
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

// UpdateNFSVolumeExports request to update the NFS volume exports of `volumeID`
// in the namespace `namespaceID`. It can return an error or nil if it succeeds.
func (o *OpenAPI) UpdateNFSVolumeExports(
	ctx context.Context,
	namespaceID id.Namespace,
	volumeID id.Volume,
	exports []volume.NFSExportConfig,
	params *apiclient.UpdateNFSVolumeExportsRequestParams,
) error {

	o.mu.RLock()
	defer o.mu.RUnlock()

	var casVersion string
	var ignoreVersion optional.Bool = optional.NewBool(true)
	var asyncMax optional.String = optional.EmptyString()

	if params != nil {
		if params.CASVersion.String() != "" {
			ignoreVersion = optional.NewBool(false)
			casVersion = params.CASVersion.String()
		}

		if params.AsyncMax != 0 {
			asyncMax = optional.NewString(params.AsyncMax.String())
		}
	}

	openapiExports := make([]openapi.NfsExportConfig, 0, len(exports))
	for _, e := range exports {
		openapiExports = append(openapiExports, o.codec.encodeNFSExport(e))
	}

	resp, err := o.client.DefaultApi.UpdateNFSVolumeExports(
		ctx,
		namespaceID.String(),
		volumeID.String(),
		openapi.NfsVolumeExports{
			Exports: openapiExports,
			Version: casVersion,
		},
		&openapi.UpdateNFSVolumeExportsOpts{
			IgnoreVersion: ignoreVersion,
			AsyncMax:      asyncMax,
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

// UpdateNFSVolumeMountEndpoint request to update the NFS mount endpoint of
// `volumeID` in the namespace `namespaceID`. It can return an error or nil if
// it succeeds.
func (o *OpenAPI) UpdateNFSVolumeMountEndpoint(
	ctx context.Context,
	namespaceID id.Namespace,
	volumeID id.Volume,
	endpoint string,
	params *apiclient.UpdateNFSVolumeMountEndpointRequestParams,
) error {

	o.mu.RLock()
	defer o.mu.RUnlock()

	var casVersion string
	var ignoreVersion optional.Bool = optional.NewBool(true)
	var asyncMax optional.String = optional.EmptyString()

	if params != nil {
		if params.CASVersion.String() != "" {
			ignoreVersion = optional.NewBool(false)
			casVersion = params.CASVersion.String()
		}

		if params.AsyncMax != 0 {
			asyncMax = optional.NewString(params.AsyncMax.String())
		}
	}

	resp, err := o.client.DefaultApi.UpdateNFSVolumeMountEndpoint(
		ctx,
		namespaceID.String(),
		volumeID.String(),
		openapi.NfsVolumeMountEndpoint{
			MountEndpoint: endpoint,
			Version:       casVersion,
		},
		&openapi.UpdateNFSVolumeMountEndpointOpts{
			IgnoreVersion: ignoreVersion,
			AsyncMax:      asyncMax,
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

// SetFailureModeIntent attempts to perform an update of the failure mode
// for the target volume to the provided intent-based behaviour.
func (o *OpenAPI) SetFailureModeIntent(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, intent string, params *apiclient.SetFailureModeRequestParams) (*volume.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	partialRequest := openapi.SetFailureModeRequest{
		Mode: openapi.FailureModeIntent(intent),
	}

	return o.setFailureMode(ctx, namespaceID, volumeID, partialRequest, params)
}

// SetFailureThreshold attempts to perform an update of the failure mode
// for the target volume to the provided numerical threshold.
func (o *OpenAPI) SetFailureThreshold(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, threshold uint64, params *apiclient.SetFailureModeRequestParams) (*volume.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	partialRequest := openapi.SetFailureModeRequest{
		FailureThreshold: threshold,
	}

	return o.setFailureMode(ctx, namespaceID, volumeID, partialRequest, params)
}

// setFailureMode takes the partial request provided to it, completing it as
// appropriate for the given params.
func (o *OpenAPI) setFailureMode(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, partialRequest openapi.SetFailureModeRequest, params *apiclient.SetFailureModeRequestParams) (*volume.Resource, error) {
	opts := &openapi.SetFailureModeOpts{
		IgnoreVersion: optional.NewBool(true),
	}

	if params != nil && params.CASVersion != "" {
		partialRequest.Version = params.CASVersion.String()
		opts.IgnoreVersion = optional.NewBool(false)
	}

	model, resp, err := o.client.DefaultApi.SetFailureMode(
		ctx,
		namespaceID.String(),
		volumeID.String(),
		partialRequest,
		opts,
	)
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

// DetachVolume makes a detach request for volumeID in namespaceID.
//
// The behaviour of the operation is dictated by params:
//
//
//  Version constraints:
// 	- If params is nil or params.CASVersion is empty then the detach request is
// 	unconditional
// 	- If params.CASVersion is set, the request is conditional upon it matching
// 	the volume entity's version as seen by the server.
//
//  Asynchrony:
//  - If params is nil or params.AsyncMax is empty/zero valued then the delete
//  request is performed synchronously.
//  - If params.AsyncMax is set, the request is performed asynchronously using
//  the duration given as the maximum amount of time allowed for the request
//  before it times out.
func (o *OpenAPI) DetachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *apiclient.DetachVolumeRequestParams) error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	var casVersion string

	opts := &openapi.DetachVolumeOpts{
		IgnoreVersion: optional.NewBool(true),
		AsyncMax:      optional.EmptyString(),
	}

	// Set the CAS version constraint if provided
	if params != nil {
		if params.CASVersion.String() != "" {
			opts.IgnoreVersion = optional.NewBool(false)
			casVersion = params.CASVersion.String()
		}
		if params.AsyncMax != 0 {
			opts.AsyncMax = optional.NewString(params.AsyncMax.String())
		}
	}

	resp, err := o.client.DefaultApi.DetachVolume(
		ctx,
		namespaceID.String(),
		volumeID.String(),
		casVersion,
		opts,
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

// SetReplicas changes the number of the replicas of a specified volume.
// Operation is asynchronous, we return nil if the request has been accepted.
func (o *OpenAPI) SetReplicas(ctx context.Context, nsID id.Namespace, volID id.Volume, numReplicas uint64, params *apiclient.SetReplicasRequestParams) error {

	o.mu.RLock()
	defer o.mu.RUnlock()

	// default
	request := openapi.SetReplicasRequest{Replicas: numReplicas}
	opts := &openapi.SetReplicasOpts{
		IgnoreVersion: optional.NewBool(true),
	}

	// check optional params
	if params != nil && params.CASVersion != "" {
		request.Version = params.CASVersion.String()
		opts.IgnoreVersion = optional.NewBool(false)
	}

	_, resp, err := o.client.DefaultApi.SetReplicas(ctx, nsID.String(), volID.String(), request, opts)
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

// UpdateVolume changes the description of a specified volume.
//
//  Version constraints:
// 	- If params is nil or params.CASVersion is empty then the detach request is
// 	unconditional
// 	- If params.CASVersion is set, the request is conditional upon it matching
// 	the volume entity's version as seen by the server.
//
//  Asynchrony:
//  - If params is nil or params.AsyncMax is empty/zero valued then the delete
//  request is performed synchronously.
//  - If params.AsyncMax is set, the request is performed asynchronously using
//  the duration given as the maximum amount of time allowed for the request
//  before it times out.
func (o *OpenAPI) UpdateVolume(
	ctx context.Context,
	nsID id.Namespace,
	volID id.Volume,
	description string,
	labels labels.Set,
	params *apiclient.UpdateVolumeRequestParams,
) (*volume.Resource, error) {

	o.mu.RLock()
	defer o.mu.RUnlock()

	// default
	request := openapi.UpdateVolumeData{
		Labels:      labels,
		Description: description,
	}
	opts := &openapi.UpdateVolumeOpts{
		IgnoreVersion: optional.NewBool(true),
		AsyncMax:      optional.EmptyString(),
	}

	// check optional params
	if params != nil && params.CASVersion != "" {
		request.Version = params.CASVersion.String()
		opts.IgnoreVersion = optional.NewBool(false)
	}

	model, resp, err := o.client.DefaultApi.UpdateVolume(ctx, nsID.String(), volID.String(), request, opts)
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

// ResizeVolume changes the size of a specified volume.
// Operation is asynchronous, we return nil if the request has been accepted.
func (o *OpenAPI) ResizeVolume(
	ctx context.Context,
	nsID id.Namespace,
	volID id.Volume,
	sizeBytes uint64,
	params *apiclient.ResizeVolumeRequestParams,
) (*volume.Resource, error) {

	o.mu.RLock()
	defer o.mu.RUnlock()

	// default
	request := openapi.ResizeVolumeRequest{
		SizeBytes: sizeBytes,
	}
	opts := &openapi.ResizeVolumeOpts{
		AsyncMax:      optional.EmptyString(),
		IgnoreVersion: optional.Bool{},
	}

	// check optional params
	if params != nil {
		if params.AsyncMax != 0 {
			opts.AsyncMax = optional.NewString(params.AsyncMax.String())
		}

		if params.CASVersion != "" {
			request.Version = params.CASVersion.String()
			opts.IgnoreVersion = optional.NewBool(false)
		}

		if params.AsyncMax != 0 {
			opts.AsyncMax = optional.NewString(params.AsyncMax.String())
		}
	}

	model, resp, err := o.client.DefaultApi.ResizeVolume(ctx, nsID.String(), volID.String(), request, opts)
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
