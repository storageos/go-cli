package node

import (
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/capacity"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

// Resource encapsulates a StorageOS node API resource as a data type.
type Resource struct {
	ID       id.Node          `json:"id"`
	Name     string           `json:"name"`
	Health   health.NodeState `json:"health"`
	Capacity capacity.Stats   `json:"capacity,omitempty"`

	IOAddr         string `json:"ioAddress"`
	SupervisorAddr string `json:"supervisorAddress"`
	GossipAddr     string `json:"gossipAddress"`
	ClusteringAddr string `json:"clusteringAddress"`

	Labels labels.Set `json:"labels"`

	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Version   version.Version `json:"version"`
}
