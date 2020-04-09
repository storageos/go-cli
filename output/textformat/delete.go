package textformat

import (
	"context"
	"fmt"
	"io"

	"code.storageos.net/storageos/c2-cli/output"
)

// DeleteUser writes a message containing the user deletion confirmation to w.
func (d *Displayer) DeleteUser(ctx context.Context, w io.Writer, confirmation output.UserDeletion) error {
	_, err := fmt.Fprintf(w, "deleted user %s\n", confirmation.ID.String())
	return err
}

// DeleteVolume writes a message containing the volume deletion confirmation
// to w.
func (d *Displayer) DeleteVolume(ctx context.Context, w io.Writer, confirmation output.VolumeDeletion) error {
	_, err := fmt.Fprintf(w, "deleted volume %v from namespace %v\n", confirmation.ID, confirmation.Namespace)
	return err
}

// DeleteNamespace writes a message containing the namespace deletion
// confirmation to w.
func (d *Displayer) DeleteNamespace(ctx context.Context, w io.Writer, confirmation output.NamespaceDeletion) error {
	_, err := fmt.Fprintf(w, "deleted namespace %s\n", confirmation.ID.String())
	return err
}

// DeleteVolumeAsync writes a successful request submission string to w.
func (d *Displayer) DeleteVolumeAsync(ctx context.Context, w io.Writer) error {
	_, err := fmt.Fprintln(w, "volume deletion request accepted")
	return err
}

// DeletePolicyGroup encodes the policy group deletion confirmation as JSON, writing
// the result to w
func (d *Displayer) DeletePolicyGroup(ctx context.Context, w io.Writer, confirmation output.PolicyGroupDeletion) error {
	_, err := fmt.Fprintf(w, "deleted policy group %s\n", confirmation.ID.String())
	return err
}
