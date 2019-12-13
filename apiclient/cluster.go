package apiclient

import (
	"context"

	"code.storageos.net/storageos/c2-cli/cluster"
)

// GetCluster requests basic information for the cluster resource from the
// StorageOS API.
func (c *Client) GetCluster(ctx context.Context) (*cluster.Resource, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.GetCluster(ctx)
}
