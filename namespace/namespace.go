package namespace

import (
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/entity"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
)

// Resource encapsulates a StorageOS namespace API resource as a data type.
type Resource struct {
	ID     id.Namespace `json:"id"`
	Name   string       `json:"name"`
	Labels labels.Set   `json:"labels"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Version   entity.Version `json:"version"`
}
