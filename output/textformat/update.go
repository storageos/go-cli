package textformat

import (
	"context"
	"fmt"
	"io"

	"github.com/dustin/go-humanize"

	"code.storageos.net/storageos/c2-cli/output"
)

// UpdateLicence prints all the detailed information about a new licence, after
// it has been correctly updated.
func (d *Displayer) UpdateLicence(ctx context.Context, w io.Writer, licence *output.Licence) error {
	fmt.Fprintf(w, "Licence applied to cluster %s.\n\n", licence.ClusterID.String())

	table, write := createTable(nil)

	table.AddRow("Expiration:", d.timeToHuman(licence.ExpiresAt))
	table.AddRow("Capacity:", humanize.IBytes(licence.ClusterCapacityBytes))
	table.AddRow("Used:", humanize.IBytes(licence.UsedBytes))
	table.AddRow("Kind:", licence.Kind)
	table.AddRow("Customer name:", licence.CustomerName)

	return write(w)
}

// UpdateVolume prints a user friendly message showing the successful
// mutation
func (d *Displayer) UpdateVolume(ctx context.Context, w io.Writer, updatedVol output.VolumeUpdate) error {

	// if we have no updated volume data show a generic msg
	if updatedVol.ID == "" {
		fmt.Fprintf(w, "Update successful")
		return nil
	}

	fmt.Fprintf(w, "Update successful for volume %s.\n\n", updatedVol.Name)

	table, write := createTable(nil)

	table.AddRow("Name:", updatedVol.Name)
	table.AddRow("Description:", updatedVol.Description)
	table.AddRow("Namespace:", updatedVol.Namespace)
	table.AddRow("Labels:", updatedVol.Labels)

	return write(w)
}

// SetReplicas prints a user friendly message denoting that the target
// replica num has been updated
func (d *Displayer) SetReplicas(ctx context.Context, w io.Writer, new uint64) error {

	if _, err := fmt.Fprintf(w, "Target replica number accepted, it is now %d", new); err != nil {
		return err
	}

	return nil
}
