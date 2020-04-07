package output

import (
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/policygroup"
)

// PolicyGroup defines a type that contains all the info we need to output a
// namespace.
type PolicyGroup struct {
	ID        id.PolicyGroup  `json:"id" yaml:"id"`
	Name      string          `json:"name" yaml:"name"`
	Users     []*Member       `json:"users"`
	Specs     []*Spec         `json:"specs"`
	CreatedAt time.Time       `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt" yaml:"updatedAt"`
	Version   version.Version `json:"version" yaml:"version"`
}

// Spec encapsulates a policy specification API resource belonging to a policy
// group as a data type.
type Spec struct {
	NamespaceID  id.Namespace `json:"namespaceID" yaml:"namespaceID"`
	ResourceType string       `json:"resourceType" yaml:"resourceType"`
	ReadOnly     bool         `json:"readOnly" yaml:"readOnly"`
}

// Member represents the details of a user that is a member of a policy group.
type Member struct {
	ID       id.User `json:"id" yaml:"id"`
	Username string  `json:"username" yaml:"username"`
}

// NewPolicyGroup returns a new PolicyGroup object that contains all the info needed
// to be outputted.
func NewPolicyGroup(n *policygroup.Resource) *PolicyGroup {
	users := make([]*Member, 0, len(n.Users))
	for _, u := range n.Users {
		users = append(
			users,
			&Member{
				ID:       u.ID,
				Username: u.Username,
			},
		)
	}

	specs := make([]*Spec, 0, len(n.Specs))
	for _, s := range n.Specs {
		specs = append(
			specs,
			&Spec{
				NamespaceID:  s.NamespaceID,
				ResourceType: s.ResourceType,
				ReadOnly:     s.ReadOnly,
			},
		)
	}

	return &PolicyGroup{
		ID:        n.ID,
		Name:      n.Name,
		Users:     users,
		Specs:     specs,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Version:   n.Version,
	}
}

// NewPolicyGroups returns a list of PolicyGroup objects that contains all the info
// needed to be outputted.
func NewPolicyGroups(pg []*policygroup.Resource) []*PolicyGroup {
	policyGroups := make([]*PolicyGroup, 0, len(pg))
	for _, n := range pg {
		policyGroups = append(policyGroups, NewPolicyGroup(n))
	}
	return policyGroups
}

// PolicyGroupDeletion defines a policy group deletion confirmation output
// representation.
type PolicyGroupDeletion struct {
	ID id.PolicyGroup `json:"id" yaml:"id"`
}
