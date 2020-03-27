package output

import (
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/policygroup"
	"code.storageos.net/storageos/c2-cli/user"
)

// User encapsulates the information required to output a StorageOS user account
// to a display.
type User struct {
	ID       id.User `json:"id" yaml:"id"`
	Username string  `json:"name" yaml:"name"`

	IsAdmin   bool            `json:"isAdmin" yaml:"isAdmin"`
	Groups    []PolicyGroup   `json:"groups" yaml:"groups"`
	CreatedAt time.Time       `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt" yaml:"updatedAt"`
	Version   version.Version `json:"version" yaml:"version"`
}

// PolicyGroup encapsulates the information required to output policy groups
// to a display.
type PolicyGroup struct {
	ID   id.PolicyGroup `json:"id" yaml:"id"`
	Name string         `json:"name" yaml:"name"`
}

// NewUser creates a new User output representation using extra details from
// the provided parameters.
func NewUser(user *user.Resource, policyGroups map[id.PolicyGroup]*policygroup.Resource) (*User, error) {
	outputGroups := make([]PolicyGroup, 0, len(user.Groups))
	for _, gid := range user.Groups {
		group, ok := policyGroups[gid]
		if !ok {
			return nil, NewMissingRequiredPolicyGroupErr(gid)
		}

		outputGroups = append(outputGroups, PolicyGroup{
			ID:   gid,
			Name: group.Name,
		})
	}

	return &User{
		ID:        user.ID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		Groups:    outputGroups,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Version:   user.Version,
	}, nil
}
