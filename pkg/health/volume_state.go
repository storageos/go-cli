package health

import (
	"code.storageos.net/storageos/openapi"
)

// VolumeState represents the health state in which a volume could be
type VolumeState string

// All States a node could be.
const (
	ReplicaRecovering   VolumeState = VolumeState(openapi.REPLICAHEALTH_RECOVERING)
	ReplicaProvisioning             = VolumeState(openapi.REPLICAHEALTH_PROVISIONING)
	ReplicaProvisioned              = VolumeState(openapi.REPLICAHEALTH_PROVISIONED)
	ReplicaSyncing                  = VolumeState(openapi.REPLICAHEALTH_SYNCING)
	ReplicaReady                    = VolumeState(openapi.REPLICAHEALTH_READY)
	ReplicaDeleted                  = VolumeState(openapi.REPLICAHEALTH_DELETED)
	ReplicaFailed                   = VolumeState(openapi.REPLICAHEALTH_FAILED)
	ReplicaUnknown                  = VolumeState(openapi.REPLICAHEALTH_UNKNOWN)
	MasterOnline                    = VolumeState(openapi.MASTERHEALTH_ONLINE)
	MasterOffline                   = VolumeState(openapi.MASTERHEALTH_OFFLINE)
	MasterUnknown                   = VolumeState(openapi.MASTERHEALTH_UNKNOWN)
)

// ReplicaFromString returns the replica State matching the string in input.
// If the string doesn't match any of the known state, unknown is returned.
func ReplicaFromString(s string) VolumeState {
	switch s {
	case string(openapi.REPLICAHEALTH_RECOVERING):
		return ReplicaRecovering
	case string(openapi.REPLICAHEALTH_PROVISIONING):
		return ReplicaProvisioning
	case string(openapi.REPLICAHEALTH_PROVISIONED):
		return ReplicaProvisioned
	case string(openapi.REPLICAHEALTH_SYNCING):
		return ReplicaSyncing
	case string(openapi.REPLICAHEALTH_READY):
		return ReplicaReady
	case string(openapi.REPLICAHEALTH_DELETED):
		return ReplicaDeleted
	case string(openapi.REPLICAHEALTH_FAILED):
		return ReplicaFailed
	default:
		return ReplicaUnknown
	}
}

// MasterFromString returns the master State matching the string in input.
// If the string doesn't match any of the known state, unknown is returned.
func MasterFromString(s string) VolumeState {
	switch s {
	case string(openapi.MASTERHEALTH_ONLINE):
		return MasterOnline
	case string(openapi.MASTERHEALTH_OFFLINE):
		return MasterOffline
	default:
		return MasterUnknown
	}
}

// String returns the string representation of the State
func (v VolumeState) String() string {
	return string(v)
}
