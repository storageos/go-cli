package types

import (
	"time"

	apiTypes "github.com/storageos/go-api/types"
)

// ClusterCreateOps - optional fields when creating cluster
type ClusterCreateOps struct {
	AccountID string
	// optional value when to expire cluster
	TTL  int64
	Name string
	Size int
}

type Cluster struct {
	// cluster ID used for joining or getting cluster status
	ID string `json:"id,omitempty"`

	// cluster size, defaults to 3
	Size int `json:"size,omitempty"`

	Name string `json:"name,omitempty"`

	// optional account ID
	AccountID string `json:"accountID,omitempty"`

	// nodes participating in cluster
	Nodes []*Node `json:"nodes,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type Node struct {
	ID               string `json:"id,omitempty"` // node/controller UUID
	Name             string `json:"name,omitempty"`
	AdvertiseAddress string `json:"advertiseAddress,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`

	Health struct {
		CP *apiTypes.CPHealthStatus
		DP *apiTypes.DPHealthStatus
	}
}

// NodeByName sorts node list by hostname
type NodeByName []*Node

func (n NodeByName) Len() int           { return len(n) }
func (n NodeByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n NodeByName) Less(i, j int) bool { return n[i].Name < n[j].Name }
