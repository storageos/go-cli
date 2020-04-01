package apiclient

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/volume"
)

// VolumeExistsError is returned when a volume creation request is sent to the
// StorageOS API for a namespace where name is already in use.
type VolumeExistsError struct {
	name        string
	namespaceID id.Namespace
}

// Error returns an error message indicating that a volume name is already in
// use for the target namespace.
func (e VolumeExistsError) Error() string {
	return fmt.Sprintf("volume name %v is already in use for namespace with ID %v", e.name, e.namespaceID)
}

// NewVolumeExistsError returns an error indicating that a volume with name
// already exists in namespaceID.
func NewVolumeExistsError(name string, namespaceID id.Namespace) VolumeExistsError {
	return VolumeExistsError{
		name:        name,
		namespaceID: namespaceID,
	}
}

// InvalidVolumeCreationError is returned when a volume creation request sent
// to the StorageOS API is invalid.
type InvalidVolumeCreationError struct {
	details string
}

// Error returns an error message indicating that a volume creation request
// made to the StorageOS API is invalid, including details if available.
func (e InvalidVolumeCreationError) Error() string {
	msg := "volume creation request is invalid"
	if e.details != "" {
		msg = fmt.Sprintf("%v: %v", msg, e.details)
	}
	return msg
}

// NewInvalidVolumeCreationError returns an InvalidVolumeCreationError, using
// details to provide information about what must be corrected.
func NewInvalidVolumeCreationError(details string) InvalidVolumeCreationError {
	return InvalidVolumeCreationError{
		details: details,
	}
}

// VolumeNotFoundError indicates that the API could not find the StorageOS volume
// specified.
type VolumeNotFoundError struct {
	msg string

	uid  id.Volume
	name string
}

// Error returns an error message indicating that the volume with a given
// ID or name was not found, as configured.
func (e VolumeNotFoundError) Error() string {
	return e.msg
}

// NewVolumeNotFoundError returns a VolumeNotFoundError using details as the
// the error message. This can be used when provided an opaque but detailed
// error strings.
func NewVolumeNotFoundError(details string) VolumeNotFoundError {
	return VolumeNotFoundError{
		msg: details,
	}
}

// NewVolumeIDNotFoundError returns a VolumeNotFoundError for the volume with uid,
// constructing a user friendly message and storing the ID inside the error.
func NewVolumeIDNotFoundError(volumeID id.Volume) VolumeNotFoundError {
	return VolumeNotFoundError{
		msg: fmt.Sprintf("volume with ID %v not found for target namespace", volumeID),
		uid: volumeID,
	}
}

// NewVolumeNameNotFoundError returns a VolumeNotFoundError for the volume
// with name, constructing a user friendly message and storing the name inside
// the error.
func NewVolumeNameNotFoundError(name string) VolumeNotFoundError {
	return VolumeNotFoundError{
		msg:  fmt.Sprintf("volume with name %v not found for target namespace", name),
		name: name,
	}
}

// CreateVolumeRequestParams contains optional request parameters for a create
// volume operation.
type CreateVolumeRequestParams struct {
	AsyncMax time.Duration
}

// CreateVolume requests the creation of a new StorageOS volume in namespace
// from the provided fields. If successful the created resource for the volume
// is returned to the caller.
//
// The behaviour of the operation is dictated by params:
//
//
//  Asynchrony:
//  - If params is nil or params.AsyncMax is empty/zero valued then the create
//  request is performed synchronously.
//  - If params.AsyncMax is set, the request is performed asynchronously using
//  the duration given as the maximum amount of time allowed for the request
//  before it times out.
func (c *Client) CreateVolume(ctx context.Context, namespace id.Namespace, name, description string, fs volume.FsType, sizeBytes uint64, labelSet labels.Set, params *CreateVolumeRequestParams) (*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.CreateVolume(ctx, namespace, name, description, fs, sizeBytes, labelSet, params)
}

// GetVolume requests basic information for the volume resource which
// corresponds to uid in namespace from the StorageOS API.
func (c *Client) GetVolume(ctx context.Context, namespace id.Namespace, volumeID id.Volume) (*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.GetVolume(ctx, namespace, volumeID)
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
func (c *Client) GetVolumeByName(ctx context.Context, namespaceID id.Namespace, name string) (*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	volumes, err := c.transport.ListVolumes(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	for _, v := range volumes {
		if v.Name == name {
			return v, nil
		}
	}

	return nil, NewVolumeNameNotFoundError(name)
}

// GetNamespaceVolumes requests basic information for each volume resource in
// namespace from the StorageOS API.
//
// The returned list is filtered using uids so that it contains only those
// resources which have a matching ID. Omitting uids will skip the filtering.
func (c *Client) GetNamespaceVolumes(ctx context.Context, namespaceID id.Namespace, uids ...id.Volume) ([]*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	volumes, err := c.transport.ListVolumes(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return filterVolumesForUIDs(volumes, uids...)
}

// GetNamespaceVolumesByName requests basic information for each volume resource in
// namespace from the StorageOS API.
//
// The returned list is filtered using names so that it contains only those
// resources which have a matching name. Omitting names will skip the filtering.
func (c *Client) GetNamespaceVolumesByName(ctx context.Context, namespaceID id.Namespace, names ...string) ([]*volume.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	volumes, err := c.transport.ListVolumes(ctx, namespaceID)
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

	return c.fetchAllVolumesParallel(ctx)
}

// fetchAllVolumesParallel requests the list of all namespaces from the
// StorageOS API, then requests the list of volumes within each namespace,
// calling all of them in parallel, returning an aggregate list of the volumes
// returned.
//
// If access is not granted when listing volumes for a retrieved namespace it
// is noted but will not return an error. Only if access is denied for all
// attempts will this return a permissions error.
//
// If any of the call returns an error:
//  - the context is canceled so all pending requests are cut
//  - this method returns an error
func (c *Client) fetchAllVolumesParallel(ctx context.Context) ([]*volume.Resource, error) {
	namespaces, err := c.transport.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	// The derived Context is canceled the first time a function passed to Go
	// returns a non-nil error or the first time Wait returns, whichever occurs
	// first.
	group, ctx := errgroup.WithContext(ctx)

	results := make(chan []*volume.Resource, len(namespaces))

	for _, ns := range namespaces {
		ns := ns

		// Go calls the given function in a new goroutine.
		//
		// The first call to return a non-nil error cancels the group; its error
		// will be returned by Wait.
		group.Go(func() error {

			nsvols, err := c.transport.ListVolumes(ctx, ns.ID)
			switch {
			case err == nil, errors.As(err, &UnauthorisedError{}):
				// For an unauthorised error, ignore - its not fatal to the operation.
			default:
				return err
			}

			results <- nsvols
			return nil
		})
	}

	// blocks until all function calls from the Go method have returned
	if err := group.Wait(); err != nil {
		return nil, err
	}

	close(results)

	// merge the results
	volumes := []*volume.Resource{}
	for r := range results {
		volumes = append(volumes, r...)
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

	filtered := make([]*volume.Resource, 0, len(uids))

	for _, id := range uids {
		v, ok := retrieved[id]
		if !ok {
			return nil, NewVolumeIDNotFoundError(id)
		}
		filtered = append(filtered, v)
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

	filtered := make([]*volume.Resource, 0, len(names))

	for _, name := range names {
		v, ok := retrieved[name]
		if !ok {
			return nil, NewVolumeNameNotFoundError(name)
		}
		filtered = append(filtered, v)
	}

	return filtered, nil
}

// DeleteVolumeRequestParams contains optional request parameters for a delete
// volume operation.
type DeleteVolumeRequestParams struct {
	CASVersion version.Version
	AsyncMax   time.Duration
}

// DeleteVolume makes a delete request for volumeID in namespaceID.
//
// The behaviour of the operation is dictated by params:
//
//
// 	Version constraints:
// 	- If params is nil or params.CASVersion is empty then the delete request is
// 	unconditional
// 	- If params.CASVersion is set, the request is conditional upon it matching
// 	the volume entity's version as seen by the server.
//
//  Asynchrony:
//  - If params is nil or params.AsyncMax is empty/zero valued then the delete
//  request is performed synchronously.
//  - If params.AsyncMax is set, the request is performed asynchronously using
//  the duration given as the maximum amount of time allowed for the request
//  before it times out.
func (c *Client) DeleteVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *DeleteVolumeRequestParams) error {
	_, err := c.authenticate(ctx)
	if err != nil {
		return err
	}

	return c.transport.DeleteVolume(ctx, namespaceID, volumeID, params)
}

// AttachVolume requests to attach a volume (namespace/volume) to a node
// It requires authentication.
func (c *Client) AttachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, nodeID id.Node) error {
	_, err := c.authenticate(ctx)
	if err != nil {
		return err
	}

	return c.transport.AttachVolume(ctx, namespaceID, volumeID, nodeID)
}

// DetachVolumeRequestParams contains optional request parameters for a detach
// volume operation.
type DetachVolumeRequestParams struct {
	CASVersion version.Version
}

// DetachVolume makes a detach request for volumeID in namespaceID.
//
// The behaviour of the operation is dictated by params:
//
//
//  Version constraints:
// 	- If params is nil or params.CASVersion is empty then the detach request is
// 	unconditional
// 	- If params.CASVersion is set, the request is conditional upon it matching
// 	the volume entity's version as seen by the server.
func (c *Client) DetachVolume(ctx context.Context, namespaceID id.Namespace, volumeID id.Volume, params *DetachVolumeRequestParams) error {
	_, err := c.authenticate(ctx)
	if err != nil {
		return err
	}

	return c.transport.DetachVolume(ctx, namespaceID, volumeID, params)
}
