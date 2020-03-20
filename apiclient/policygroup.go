package apiclient

import (
	"context"
	"fmt"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/policygroup"
)

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

// GetListPolicyGroups requests a list containing basic information on each
// policy group configured for the cluster.
//
// The returned list is filtered using gids so that it contains  only those
// resources which have a matching GID. Omitting gids  will skip the filtering.
func (c *Client) GetListPolicyGroups(ctx context.Context, gids ...id.PolicyGroup) ([]*policygroup.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	policyGroups, err := c.transport.ListPolicyGroups(ctx)
	if err != nil {
		return nil, err
	}

	return filterPolicyGroupsForIDs(policyGroups, gids...)
}

// GetListPolicyGroupsByName requests a list containing basic information on each
// policy group configured for the cluster.
//
// The returned list is filtered using name so that it contains  only those
// resources which have a matching name. Omitting gids  will skip the filtering.
func (c *Client) GetListPolicyGroupsByName(ctx context.Context, names ...string) ([]*policygroup.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	policyGroups, err := c.transport.ListPolicyGroups(ctx)
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
