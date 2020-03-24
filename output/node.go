package output

import (
	"time"

	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

// Node defines a type that contains all the info we need to output a node.
type Node struct {
	ID             id.Node          `json:"id"`
	Name           string           `json:"name"`
	Health         health.NodeState `json:"health"`
	IOAddr         string           `json:"ioAddress"`
	SupervisorAddr string           `json:"supervisorAddress"`
	GossipAddr     string           `json:"gossipAddress"`
	ClusteringAddr string           `json:"clusteringAddress"`
	Labels         labels.Set       `json:"labels"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Version        version.Version  `json:"version"`
}

// NewNode creates a new Node output representation containing all the info to
// be outputted.
func NewNode(n *node.Resource) *Node {
	return &Node{
		ID:             n.ID,
		Name:           n.Name,
		Health:         n.Health,
		IOAddr:         n.IOAddr,
		SupervisorAddr: n.SupervisorAddr,
		GossipAddr:     n.GossipAddr,
		ClusteringAddr: n.ClusteringAddr,
		Labels:         n.Labels,
		CreatedAt:      n.CreatedAt,
		UpdatedAt:      n.UpdatedAt,
		Version:        n.Version,
	}
}

// NewNodes creates a new list of output representation of nodes containing all
// the info to be outputted.
func NewNodes(nodes []*node.Resource) []*Node {
	ns := make([]*Node, 0, len(nodes))
	for _, n := range nodes {
		ns = append(ns, NewNode(n))
	}
	return ns
}
