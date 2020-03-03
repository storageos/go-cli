package health

import "code.storageos.net/storageos/openapi"

// NodeState represents the health state in which a node could be
type NodeState string

// Are all States a node could be.
const (
	NodeOnline  NodeState = NodeState(openapi.NODEHEALTH_ONLINE)
	NodeOffline           = NodeState(openapi.NODEHEALTH_OFFLINE)
	NodeUnknown           = NodeState(openapi.NODEHEALTH_UNKNOWN)
)

// NodeFromString parses a string and return the matching node state.
// If the string is not recognized, Unknown is returned
func NodeFromString(s string) NodeState {
	switch s {
	case string(openapi.NODEHEALTH_ONLINE):
		return NodeOnline
	case string(openapi.NODEHEALTH_OFFLINE):
		return NodeOffline
	default:
		return NodeUnknown
	}
}

// String returns the string representation of the State
func (n NodeState) String() string {
	return string(n)
}
