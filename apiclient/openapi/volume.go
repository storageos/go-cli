package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
	"code.storageos.net/storageos/openapi"
)

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
		return nil, mapOpenAPIError(err, resp)
	}

	return o.codec.decodeVolume(model)
}

func (o *OpenAPI) GetVolume(ctx context.Context, namespace id.Namespace, uid id.Volume) (*volume.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, resp, err := o.client.DefaultApi.GetVolume(ctx, namespace.String(), uid.String())
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	return o.codec.decodeVolume(model)
}

func (o *OpenAPI) ListVolumes(ctx context.Context, namespace id.Namespace) ([]*volume.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	models, resp, err := o.client.DefaultApi.ListVolumes(ctx, namespace.String())
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
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
