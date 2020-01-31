package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/openapi"
)

// GetCluster requests the cluster configuration from the StorageOS API,
// translating it into a *cluster.Resource.
func (o *OpenAPI) GetCluster(ctx context.Context) (*cluster.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, resp, err := o.client.DefaultApi.GetCluster(ctx)
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	return o.codec.decodeCluster(model)
}

// UpdateCluster attempts to perform an update of the cluster configuration
// through the StorageOS API using resource and licenceKey as the update values.
func (o *OpenAPI) UpdateCluster(ctx context.Context, resource *cluster.Resource, licenceKey []byte) (*cluster.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	level, err := o.codec.encodeLogLevel(resource.LogLevel)
	if err != nil {
		return nil, err
	}

	format, err := o.codec.encodeLogFormat(resource.LogFormat)
	if err != nil {
		return nil, err
	}

	updateData := openapi.UpdateClusterData{
		LicenceKey:            string(licenceKey),
		DisableTelemetry:      resource.DisableTelemetry,
		DisableCrashReporting: resource.DisableCrashReporting,
		DisableVersionCheck:   resource.DisableVersionCheck,
		LogLevel:              level,
		LogFormat:             format,
		Version:               resource.Version.String(),
	}

	model, resp, err := o.client.DefaultApi.UpdateCluster(ctx, updateData)
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	return o.codec.decodeCluster(model)
}
