package textformat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/dustin/go-humanize"
	"github.com/gosuri/uitable"

	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/capacity"
	"code.storageos.net/storageos/c2-cli/pkg/health"
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
func (d *Displayer) GetCluster(ctx context.Context, w io.Writer, resource *output.Cluster) error {
	table, write := createTable(nil)

	table.AddRow("ID:", resource.ID)
	table.AddRow("Licence:", "")
	table.AddRow("  expiration:", d.timeToHuman(resource.Licence.ExpiresAt))
	table.AddRow("  capacity:", humanize.IBytes(resource.Licence.ClusterCapacityBytes))
	table.AddRow("  kind:", resource.Licence.Kind)
	table.AddRow("  customer name:", resource.Licence.CustomerName)
	table.AddRow("Created at:", d.timeToHuman(resource.CreatedAt))
	table.AddRow("Updated at:", d.timeToHuman(resource.UpdatedAt))

	return write(w)
}

// GetDiagnostics writes a success message displaying outputPath to w.
func (d *Displayer) GetDiagnostics(ctx context.Context, w io.Writer, outputPath string) error {
	_, err := fmt.Fprintf(w, "Diagnostic bundle written to %v \n", outputPath)
	return err
}

// GetNode creates human-readable strings, writing the result to w.
func (d *Displayer) GetNode(ctx context.Context, w io.Writer, resource *output.Node) error {
	table, write := createTable(nodeHeaders)
	d.printNode(table, resource)
	return write(w)
}

// GetListNodes creates human-readable strings, writing the result to w.
func (d *Displayer) GetListNodes(ctx context.Context, w io.Writer, resources []*output.Node) error {
	table, write := createTable(nodeHeaders)
	for _, r := range resources {
		d.printNode(table, r)
	}
	return write(w)
}

// GetNamespace creates human-readable strings, writing the result to w.
func (d *Displayer) GetNamespace(ctx context.Context, w io.Writer, resource *output.Namespace) error {
	table, write := createTable(namespaceHeaders)

	// Humanized
	age := d.timeHumanizer.TimeToHuman(resource.CreatedAt)

	table.AddRow(resource.Name, age)

	return write(w)
}

// GetListNamespaces creates human-readable strings, writing the result to w.
func (d *Displayer) GetListNamespaces(ctx context.Context, w io.Writer, resources []*output.Namespace) error {
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

func (d *Displayer) printNode(table *uitable.Table, node *output.Node) {
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

// DescribeNode prints all the detailed information about a node
func (d *Displayer) DescribeNode(ctx context.Context, w io.Writer, node *output.NodeDescription) error {
	return d.describeNode(ctx, w, node)
}

// DescribeListNodes prints all the detailed information about a list of nodes
func (d *Displayer) DescribeListNodes(ctx context.Context, w io.Writer, nodes []*output.NodeDescription) error {
	for i, node := range nodes {
		if i != 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}

		if err := d.describeNode(ctx, w, node); err != nil {
			return err
		}
	}

	return nil
}

func (d *Displayer) describeNode(ctx context.Context, w io.Writer, node *output.NodeDescription) error {
	table := uitable.New()
	table.MaxColWidth = 80
	table.Separator = "  "

	table.AddRow("ID", node.ID.String())
	table.AddRow("Name", node.Name)
	table.AddRow("Health", node.Health.String())
	// Addresses
	table.AddRow("Addresses:")
	table.AddRow("  Data Transfer address", node.IOAddr)
	table.AddRow("  Gossip address", node.GossipAddr)
	table.AddRow("  Supervisor address", node.SupervisorAddr)
	table.AddRow("  Clustering address", node.ClusteringAddr)

	table.AddRow("Labels", node.Labels.String())
	table.AddRow("Created at", d.timeToHuman(node.CreatedAt))
	table.AddRow("Updated at", d.timeToHuman(node.UpdatedAt))
	table.AddRow("Version", node.Version.String())

	capacityStats := "n/a"
	if node.Capacity != (capacity.Stats{}) {
		capacityStats = fmt.Sprintf(
			"%s/%s (%s in use)",
			humanize.IBytes(node.Capacity.Available),
			humanize.IBytes(node.Capacity.Total),
			humanize.IBytes(node.Capacity.Total-node.Capacity.Free),
		)
	}
	table.AddRow("Available capacity", capacityStats)

	// Volumes
	if len(node.HostedVolumes) > 0 {
		table.AddRow()
		table.AddRow("Local volume deployments:")
		table.AddRow("  DEPLOYMENT ID", "VOLUME", "NAMESPACE", "HEALTH", "TYPE", "SIZE")
		for _, vol := range node.HostedVolumes {
			table.AddRow(
				"  "+vol.LocalDeployment.ID,
				vol.Name,
				vol.NamespaceName,
				vol.LocalDeployment.Health,
				vol.LocalDeployment.Kind,
				humanize.IBytes(vol.SizeBytes),
			)
		}
	}

	_, err := fmt.Fprintln(w, table)
	return err
}

// DescribeVolume prints in the output writer a tabular representation, in a key
// value shape, of all details about a volumes and its master and replicas.
func (d *Displayer) DescribeVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	return d.describeVolume(ctx, w, volume)
}

// DescribeListVolumes writes a detailed, yet human-friendly table
// representation to w for each item in volumes.
func (d *Displayer) DescribeListVolumes(ctx context.Context, w io.Writer, volumes []*output.Volume) error {
	for i, v := range volumes {
		if i > 0 {
			_, err := w.Write([]byte("\n"))
			if err != nil {
				return err
			}
		}

		err := d.describeVolume(ctx, w, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Displayer) describeVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	table := uitable.New()
	table.Separator = "  "

	table.AddRow("ID", volume.ID.String())
	table.AddRow("Name", volume.Name)
	table.AddRow("Description", volume.Description)

	var attachedOnString string
	if volume.AttachedOn != "" {
		attachedOnString = fmt.Sprintf("%s (%s)", volume.AttachedOnName, volume.AttachedOn)
	}

	table.AddRow("AttachedOn", attachedOnString)
	table.AddRow("Namespace", fmt.Sprintf("%s (%s)", volume.NamespaceName, volume.Namespace))
	table.AddRow("Labels", volume.Labels.String())
	table.AddRow("FileSystem", volume.Filesystem.String())
	table.AddRow("Size", volume.SizeBytes)

	table.AddRow("Version", volume.Version)
	table.AddRow("Created at", d.timeToHuman(volume.CreatedAt))
	table.AddRow("Updated at", d.timeToHuman(volume.UpdatedAt))

	table.AddRow("", "")
	table.AddRow("Master:")
	d.describeMaster(table, volume.Master)

	table.AddRow("", "")
	table.AddRow("Replicas:")
	for i, rep := range volume.Replicas {
		if i > 0 {
			table.AddRow("", "")
		}
		d.describeReplica(table, volume.SizeBytes, rep)
	}

	_, err := fmt.Fprintln(w, table)
	return err
}

func (d *Displayer) describeMaster(table *uitable.Table, master *output.Deployment) {
	table.AddRow("  ID", master.ID.String())
	table.AddRow("  Node", fmt.Sprintf("%s (%s)", master.NodeName, master.Node))
	table.AddRow("  Health", master.Health.String())
}

func (d *Displayer) describeReplica(table *uitable.Table, size uint64, replica *output.Deployment) {
	table.AddRow("  ID", replica.ID.String())
	table.AddRow("  Node", fmt.Sprintf("%s (%s)", replica.NodeName, replica.Node))
	table.AddRow("  Health", replica.Health.String())
	table.AddRow("  Promotable", replica.Promotable)
	if replica.Health == health.ReplicaSyncing {
		barStr, err := syncProgressBarString(
			replica.SyncProgress.BytesRemaining,
			size,
			replica.SyncProgress.EstimatedSecondsRemaining,
		)
		if err != nil {
			recap := fmt.Sprintf("%d/%d", size-replica.SyncProgress.BytesRemaining, size)
			table.AddRow("  Sync Progress", recap)
		} else {
			table.AddRow("  Sync Progress", barStr)
		}
	}
}

const format pb.ProgressBarTemplate = `{{counters . }} {{bar . "[" "#" "#" "." "]"}} {{percent . }}  -  {{string . "suffix"}}`

func syncProgressBarString(current, max, secondsRemaining uint64) (string, error) {
	var maxInt64 uint64 = math.MaxInt64
	if current > maxInt64 || max > maxInt64 {
		return "", errors.New("invalid sync progress value received: int64 overflow")
	}

	if current > max {
		return "", errors.New("invalid sync progress value received: current < max")
	}

	bar := format.Start64(int64(max))
	bar.Set(pb.Bytes, true)
	bar.SetWidth(80)
	bar.SetCurrent(int64(max - current))

	etaString := time.Duration(secondsRemaining) * time.Second
	bar.Set("suffix", fmt.Sprintf("ETA: %s", etaString))
	return bar.String(), nil
}

func (d *Displayer) timeToHuman(t time.Time) string {
	humanized := d.timeHumanizer.TimeToHuman(t)
	rfc := t.Format(time.RFC3339)
	return fmt.Sprintf("%s (%s)", rfc, humanized)
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
