package output

import (
	"time"

	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/volume"
)

// Volume defines a type that includes all the info that a volume should have to
// be outputted
type Volume struct {
	ID            id.Volume       `json:"id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	AttachedOn    id.Node         `json:"attachedOn"`
	Namespace     id.Namespace    `json:"namespaceID"`
	NamespaceName string          `json:"namespaceName"`
	Labels        labels.Set      `json:"labels"`
	Filesystem    volume.FsType   `json:"filesystem"`
	SizeBytes     uint64          `json:"sizeBytes"`
	Master        *Deployment     `json:"master"`
	Replicas      []*Deployment   `json:"replicas"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	Version       version.Version `json:"version"`
}

// Deployment defines a type that includes all the info that a deployment should
// have to be outputted
type Deployment struct {
	ID           id.Deployment      `json:"id"`
	Node         id.Node            `json:"nodeID"`
	Health       health.VolumeState `json:"health"`
	Promotable   bool               `json:"promotable"`
	SyncProgress *SyncProgress      `json:"syncProgress,omitempty"`
}

// SyncProgress defines a type that includes all the info that a SyncProgress
// should have to be outputted
type SyncProgress struct {
	BytesRemaining            uint64 `json:"bytesRemaining"`
	ThroughputBytes           uint64 `json:"throughputBytes"`
	EstimatedSecondsRemaining uint64 `json:"estimatedSecondsRemaining"`
}

// NewVolume creates a new Volume object for output purpose
func NewVolume(vol *volume.Resource, ns *namespace.Resource) *Volume {
	return &Volume{
		ID:            vol.ID,
		Name:          vol.Name,
		Description:   vol.Description,
		AttachedOn:    vol.AttachedOn,
		Namespace:     vol.Namespace,
		NamespaceName: ns.Name,
		Labels:        vol.Labels,
		Filesystem:    vol.Filesystem,
		SizeBytes:     vol.SizeBytes,
		Master:        newDeployment(vol.Master),
		Replicas:      newDeployments(vol.Replicas),
		CreatedAt:     vol.CreatedAt,
		UpdatedAt:     vol.UpdatedAt,
		Version:       vol.Version,
	}
}

func newDeployment(dep *volume.Deployment) *Deployment {
	outputDep := &Deployment{
		ID:         dep.ID,
		Node:       dep.Node,
		Health:     dep.Health,
		Promotable: dep.Promotable,
	}

	if dep.SyncProgress != nil {
		outputDep.SyncProgress = newSyncProgress(dep.SyncProgress)
	}

	return outputDep
}

func newDeployments(deployments []*volume.Deployment) []*Deployment {
	deps := make([]*Deployment, 0, len(deployments))
	for _, d := range deployments {
		deps = append(deps, newDeployment(d))
	}
	return deps
}

func newSyncProgress(sync *volume.SyncProgress) *SyncProgress {
	return &SyncProgress{
		BytesRemaining:            sync.BytesRemaining,
		ThroughputBytes:           sync.ThroughputBytes,
		EstimatedSecondsRemaining: sync.EstimatedSecondsRemaining,
	}
}
