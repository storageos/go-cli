package user

import (
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

// Resource encapsulates a StorageOS user API resource as data type.
type Resource struct {
	ID       id.User `json:"id"`
	Username string  `json:"name"`

	IsAdmin bool             `json:"isAdmin"`
	Groups  []id.PolicyGroup `json:"groups"`

	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Version   version.Version `json:"version"`
}
