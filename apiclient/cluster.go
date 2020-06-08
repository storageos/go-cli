package apiclient

import "code.storageos.net/storageos/c2-cli/pkg/version"

// UpdateClusterRequestParams contains optional request parameters for a update
// cluster operation.
type UpdateClusterRequestParams struct {
	CASVersion version.Version
}
