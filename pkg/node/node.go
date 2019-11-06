package node

import (
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/entity"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
)

// Resource encapsulates a StorageOS node as a data type.
type Resource struct {
	ID     id.Node       `json:"id"`
	Name   string        `json:"name"`
	Health entity.Health `json:"health"`

	Configuration *Configuration `json:"configuration,omitempty"`

	Labels labels.Set `json:"labels"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Version   entity.Version `json:"version"`
}

// Configuration encapsulates the detailed configuration settings of a
// StorageOS node.
type Configuration struct {
	IOAddr         string `json:"ioAddress"`
	SupervisorAddr string `json:"supervisorAddress"`
	GossipAddr     string `json:"gossipAddress"`
	ClusteringAddr string `json:"clusteringAddress"`
}
