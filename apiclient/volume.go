package apiclient

import (
	"context"
	"fmt"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
)

// GetVolume requests basic information for the volume resource which
// corresponds to uid in namespace from the StorageOS API.
func (c *Client) GetVolume(ctx context.Context, namespace id.Namespace, uid id.Volume) (*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.GetVolume(ctx, namespace, uid)
}

// GetVolumeByName requests basic information for the volume resource which has
// name in namespace.
//
// The resource model for the API is build around using unique identifiers,
// so this operation is inherently more expensive than the corresponding
// GetVolume() operation.
//
// Retrieving a volume resource by name involves requesting a list of all
// volumes in the namespace from the StorageOS API and returning the first one
// where the name matches.
func (c *Client) GetVolumeByName(ctx context.Context, namespace id.Namespace, name string) (*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	volumes, err := c.transport.ListVolumes(ctx, namespace)
	if err != nil {
		return nil, err
	}

	for _, v := range volumes {
		if v.Name == name {
			return v, nil
		}
	}

	return nil, NewNotFoundError(fmt.Sprintf("volume with name %v not found", name))
}

// GetNamespaceVolumes requests basic information for each volume resource in
// namespace from the StorageOS API.
//
// The returned list is filtered using uids so that it contains only those
// resources which have a matching ID. Omitting uids will skip the filtering.
func (c *Client) GetNamespaceVolumes(ctx context.Context, namespace id.Namespace, uids ...id.Volume) ([]*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	volumes, err := c.transport.ListVolumes(ctx, namespace)
	if err != nil {
		return nil, err
	}

	if len(uids) == 0 {
		return volumes, nil
	}

	// Filter uids have been provided:
	retrieved := map[id.Volume]*volume.Resource{}

	for _, v := range volumes {
		retrieved[v.ID] = v
	}

	filtered := make([]*volume.Resource, len(uids))

	i := 0
	for _, id := range uids {
		v, ok := retrieved[id]
		if ok {
			filtered[i] = v
			i++
		} else {
			return nil, NewNotFoundError(fmt.Sprintf("volume %v not found", id))
		}
	}

	return filtered, nil
}

// GetAllVolumes requests basic information for each volume resource in every
// namespace exposed by the StorageOS API to the authenticated user.
func (c *Client) GetAllVolumes(ctx context.Context) ([]*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.fetchAllVolumes(ctx)
}
