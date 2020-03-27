package output

import (
	"fmt"

	"code.storageos.net/storageos/c2-cli/pkg/id"
)

// MissingRequiredNodeErr is an error type which indicates that an output could
// not be constructed because it is missing details of the given node.
type MissingRequiredNodeErr struct {
	uid id.Node
}

func (e MissingRequiredNodeErr) Error() string {
	return fmt.Sprintf("missing required details for node with id %v", e.uid)
}

// NewMissingRequiredNodeErr returns a new error indicating that the required
// details for node with uid are missing.
func NewMissingRequiredNodeErr(uid id.Node) MissingRequiredNodeErr {
	return MissingRequiredNodeErr{
		uid: uid,
	}
}

// MissingRequiredPolicyGroupErr is an error type which indicates that an
// output could not be constructed because it is missing details of the given
// policy group.
type MissingRequiredPolicyGroupErr struct {
	uid id.PolicyGroup
}

func (e MissingRequiredPolicyGroupErr) Error() string {
	return fmt.Sprintf("missing required details for policy group with id %v", e.uid)
}

// NewMissingRequiredPolicyGroupErr returns a new error indicating that the
// required details for policy group with uid are missing.
func NewMissingRequiredPolicyGroupErr(uid id.PolicyGroup) MissingRequiredPolicyGroupErr {
	return MissingRequiredPolicyGroupErr{
		uid: uid,
	}
}

// NodeDoesNotHostVolumeErr is an error type indicating that the given node
// does not host a local deployment for volID.
type NodeDoesNotHostVolumeErr struct {
	nodeID id.Node
	volID  id.Volume
}

func (e NodeDoesNotHostVolumeErr) Error() string {
	return fmt.Sprintf("node with id %v does not host volume with id %v", e.nodeID, e.volID)
}

// NewNodeDoesNotHostVolumeErr returns a new error indicating that nodeID does not have a
// local deployment for volID.
func NewNodeDoesNotHostVolumeErr(nodeID id.Node, volID id.Volume) NodeDoesNotHostVolumeErr {
	return NodeDoesNotHostVolumeErr{
		nodeID: nodeID,
		volID:  volID,
	}
}
