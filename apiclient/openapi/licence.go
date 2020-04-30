package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/licence"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/openapi"
)

// GetLicence requests the current cluster licence from the StorageOS API,
// translating it into a *licence.Resource.
func (o *OpenAPI) GetLicence(ctx context.Context) (*licence.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, resp, err := o.client.DefaultApi.GetLicence(ctx)
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	return o.codec.decodeLicence(model)
}

// UpdateLicence sends a new version of the licence to apply to the current
// cluster. It returns the new licence resource if correctly applied.
func (o *OpenAPI) UpdateLicence(ctx context.Context, licence []byte, casVersion version.Version) (*licence.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	updateData := openapi.UpdateLicence{
		LicenceKey: string(licence),
		Version:    casVersion.String(),
	}

	lic, resp, err := o.client.DefaultApi.UpdateLicence(ctx, updateData)
	if err != nil {
		return nil, mapOpenAPIError(err, resp)
	}

	return o.codec.decodeLicence(lic)
}
