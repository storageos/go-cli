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
