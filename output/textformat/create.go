package textformat

import (
	"context"
	"fmt"
	"io"

	"code.storageos.net/storageos/c2-cli/output"
)

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

// CreateVolumeAsync writes a successful request submission string to w.
func (d *Displayer) CreateVolumeAsync(ctx context.Context, w io.Writer) error {
	_, err := fmt.Fprintln(w, "volume provision request accepted")
	return err
}

// CreateNamespace builds a human friendly representation of resource, writing
// the result to w.
func (d *Displayer) CreateNamespace(ctx context.Context, w io.Writer, namespace *output.Namespace) error {
	table, write := createTable(namespaceHeaders)

	// Humanized
	age := d.timeHumanizer.TimeToHuman(namespace.CreatedAt)

	table.AddRow(namespace.Name, age)

	return write(w)
}

// CreatePolicyGroup builds a human friendly representation of resource, writing
// the result to w.
func (d *Displayer) CreatePolicyGroup(ctx context.Context, w io.Writer, group *output.PolicyGroup) error {
	table, write := createTable(policyGroupHeaders)

	// Humanized
	age := d.timeHumanizer.TimeToHuman(group.CreatedAt)
	table.AddRow(group.Name, len(group.Users), len(group.Specs), age)

	return write(w)
}
