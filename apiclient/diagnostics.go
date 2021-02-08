package apiclient

import (
	"context"

	"code.storageos.net/storageos/c2-cli/diagnostics"
)

// GetSingleNodeDiagnosticsByName requests a local diagnostics bundle from the given
// node.
func (c *Client) GetSingleNodeDiagnosticsByName(ctx context.Context, name string) (*diagnostics.BundleReadCloser, error) {

	node, err := c.GetNodeByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return c.Transport.GetSingleNodeDiagnostics(ctx, node.ID)

}
