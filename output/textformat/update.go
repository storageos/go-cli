package textformat

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/gosuri/uitable"

	"code.storageos.net/storageos/c2-cli/pkg/size"

	"code.storageos.net/storageos/c2-cli/pkg/health"

	"code.storageos.net/storageos/c2-cli/pkg/labels"

	"code.storageos.net/storageos/c2-cli/output"
)

// UpdateLicence prints all the detailed information about a new licence, after
// it has been correctly updated.
func (d *Displayer) UpdateLicence(ctx context.Context, w io.Writer, licence *output.Licence) error {
	fmt.Fprintf(w, "Licence applied to cluster %s.\n\n", licence.ClusterID.String())

	table, write := createTable(nil)

	table.AddRow("Expiration:", d.timeToHuman(licence.ExpiresAt))
	table.AddRow("Capacity:", size.Format(licence.ClusterCapacityBytes))
	table.AddRow("Used:", size.Format(licence.UsedBytes))
	table.AddRow("Kind:", licence.Kind)
	table.AddRow("Customer name:", licence.CustomerName)

	return write(w)
}

// SetReplicas prints a user friendly message denoting that the target
// replica num has been updated
func (d *Displayer) SetReplicas(ctx context.Context, w io.Writer, new uint64) error {
	if _, err := fmt.Fprintf(w, "Target replica number accepted, converging to %d", new); err != nil {
		return err
	}
	return nil
}

// ResizeVolume prints a user friendly message denoting that the target size
// has been updated
func (d *Displayer) ResizeVolume(ctx context.Context, w io.Writer, updatedVol output.VolumeUpdate) error {
	err := printVolumeUpdate(w, updatedVol)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "\nVolume %s (%s) updated. Size changed.\n",
		updatedVol.Name,
		updatedVol.ID,
	)

	return err
}

// UpdateVolumeDescription writes a success message to w with the updated description.
func (d *Displayer) UpdateVolumeDescription(ctx context.Context, w io.Writer, updatedVol output.VolumeUpdate) error {
	err := printVolumeUpdate(w, updatedVol)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "\nVolume %s (%s) updated. Description changed.\n",
		updatedVol.Name,
		updatedVol.ID,
	)

	return err
}

// UpdateVolumeLabels writes a success message to w with the updated labels.
func (d *Displayer) UpdateVolumeLabels(ctx context.Context, w io.Writer, updatedVol output.VolumeUpdate) error {
	err := printVolumeUpdate(w, updatedVol)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "\nVolume %s (%s) updated. Labels changed.\n",
		updatedVol.Name,
		updatedVol.ID,
	)

	return err
}

func printVolumeUpdate(w io.Writer, updateVol output.VolumeUpdate) error {
	table, write := createTable(nil)
	table.AddRow("Name:", updateVol.Name)
	table.AddRow("ID:", updateVol.ID)
	table.AddRow("Size:", size.Format(updateVol.SizeBytes))
	table.AddRow("Description:", updateVol.Description)
	table.AddRow("AttachedOn:", updateVol.AttachedOn)
	table.AddRow("Replicas:", deploymentsToString(updateVol.Replicas))
	addLabels(table, updateVol.Labels)
	return write(w)
}

func addLabels(table *uitable.Table, lab labels.Set) {
	table.AddRow("Labels:", "")

	keys := []string{}
	for k := range lab {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		table.AddRow("  - "+k, lab[k])
	}
}

func deploymentsToString(deps []*output.VolumeUpdateDeployment) string {
	var (
		recovering,
		provisioning,
		provisioned,
		syncing,
		ready,
		deleted,
		failed,
		unknown int
	)
	for _, d := range deps {

		switch d.Health {
		case health.ReplicaRecovering:
			recovering++
		case health.ReplicaProvisioning:
			provisioning++
		case health.ReplicaProvisioned:
			provisioned++
		case health.ReplicaSyncing:
			syncing++
		case health.ReplicaReady:
			ready++
		case health.ReplicaDeleted:
			deleted++
		case health.ReplicaFailed:
			failed++
		case health.ReplicaUnknown:
			unknown++
		}
	}

	s := []string{}
	if recovering > 0 {
		s = append(s, fmt.Sprintf("%dx %s", recovering, health.ReplicaRecovering.String()))
	}
	if provisioning > 0 {
		s = append(s, fmt.Sprintf("%dx %s", provisioning, health.ReplicaProvisioning.String()))
	}
	if provisioned > 0 {
		s = append(s, fmt.Sprintf("%dx %s", provisioned, health.ReplicaProvisioned.String()))
	}
	if syncing > 0 {
		s = append(s, fmt.Sprintf("%dx %s", syncing, health.ReplicaSyncing.String()))
	}
	if ready > 0 {
		s = append(s, fmt.Sprintf("%dx %s", ready, health.ReplicaReady.String()))
	}
	if deleted > 0 {
		s = append(s, fmt.Sprintf("%dx %s", deleted, health.ReplicaDeleted.String()))
	}
	if failed > 0 {
		s = append(s, fmt.Sprintf("%dx %s", failed, health.ReplicaFailed.String()))
	}
	if unknown > 0 {
		s = append(s, fmt.Sprintf("%dx %s", unknown, health.ReplicaUnknown.String()))
	}

	return strings.Join(s, ", ")
}
