package node

import (
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/entity"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/volume"
)

// State aggregates information that can be used to provide a detailed picture
// of a node's state.
type State struct {
	Resource    *Resource            `json:"resource"`
	Deployments []*volume.Deployment `json:"deployments"`
}

// Resource encapsulates a StorageOS node API resource as a data type.
type Resource struct {
	ID     id.Node       `json:"id"`
	Name   string        `json:"name"`
	Health entity.Health `json:"health"`

	IOAddr         string `json:"ioAddress"`
	SupervisorAddr string `json:"supervisorAddress"`
	GossipAddr     string `json:"gossipAddress"`
	ClusteringAddr string `json:"clusteringAddress"`

	Labels labels.Set `json:"labels"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Version   entity.Version `json:"version"`
}
