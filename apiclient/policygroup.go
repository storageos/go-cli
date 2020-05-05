package apiclient

import (
	"context"
	"fmt"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/policygroup"
)

// DeletePolicyGroupRequestParams contains optional request parameters for a
// delete policy group operation.
type DeletePolicyGroupRequestParams struct {
	CASVersion version.Version
}

// PolicyGroupNotFoundError indicates that the API could not find the policy
// group specified.
type PolicyGroupNotFoundError struct {
	msg string

	gid  id.PolicyGroup
	name string
}

// Error returns an error message indicating that the policy group with a given
// ID or name was not found, as configured.
func (e PolicyGroupNotFoundError) Error() string {
	switch {
	case e.gid != "":
		return fmt.Sprintf("policy group with ID %v not found", e.gid)
	case e.name != "":
		return fmt.Sprintf("policy group with name %v not found", e.name)
	}

	return e.msg
}

// NewPolicyGroupIDNotFoundError returns a PolicyGroupNotFoundError for the
// policy group with gid, constructing a user friendly message and storing
// the ID inside the error.
func NewPolicyGroupIDNotFoundError(gid id.PolicyGroup) PolicyGroupNotFoundError {
	return PolicyGroupNotFoundError{
		gid: gid,
	}
}

// NewPolicyGroupNameNotFoundError returns a PolicyGroupNotFoundError for the
// policy group with name, constructing a user friendly message and storing
// the name inside the error.
func NewPolicyGroupNameNotFoundError(name string) PolicyGroupNotFoundError {
	return PolicyGroupNotFoundError{
		name: name,
	}
}

// PolicyGroupExistsError is returned when a policy group creation request is sent to
// a cluster where that name is already in use.
type PolicyGroupExistsError struct {
	name string
}

// Error returns an error message indicating that a policy group name is already in
// use.
func (e PolicyGroupExistsError) Error() string {
	return fmt.Sprintf("policy group name %s is already in use", e.name)
}

// NewPolicyGroupExistsError returns an error indicating that a policy group with
// that name already exists.
func NewPolicyGroupExistsError(name string) PolicyGroupExistsError {
	return PolicyGroupExistsError{
		name: name,
	}
}

// InvalidPolicyGroupCreationError is returned when a policy group creation
// request sent to the StorageOS API is invalid.
type InvalidPolicyGroupCreationError struct {
	details string
}

// Error returns an error message indicating that a policy group creation
// request made to the StorageOS API is invalid, including details if available.
func (e InvalidPolicyGroupCreationError) Error() string {
	msg := "policy group creation request is invalid"
	if e.details != "" {
		msg = fmt.Sprintf("%v: %v", msg, e.details)
	}
	return msg
}

// NewInvalidPolicyGroupCreationError returns an InvalidPolicyGroupCreationError,
// using details to provide information about what must be corrected.
func NewInvalidPolicyGroupCreationError(details string) InvalidPolicyGroupCreationError {
	return InvalidPolicyGroupCreationError{
		details: details,
	}
}

// GetListPolicyGroupsByUID requests a list containing basic information on each
// policy group configured for the cluster.
//
// The returned list is filtered using gids so that it contains  only those
// resources which have a matching GID. Omitting gids  will skip the filtering.
func (c *Client) GetListPolicyGroupsByUID(ctx context.Context, gids ...id.PolicyGroup) ([]*policygroup.Resource, error) {
	policyGroups, err := c.Transport.ListPolicyGroups(ctx)
	if err != nil {
		return nil, err
	}

	return filterPolicyGroupsForIDs(policyGroups, gids...)
}

// GetPolicyGroupByName requests a policy group given its name
func (c *Client) GetPolicyGroupByName(ctx context.Context, name string) (*policygroup.Resource, error) {
	policyGroups, err := c.Transport.ListPolicyGroups(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range policyGroups {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, NewPolicyGroupNameNotFoundError(name)
}

// GetListPolicyGroupsByName requests a list containing basic information on each
// policy group configured for the cluster.
//
// The returned list is filtered using name so that it contains  only those
// resources which have a matching name. Omitting gids  will skip the filtering.
func (c *Client) GetListPolicyGroupsByName(ctx context.Context, names ...string) ([]*policygroup.Resource, error) {
	policyGroups, err := c.Transport.ListPolicyGroups(ctx)
	if err != nil {
		return nil, err
	}

	return filterPolicyGroupsForNames(policyGroups, names...)
}

// filterPolicyGroupsForNames will return a subset of policyGroups containing
// resources which have one of the provided names. If names is not provided,
// policyGroups is returned as is.
//
// If there is no resource for a given name then an error is returned, thus
// this is a strict helper.
func filterPolicyGroupsForNames(policyGroups []*policygroup.Resource, names ...string) ([]*policygroup.Resource, error) {
	// return everything if no filter names given
	if len(names) == 0 {
		return policyGroups, nil
	}

	// implicitly removes also duplicates, if any, in the names list
	retrieved := map[string]*policygroup.Resource{}
	for _, g := range policyGroups {
		retrieved[g.Name] = g
	}

	filtered := make([]*policygroup.Resource, 0, len(names))

	for _, name := range names {
		g, ok := retrieved[name]
		if !ok {
			return nil, NewPolicyGroupNameNotFoundError(name)
		}
		filtered = append(filtered, g)
	}

	return filtered, nil
}

// filterPolicyGroupsForIDs will return a subset of policyGroups containing
// resources which have one of the provided gids. If gids is not provided,
// policyGroups is returned as is.
//
// If there is no resource for a given gid then an error is returned, thus
// this is a strict helper.
func filterPolicyGroupsForIDs(policyGroups []*policygroup.Resource, gids ...id.PolicyGroup) ([]*policygroup.Resource, error) {
	if len(gids) == 0 {
		return policyGroups, nil
	}

	// implicitly removes also duplicates, if any, in the gids list
	retrieved := map[id.PolicyGroup]*policygroup.Resource{}
	for _, g := range policyGroups {
		retrieved[g.ID] = g
	}

	filtered := make([]*policygroup.Resource, 0, len(gids))

	for _, gid := range gids {
		g, ok := retrieved[gid]
		if !ok {
			return nil, NewPolicyGroupIDNotFoundError(gid)
		}
		filtered = append(filtered, g)
	}

	return filtered, nil
}
