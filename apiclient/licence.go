package apiclient

import (
	"context"

	"code.storageos.net/storageos/c2-cli/licence"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

// UpdateLicenceRequestParams contains optional request parameters for a update
// licence operation.
type UpdateLicenceRequestParams struct {
	CASVersion version.Version
}

// UpdateLicence sends a new version of the licence to apply to the current
// cluster. It returns the new licence resource if correctly applied.
// It doesn't require a version but overwrite the licence using the last
// available version from the cluster.
func (c *Client) UpdateLicence(ctx context.Context, licence []byte, params *UpdateLicenceRequestParams) (*licence.Resource, error) {

	if params == nil {
		cl, err := c.Transport.GetCluster(ctx)
		if err != nil {
			return nil, err
		}
		return c.Transport.UpdateLicence(ctx, licence, cl.Version)
	}

	return c.Transport.UpdateLicence(ctx, licence, params.CASVersion)
}
