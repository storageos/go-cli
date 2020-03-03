package textformat

import (
	"context"
	"fmt"
	"io"
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/health"

	"github.com/dustin/go-humanize"
	"github.com/gosuri/uitable"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/volume"
)

var (
	nodeHeaders      = []interface{}{"NAME", "HEALTH", "AGE", "LABELS"}
	namespaceHeaders = []interface{}{"NAME", "AGE"}
	volumeHeaders    = []interface{}{"NAMESPACE", "NAME", "SIZE", "LOCATION", "REPLICAS", "AGE"}
)

// Displayer is a type which creates human-readable strings and writes them to
// io.Writers.
type Displayer struct {
	timeHumanizer output.TimeHumanizer
}

// -----------------------------------------------------------------------------
// GET
// -----------------------------------------------------------------------------

// GetCluster creates human-readable strings, writing the result to w.
func (d *Displayer) GetCluster(ctx context.Context, w io.Writer, resource *cluster.Resource) error {
	table, write := createTable(nil)

	// Dates
	expiration := d.timeHumanizer.TimeToHuman(resource.Licence.ExpiresAt)
	expirationString := fmt.Sprintf("%s (%s)", resource.Licence.ExpiresAt.Format(time.RFC3339), expiration)

	age := d.timeHumanizer.TimeToHuman(resource.CreatedAt)
	ageString := fmt.Sprintf("%s (%s)", resource.CreatedAt.Format(time.RFC3339), age)

	table.AddRow("ID:", resource.ID)
	table.AddRow("Licence:", "")
	table.AddRow("  expiration:", expirationString)
	table.AddRow("  capacity:", humanize.IBytes(resource.Licence.ClusterCapacityBytes))
	table.AddRow("  kind:", resource.Licence.Kind)
	table.AddRow("  customer name:", resource.Licence.CustomerName)
	table.AddRow("Created At:", ageString)

	return write(w)
}

// GetNode creates human-readable strings, writing the result to w.
func (d *Displayer) GetNode(ctx context.Context, w io.Writer, resource *node.Resource) error {
	table, write := createTable(nodeHeaders)
	d.printNode(table, resource)
	return write(w)
}

// GetListNodes creates human-readable strings, writing the result to w.
func (d *Displayer) GetListNodes(ctx context.Context, w io.Writer, resources []*node.Resource) error {
	table, write := createTable(nodeHeaders)
	for _, r := range resources {
		d.printNode(table, r)
	}
	return write(w)
}

// GetNamespace creates human-readable strings, writing the result to w.
func (d *Displayer) GetNamespace(ctx context.Context, w io.Writer, resource *namespace.Resource) error {
	table, write := createTable(namespaceHeaders)

	// Humanized
	age := d.timeHumanizer.TimeToHuman(resource.CreatedAt)

	table.AddRow(resource.Name, age)

	return write(w)
}

// GetListNamespaces creates human-readable strings, writing the result to w.
func (d *Displayer) GetListNamespaces(ctx context.Context, w io.Writer, resources []*namespace.Resource) error {
	table, write := createTable(namespaceHeaders)

	for _, ns := range resources {
		// Humanized
		age := d.timeHumanizer.TimeToHuman(ns.CreatedAt)

		table.AddRow(ns.Name, age)
	}

	return write(w)
}

// GetVolume creates human-readable strings, writing the result to w.
func (d *Displayer) GetVolume(ctx context.Context, w io.Writer, resource *volume.Resource) error {
	table, write := createTable(volumeHeaders)
	d.printVolume(table, resource)
	return write(w)
}

// GetListVolumes creates human-readable strings, writing the result to w.
func (d *Displayer) GetListVolumes(ctx context.Context, w io.Writer, resources []*volume.Resource) error {
	table, write := createTable(volumeHeaders)
	for _, vol := range resources {
		d.printVolume(table, vol)
	}
	return write(w)
}

func (d *Displayer) printVolume(table *uitable.Table, vol *volume.Resource) {
	location := fmt.Sprintf("%s (%s)", vol.Master.Node, vol.Master.Health)

	// Replicas
	readyReplicas := 0
	for _, r := range vol.Replicas {
		if r.Health == health.ReplicaReady {
			readyReplicas++
		}
	}
	replicas := fmt.Sprintf("%d/%d", readyReplicas, len(vol.Replicas))

	// Humanized
	size := humanize.IBytes(vol.SizeBytes)
	age := d.timeHumanizer.TimeToHuman(vol.CreatedAt)

	table.AddRow(vol.Namespace, vol.Name, size, location, replicas, age)
}

func (d *Displayer) printNode(table *uitable.Table, node *node.Resource) {
	age := d.timeHumanizer.TimeToHuman(node.CreatedAt)
	table.AddRow(node.Name, node.Health.String(), age, node.Labels.String())
}

func createTable(headers []interface{}) (*uitable.Table, func(io.Writer) error) {
	table := uitable.New()
	table.MaxColWidth = 50
	table.Separator = "  "

	// header
	if headers != nil {
		table.AddRow(headers...)
	}

	return table, func(w io.Writer) error {
		_, err := fmt.Fprintln(w, table)
		return err
	}
}

// NewDisplayer initialises a Displayer which prints human readable strings
// StorageOS to output CLI results.
func NewDisplayer(timeHumanizer output.TimeHumanizer) *Displayer {
	return &Displayer{
		timeHumanizer: timeHumanizer,
	}
}
