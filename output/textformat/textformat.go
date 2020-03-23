package textformat

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gosuri/uitable"

	"code.storageos.net/storageos/c2-cli/pkg/health"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
)

var (
	nodeHeaders      = []interface{}{"NAME", "HEALTH", "AGE", "LABELS"}
	namespaceHeaders = []interface{}{"NAME", "AGE"}
	userHeaders      = []interface{}{"NAME", "ROLE", "AGE", "GROUPS"}
	volumeHeaders    = []interface{}{"NAMESPACE", "NAME", "SIZE", "LOCATION", "ATTACHED ON", "REPLICAS", "AGE"}
)

// Displayer is a type which creates human-readable strings and writes them to
// io.Writers.
type Displayer struct {
	timeHumanizer output.TimeHumanizer
}

// -----------------------------------------------------------------------------
// CREATE
// -----------------------------------------------------------------------------

// CreateUser builds a human friendly representation of resource, writing the
// result to w.
func (d *Displayer) CreateUser(ctx context.Context, w io.Writer, user *output.User) error {
	table, write := createTable(userHeaders)
	d.printUser(table, user)
	return write(w)
}

// CreateVolume builds a human friendly string from volume, writing the result to w.
func (d *Displayer) CreateVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	table, write := createTable(volumeHeaders)
	d.printVolume(table, volume)
	return write(w)
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
func (d *Displayer) GetVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	table, write := createTable(volumeHeaders)
	d.printVolume(table, volume)
	return write(w)
}

// GetListVolumes creates human-readable strings, writing the result to w.
func (d *Displayer) GetListVolumes(ctx context.Context, w io.Writer, resources []*output.Volume) error {
	table, write := createTable(volumeHeaders)
	for _, vol := range resources {
		d.printVolume(table, vol)
	}
	return write(w)
}

func (d *Displayer) printVolume(table *uitable.Table, vol *output.Volume) {
	location := fmt.Sprintf("%s (%s)", vol.Master.NodeName, vol.Master.Health)

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

	table.AddRow(vol.NamespaceName, vol.Name, size, location, vol.AttachedOnName, replicas, age)
}

func (d *Displayer) printNode(table *uitable.Table, node *node.Resource) {
	age := d.timeHumanizer.TimeToHuman(node.CreatedAt)
	table.AddRow(node.Name, node.Health.String(), age, node.Labels.String())
}

func (d *Displayer) printUser(table *uitable.Table, user *output.User) {
	age := d.timeHumanizer.TimeToHuman(user.CreatedAt)

	role := "user"
	if user.IsAdmin {
		role = "admin"
	}

	groupNames := make([]string, 0, len(user.Groups))
	for _, g := range user.Groups {
		groupNames = append(groupNames, g.Name)
	}

	table.AddRow(user.Username, role, age, strings.Join(groupNames, ","))
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

// AttachVolume writes a success message in the writer
func (d *Displayer) AttachVolume(ctx context.Context, w io.Writer) error {
	_, err := fmt.Fprintln(w, "volume attached")
	return err
}

// DetachVolume writes a success message to the writer
func (d *Displayer) DetachVolume(ctx context.Context, w io.Writer) error {
	_, err := fmt.Fprintln(w, "volume detached")
	return err
}

// NewDisplayer initialises a Displayer which prints human readable strings
// StorageOS to output CLI results.
func NewDisplayer(timeHumanizer output.TimeHumanizer) *Displayer {
	return &Displayer{
		timeHumanizer: timeHumanizer,
	}
}
