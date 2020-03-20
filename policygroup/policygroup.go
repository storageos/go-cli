package policygroup

import (
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

// Resource encapsulates a StorageOS policy group API resource as a data type.
type Resource struct {
	ID    id.PolicyGroup `json:"id"`
	Name  string         `json:"name"`
	Users []*Member      `json:"users"`
	Specs []*Spec        `json:"specs"`

	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Version   version.Version `json:"version"`
}

// Spec encapsulates a policy specification API resource belonging to a policy
// group as a data type.
type Spec struct {
	NamespaceID  id.Namespace `json:"namespaceID"`
	ResourceType string       `json:"resourceType"`
	ReadOnly     bool         `json:"readOnly"`
}

// Member represents the details of a user that is a member of a policy group.
type Member struct {
	ID       id.User `json:"id"`
	Username string  `json:"username"`
}
