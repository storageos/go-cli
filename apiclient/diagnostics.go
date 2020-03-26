package apiclient

import (
	"context"
	"io"
)

// GetDiagnostics requests a new diagnostics bundle for the cluster
// from the StorageOS API.
func (c *Client) GetDiagnostics(ctx context.Context) (io.ReadCloser, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.GetDiagnostics(ctx)
}
