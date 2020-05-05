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
