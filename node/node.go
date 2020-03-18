package node

import (
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/volume"
)

// State aggregates information that can be used to provide a detailed picture
// of a node's state.
type State struct {
	Resource    *Resource     `json:"node"`
	Deployments []*Deployment `json:"deployments"`
}

// Deployment augments a volume.Deployment with the ID of the volume it
// belongs to.
type Deployment struct {
	VolumeID   id.Volume          `json:"volumeID"`
	Deployment *volume.Deployment `json:"deployment"`
}

// Resource encapsulates a StorageOS node API resource as a data type.
type Resource struct {
	ID     id.Node          `json:"id"`
	Name   string           `json:"name"`
	Health health.NodeState `json:"health"`

	IOAddr         string `json:"ioAddress"`
	SupervisorAddr string `json:"supervisorAddress"`
	GossipAddr     string `json:"gossipAddress"`
	ClusteringAddr string `json:"clusteringAddress"`

	Labels labels.Set `json:"labels"`

	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Version   version.Version `json:"version"`
}
