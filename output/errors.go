package output

import (
	"fmt"

	"code.storageos.net/storageos/c2-cli/pkg/id"
)

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
