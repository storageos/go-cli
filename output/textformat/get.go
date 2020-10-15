package textformat

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gosuri/uitable"

	"code.storageos.net/storageos/c2-cli/pkg/size"
	"code.storageos.net/storageos/c2-cli/volume"

	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/health"
)

// GetCluster creates human-readable strings, writing the result to w.
func (d *Displayer) GetCluster(ctx context.Context, w io.Writer, resource *output.Cluster) error {
	table, write := createTable(nil)

	table.AddRow("ID:", resource.ID)
	table.AddRow("Created at:", d.timeToHuman(resource.CreatedAt))
	table.AddRow("Updated at:", d.timeToHuman(resource.UpdatedAt))

	// nodes
	healthy, unhealthy := 0, 0
	for _, n := range resource.Nodes {
		if n.Health == health.NodeOnline {
			healthy++
		} else {
			unhealthy++
		}
	}

	table.AddRow("Nodes:", len(resource.Nodes))
	table.AddRow("  Healthy:", healthy)
	table.AddRow("  Unhealthy:", unhealthy)

	return write(w)
}

// GetLicence creates human-readable strings, writing the result to w.
func (d *Displayer) GetLicence(ctx context.Context, w io.Writer, l *output.Licence) error {
	table, write := createTable(nil)

	clusterCapacity := fmt.Sprintf("%s (%d)", size.Format(l.ClusterCapacityBytes), l.ClusterCapacityBytes)
	used := fmt.Sprintf("%s (%d)", size.Format(l.UsedBytes), l.UsedBytes)

	table.AddRow("ClusterID:", l.ClusterID)
	table.AddRow("Expiration:", d.timeToHuman(l.ExpiresAt))
	table.AddRow("Capacity:", clusterCapacity)
	table.AddRow("Used:", used)
	table.AddRow("Kind:", l.Kind)
	table.AddRow("Features:", fmt.Sprintf("%v", l.Features))
	table.AddRow("Customer name:", l.CustomerName)

	return write(w)
}

// GetDiagnostics writes a success message displaying outputPath to w.
func (d *Displayer) GetDiagnostics(ctx context.Context, w io.Writer, outputPath string) error {
	_, err := fmt.Fprintf(w, "Diagnostic bundle written to %v \n", outputPath)
	return err
}

// GetUser creates human-readable strings, writing the result to w.
func (d *Displayer) GetUser(ctx context.Context, w io.Writer, user *output.User) error {
	table, write := createTable(userHeaders)
	d.printUser(table, user)
	return write(w)
}

// GetUsers creates human-readable strings, writing the result to w.
func (d *Displayer) GetUsers(ctx context.Context, w io.Writer, users []*output.User) error {
	table, write := createTable(userHeaders)
	for _, u := range users {
		d.printUser(table, u)
	}
	return write(w)
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

func (d *Displayer) printNode(table *uitable.Table, node *output.Node) {
	age := d.timeHumanizer.TimeToHuman(node.CreatedAt)
	table.AddRow(node.Name, node.Health.String(), age, node.Labels.String())
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

	// Attempt to extract the target replica number from the volume label set.
	targetReplicas, err := strconv.ParseUint(vol.Labels[volume.LabelReplicas], 10, 0)
	if err != nil {
		// If this fails then the length of the returned volume replica set
		// provides the best estimate.
		targetReplicas = uint64(len(vol.Replicas))
	}

	replicas := fmt.Sprintf("%d/%d", readyReplicas, targetReplicas)

	// Humanized
	size := size.Format(vol.SizeBytes)
	age := d.timeHumanizer.TimeToHuman(vol.CreatedAt)

	table.AddRow(vol.NamespaceName, vol.Name, size, location, vol.AttachedOnName, replicas, age)
}

// GetPolicyGroup creates human-readable strings, writing the result to w.
func (d *Displayer) GetPolicyGroup(ctx context.Context, w io.Writer, group *output.PolicyGroup) error {
	table, write := createTable(policyGroupHeaders)

	// Humanized
	age := d.timeHumanizer.TimeToHuman(group.CreatedAt)
	table.AddRow(group.Name, len(group.Users), len(group.Specs), age)

	return write(w)
}

// GetListPolicyGroups creates human-readable strings, writing the result to w.
func (d *Displayer) GetListPolicyGroups(ctx context.Context, w io.Writer, groups []*output.PolicyGroup) error {
	table, write := createTable(policyGroupHeaders)

	for _, pg := range groups {
		// Humanized
		age := d.timeHumanizer.TimeToHuman(pg.CreatedAt)
		table.AddRow(pg.Name, len(pg.Users), len(pg.Specs), age)
	}

	return write(w)
}
