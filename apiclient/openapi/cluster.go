package openapi

import (
	"context"

	"github.com/antihax/optional"

	"code.storageos.net/storageos/c2-cli/apiclient"
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
// through the StorageOS API using resource as the update value.
func (o *OpenAPI) UpdateCluster(ctx context.Context, resource *cluster.Resource, params *apiclient.UpdateClusterRequestParams) (*cluster.Resource, error) {
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
		DisableTelemetry:      resource.DisableTelemetry,
		DisableCrashReporting: resource.DisableCrashReporting,
		DisableVersionCheck:   resource.DisableVersionCheck,
		LogLevel:              level,
		LogFormat:             format,
	}

	opts := &openapi.UpdateClusterOpts{
		IgnoreVersion: optional.NewBool(true),
	}

	// check optional params
	if params != nil && params.CASVersion != "" {
		updateData.Version = params.CASVersion.String()
		opts.IgnoreVersion = optional.NewBool(false)
	}

	model, resp, err := o.client.DefaultApi.UpdateCluster(ctx, updateData, opts)
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	return o.codec.decodeCluster(model)
}
