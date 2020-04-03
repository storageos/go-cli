package apiclient

import (
	"context"

	"code.storageos.net/storageos/c2-cli/cluster"
)

// UpdateLicence fetches the current cluster configuration, then applies an
// licenceKey to it. If successful the newly applied licence configuration is
// returned.
func (c *Client) UpdateLicence(ctx context.Context, licenceKey []byte) (*cluster.Licence, error) {
	config, err := c.GetCluster(ctx)
	if err != nil {
		return nil, err
	}

	updated, err := c.Transport.UpdateCluster(ctx, config, licenceKey)
	if err != nil {
		return nil, err
	}

	return updated.Licence, nil
}
