package volume

import (
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

// FsType indicates the kind of filesystem which a volume has been given.
type FsType string

func FsTypeFromString(name string) FsType {
	return FsType(name)
}

// Resource encapsulates a StorageOS volume API resource as a data type.
type Resource struct {
	ID          id.Volume `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AttachedOn  id.Node   `json:"attachedOn"`

	Namespace  id.Namespace `json:"namespaceID"`
	Labels     labels.Set   `json:"labels"`
	Filesystem FsType       `json:"filesystem"`
	Inode      uint32       `json:"inode"`
	SizeBytes  uint64       `json:"sizeBytes"`

	Master   *Deployment   `json:"master"`
	Replicas []*Deployment `json:"replicas"`

	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Version   version.Version `json:"version"`
}

// Deployment encapsulates a deployment instance for a
// volume as a data type.
type Deployment struct {
	ID      id.Deployment `json:"id"`
	Node    id.Node       `json:"nodeID"`
	Inode   uint32        `json:"inode"`
	Health  health.State  `json:"health"`
	Syncing bool          `json:"syncing"`
}
