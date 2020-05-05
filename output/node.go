package output

import (
	"time"

	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/capacity"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/volume"
)

// Node defines a type that contains all the info we need to output a node.
type Node struct {
	ID             id.Node          `json:"id" yaml:"id"`
	Name           string           `json:"name" yaml:"name"`
	Health         health.NodeState `json:"health" yaml:"health"`
	Capacity       capacity.Stats   `json:"capacity,omitempty" yaml:"capacity,omitempty"`
	IOAddr         string           `json:"ioAddress" yaml:"ioAddress"`
	SupervisorAddr string           `json:"supervisorAddress" yaml:"supervisorAddress"`
	GossipAddr     string           `json:"gossipAddress" yaml:"gossipAddr"`
	ClusteringAddr string           `json:"clusteringAddress" yaml:"clusteringAddr"`
	Labels         labels.Set       `json:"labels" yaml:"labels"`
	CreatedAt      time.Time        `json:"createdAt" yaml:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt" yaml:"updatedAt"`
	Version        version.Version  `json:"version" yaml:"version"`
}

// NodeDescription decorates a Node's output representation with additional
// details.
type NodeDescription struct {
	Node          `yaml:",inline"`
	HostedVolumes []*HostedVolume `json:"volumes" yaml:"volumes"`
}

// HostedVolume encapsulates a volume with the local deployment for a node.
type HostedVolume struct {
	ID              id.Volume       `json:"id" yaml:"id"`
	Name            string          `json:"name" yaml:"name"`
	Description     string          `json:"description" yaml:"description"`
	Namespace       id.Namespace    `json:"namespaceID" yaml:"namespaceID"`
	NamespaceName   string          `json:"namespaceName" yaml:"namespaceName"`
	Labels          labels.Set      `json:"labels" yaml:"labels"`
	Filesystem      volume.FsType   `json:"filesystem" yaml:"filesystem"`
	SizeBytes       uint64          `json:"sizeBytes" yaml:"sizeBytes"`
	LocalDeployment LocalDeployment `json:"localDeployment" yaml:"localDeployment"`
	CreatedAt       time.Time       `json:"createdAt" yaml:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt" yaml:"updatedAt"`
	Version         version.Version `json:"version" yaml:"version"`
}

// LocalDeployment contains information about the local deployment for a volume
// hosted by a node.
type LocalDeployment struct {
	ID           id.Deployment      `json:"id" yaml:"id"`
	Kind         string             `json:"kind" yaml:"kind"`
	Health       health.VolumeState `json:"health" yaml:"health"`
	Promotable   bool               `json:"promotable" yaml:"promotable"`
	SyncProgress *SyncProgress      `json:"syncProgress,omitempty" yaml:"syncProgress,omitempty"`
}

// NewNode creates a new Node output representation containing all the info to
// be outputted.
func NewNode(n *node.Resource) *Node {
	return &Node{
		ID:             n.ID,
		Name:           n.Name,
		Health:         n.Health,
		Capacity:       n.Capacity,
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

// NewNodeDescription constructs a new NodeDescription output representation
// type for n, decorating it with details about all its local volume
// deployments which are in hostedVolumes.
func NewNodeDescription(n *node.Resource, hostedVolumes []*volume.Resource, namespaceForID map[id.Namespace]*namespace.Resource) *NodeDescription {

	nodeDescription := &NodeDescription{
		Node:          *NewNode(n),
		HostedVolumes: []*HostedVolume{},
	}

	for _, volResource := range hostedVolumes {
		outputVol, err := newHostedVolumeForNode(n, volResource, namespaceForID[volResource.Namespace])
		if err != nil {
			// The node does not host the volume, don't add it to the list
			continue
		}

		nodeDescription.HostedVolumes = append(nodeDescription.HostedVolumes, outputVol)
	}

	return nodeDescription
}

// newHostedVolumeForNode constructs a new HostedVolume output representation
// for volume vol on node n in namespace ns, given that the node has a
// deployment for the volume.
func newHostedVolumeForNode(n *node.Resource, vol *volume.Resource, ns *namespace.Resource) (*HostedVolume, error) {
	var deploy *LocalDeployment

	if vol.Master.Node == n.ID {
		deploy = &LocalDeployment{
			ID:         vol.Master.ID,
			Kind:       "master",
			Health:     vol.Master.Health,
			Promotable: vol.Master.Promotable,
		}
	}
	for _, r := range vol.Replicas {
		if r.Node == n.ID {
			deploy = &LocalDeployment{
				ID:         r.ID,
				Kind:       "replica",
				Health:     r.Health,
				Promotable: r.Promotable,
			}
			if r.SyncProgress != nil {
				deploy.SyncProgress = newSyncProgress(r.SyncProgress)
			}
			break
		}
	}
	if deploy == nil {
		return nil, NewNodeDoesNotHostVolumeErr(n.ID, vol.ID)
	}

	return &HostedVolume{
		ID:              vol.ID,
		Name:            vol.Name,
		Description:     vol.Description,
		Namespace:       vol.Namespace,
		NamespaceName:   ns.Name,
		Labels:          vol.Labels,
		Filesystem:      vol.Filesystem,
		SizeBytes:       vol.SizeBytes,
		LocalDeployment: *deploy,
		CreatedAt:       vol.CreatedAt,
		UpdatedAt:       vol.UpdatedAt,
		Version:         vol.Version,
	}, nil
}
