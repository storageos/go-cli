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

// SetReplicas writes a success message to w.
func (d *Displayer) SetReplicas(ctx context.Context, w io.Writer) error {
	_, err := fmt.Fprintln(w, "request to change number of replicas accepted")
	return err
}

// UpdateVolumeDescription writes a success message to w.
func (d *Displayer) UpdateVolumeDescription(ctx context.Context, w io.Writer, volUpdate output.VolumeUpdate) error {
	_, err := fmt.Fprintf(w, "Volume %s (%s) updated.\nNew description: `%s`.", volUpdate.Name, volUpdate.ID, volUpdate.Description)
	return err
}
