package output

import (
	"time"

	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/volume"
)

// Volume defines a type that includes all the info that a volume should have to
// be outputted
type Volume struct {
	ID             id.Volume       `json:"id" yaml:"id"`
	Name           string          `json:"name" yaml:"name"`
	Description    string          `json:"description" yaml:"description"`
	AttachedOn     id.Node         `json:"attachedOn" yaml:"attachedOn"`
	AttachedOnName string          `json:"attachedOnName" yaml:"attachedOnName"`
	Namespace      id.Namespace    `json:"namespaceID" yaml:"namespaceID"`
	NamespaceName  string          `json:"namespaceName" yaml:"namespaceName"`
	Labels         labels.Set      `json:"labels" yaml:"labels"`
	Filesystem     volume.FsType   `json:"filesystem" yaml:"filesystem"`
	SizeBytes      uint64          `json:"sizeBytes" yaml:"sizeBytes"`
	Master         *Deployment     `json:"master" yaml:"master"`
	Replicas       []*Deployment   `json:"replicas" yaml:"replicas"`
	CreatedAt      time.Time       `json:"createdAt" yaml:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt" yaml:"updatedAt"`
	Version        version.Version `json:"version" yaml:"version"`
}

// Deployment defines a type that includes all the info that a deployment should
// have to be outputted
type Deployment struct {
	ID           id.Deployment      `json:"id" yaml:"id"`
	Node         id.Node            `json:"nodeID" yaml:"nodeID"`
	NodeName     string             `json:"nodeName" yaml:"nodeName"`
	Health       health.VolumeState `json:"health" yaml:"health"`
	Promotable   bool               `json:"promotable" yaml:"promotable"`
	SyncProgress *SyncProgress      `json:"syncProgress,omitempty" yaml:"syncProgress,omitempty"`
}

// SyncProgress defines a type that includes all the info that a SyncProgress
// should have to be outputted
type SyncProgress struct {
	BytesRemaining            uint64 `json:"bytesRemaining" yaml:"bytesRemaining"`
	ThroughputBytes           uint64 `json:"throughputBytes" yaml:"throughputBytes"`
	EstimatedSecondsRemaining uint64 `json:"estimatedSecondsRemaining" yaml:"estimatedSecondsRemaining"`
}

// VolumeDeletion defines a volume deletion confirmation output representation.
type VolumeDeletion struct {
	ID        id.Volume    `json:"id" yaml:"id"`
	Namespace id.Namespace `json:"namespaceID" yaml:"namespaceID"`
}

// NewVolume creates a new Volume output representation using extra details
// from the provided parameters.
func NewVolume(vol *volume.Resource, ns *namespace.Resource, nodes map[id.Node]*node.Resource) *Volume {
	var attachedOnName string

	attachedOn, ok := nodes[vol.AttachedOn]
	switch {
	case ok:
		attachedOnName = attachedOn.Name
	case vol.AttachedOn == "":
	default:
		attachedOnName = unknownResourceName
	}

	return &Volume{
		ID:             vol.ID,
		Name:           vol.Name,
		Description:    vol.Description,
		AttachedOn:     vol.AttachedOn,
		AttachedOnName: attachedOnName,
		Namespace:      vol.Namespace,
		NamespaceName:  ns.Name,
		Labels:         vol.Labels,
		Filesystem:     vol.Filesystem,
		SizeBytes:      vol.SizeBytes,
		Master:         newDeployment(vol.Master, nodes),
		Replicas:       newDeployments(vol.Replicas, nodes),
		CreatedAt:      vol.CreatedAt,
		UpdatedAt:      vol.UpdatedAt,
		Version:        vol.Version,
	}
}

func newDeployment(dep *volume.Deployment, nodes map[id.Node]*node.Resource) *Deployment {
	nodeName := unknownResourceName

	n, ok := nodes[dep.Node]
	if ok {
		nodeName = n.Name
	}

	outputDep := &Deployment{
		ID:         dep.ID,
		Node:       dep.Node,
		NodeName:   nodeName,
		Health:     dep.Health,
		Promotable: dep.Promotable,
	}

	// This field is expected to be empty in a lot of cases, so check first.
	if dep.SyncProgress != nil {
		outputDep.SyncProgress = newSyncProgress(dep.SyncProgress)
	}

	return outputDep
}

func newDeployments(deployments []*volume.Deployment, nodes map[id.Node]*node.Resource) []*Deployment {
	outputDeployments := make([]*Deployment, 0, len(deployments))
	for _, d := range deployments {
		outputDeployments = append(outputDeployments, newDeployment(d, nodes))
	}
	return outputDeployments
}

func newSyncProgress(sync *volume.SyncProgress) *SyncProgress {
	return &SyncProgress{
		BytesRemaining:            sync.BytesRemaining,
		ThroughputBytes:           sync.ThroughputBytes,
		EstimatedSecondsRemaining: sync.EstimatedSecondsRemaining,
	}
}

// NewVolumeDeletion constructs a volume deletion confirmation output representation.
func NewVolumeDeletion(volumeID id.Volume, namespaceID id.Namespace) VolumeDeletion {
	return VolumeDeletion{
		ID:        volumeID,
		Namespace: namespaceID,
	}
}
