package apiclient

import (
	"context"

	"code.storageos.net/storageos/c2-cli/cluster"
)

// GetCluster requests basic information for the cluster resource from the
// StorageOS API.
func (c *Client) GetCluster(ctx context.Context) (*cluster.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.GetCluster(ctx)
}

// UpdateLicence fetches the current cluster configuration, then applies an
// licenceKey to it. If successful the newly applied licence configuration is
// returned.
func (c *Client) UpdateLicence(ctx context.Context, licenceKey []byte) (*cluster.Licence, error) {
	config, err := c.GetCluster(ctx)
	if err != nil {
		return nil, err
	}

	updated, err := c.transport.UpdateCluster(ctx, config, licenceKey)
	if err != nil {
		return nil, err
	}

	return updated.Licence, nil
}
