package apiclient

import (
	"context"
	"errors"
	"fmt"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/volume"
)

// CreateVolume requests the creation of a new StorageOS volume in namespace
// from the provided fields. If successful the created resource for the volume
// is returned to the caller.
func (c *Client) CreateVolume(ctx context.Context, namespace id.Namespace, name, description string, fs volume.FsType, sizeBytes uint64, labelSet labels.Set) (*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.CreateVolume(ctx, namespace, name, description, fs, sizeBytes, labelSet)
}

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

	return filterVolumesForUIDs(volumes, uids...)
}

// GetNamespaceVolumesByName requests basic information for each volume resource in
// namespace from the StorageOS API.
//
// The returned list is filtered using uids so that it contains only those
// resources which have a matching ID. Omitting uids will skip the filtering.
func (c *Client) GetNamespaceVolumesByName(ctx context.Context, namespace id.Namespace, names ...string) ([]*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	volumes, err := c.transport.ListVolumes(ctx, namespace)
	if err != nil {
		return nil, err
	}

	return filterVolumesForNames(volumes, names...)
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

// fetchAllVolumes requests the list of all namespaces from the StorageOS API,
// then requests the list of volumes within each namespace, returning an
// aggregate list of the volumes returned.
//
// If access is not granted when listing volumes for a retrieved namespace it
// is noted but will not return an error. Only if access is denied for all
// attempts will this return a permissions error.
func (c *Client) fetchAllVolumes(ctx context.Context) ([]*volume.Resource, error) {
	namespaces, err := c.transport.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	volumes := []*volume.Resource{}

	for _, ns := range namespaces {
		nsvols, err := c.transport.ListVolumes(ctx, ns.ID)
		switch {
		case err == nil, errors.As(err, &UnauthorisedError{}):
			// For an unauthorised error, ignore - its not fatal to the operation.
		default:
			return nil, err
		}
		volumes = append(volumes, nsvols...)
	}

	return volumes, nil
}

// filterVolumesForUIDs will return a subset of volumes containing resources
// which have one of the provided uids. If uids is not provided, volumes is
// returned as is.
//
// If there is no resource for a given uid then an error is returned, thus
// this is a strict helper.
func filterVolumesForUIDs(volumes []*volume.Resource, uids ...id.Volume) ([]*volume.Resource, error) {
	if len(uids) == 0 {
		return volumes, nil
	}

	retrieved := map[id.Volume]*volume.Resource{}

	for _, v := range volumes {
		retrieved[v.ID] = v
	}

	filtered := make([]*volume.Resource, len(uids))

	i := 0
	for _, id := range uids {
		v, ok := retrieved[id]
		if !ok {
			return nil, NewNotFoundError(fmt.Sprintf("volume %v not found", id))
		}
		filtered[i] = v
		i++
	}

	return filtered, nil
}

// filterVolumesForNames will return a subset of volumes containing resources
// which have one of the provided names. If names is not provided, volumes is
// returned as is.
//
// If there is no resource for a given name then an error is returned, thus
// this is a strict helper.
func filterVolumesForNames(volumes []*volume.Resource, names ...string) ([]*volume.Resource, error) {
	if len(names) == 0 {
		return volumes, nil
	}

	retrieved := map[string]*volume.Resource{}

	for _, v := range volumes {
		retrieved[v.Name] = v
	}

	filtered := make([]*volume.Resource, len(names))

	i := 0
	for _, name := range names {
		v, ok := retrieved[name]
		if !ok {
			return nil, NewNotFoundError(fmt.Sprintf("volume with name %v not found", name))
		}
		filtered[i] = v
		i++
	}

	return filtered, nil
}
