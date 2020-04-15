// Package yamlformat implements a YAML format output mechanism for StorageOS
// API resources.
package yamlformat

import (
	"context"
	"io"

	"gopkg.in/yaml.v3"

	"code.storageos.net/storageos/c2-cli/output"
)

// Displayer is a type which encodes StorageOS resources to YAML and writes the
// result to io.Writers.
type Displayer struct {
	encoderIndent string
}

func (d *Displayer) encode(w io.Writer, v interface{}) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(0)
	return enc.Encode(v)
}

// -----------------------------------------------------------------------------
// CREATE
// -----------------------------------------------------------------------------

// CreateUser encodes user as YAML, writing the result to w.
func (d *Displayer) CreateUser(ctx context.Context, w io.Writer, user *output.User) error {
	return d.encode(w, user)
}

// CreateVolume encodes resource as YAML, writing the result to w.
func (d *Displayer) CreateVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	return d.encode(w, volume)
}

// CreateVolumeAsync writes nothing to w.
func (d *Displayer) CreateVolumeAsync(ctx context.Context, w io.Writer) error {
	return nil
}

// CreateNamespace encodes namespace as YAML, writing the result to w.
func (d *Displayer) CreateNamespace(ctx context.Context, w io.Writer, namespace *output.Namespace) error {
	return d.encode(w, namespace)
}

// -----------------------------------------------------------------------------
// UPDATE
// -----------------------------------------------------------------------------

// UpdateLicence encodes licence as YAML, writing the result to w.
func (d *Displayer) UpdateLicence(ctx context.Context, w io.Writer, licence *output.Licence) error {
	return d.encode(w, licence)
}

// -----------------------------------------------------------------------------
// GET
// -----------------------------------------------------------------------------

// GetCluster encodes resource as YAML, writing the result to w.
func (d *Displayer) GetCluster(ctx context.Context, w io.Writer, resource *output.Cluster) error {
	return d.encode(w, resource)
}

// GetDiagnostics encodes outputPath as YAML, writing the result to w.
func (d *Displayer) GetDiagnostics(ctx context.Context, w io.Writer, outputPath string) error {
	output := struct {
		OutputPath string `yaml:"outputPath"`
	}{
		OutputPath: outputPath,
	}

	return d.encode(w, output)
}

// GetUser encodes resources as YAML, writing the result to w.
func (d *Displayer) GetUser(ctx context.Context, w io.Writer, user *output.User) error {
	return d.encode(w, user)
}

// GetUsers encodes resources as YAML, writing the result to w.
func (d *Displayer) GetUsers(ctx context.Context, w io.Writer, users []*output.User) error {
	return d.encode(w, users)
}

// GetNode encodes resource as YAML, writing the result to w.
func (d *Displayer) GetNode(ctx context.Context, w io.Writer, resource *output.Node) error {
	return d.encode(w, resource)
}

// GetListNodes encodes resources as YAML, writing the result to w.
func (d *Displayer) GetListNodes(ctx context.Context, w io.Writer, resources []*output.Node) error {
	return d.encode(w, resources)
}

// GetNamespace encodes resource as YAML, writing the result to w.
func (d *Displayer) GetNamespace(ctx context.Context, w io.Writer, resource *output.Namespace) error {
	return d.encode(w, resource)
}

// GetListNamespaces encodes resources as YAML, writing the result to w.
func (d *Displayer) GetListNamespaces(ctx context.Context, w io.Writer, resources []*output.Namespace) error {
	return d.encode(w, resources)
}

// GetVolume encodes resource as YAML, writing the result to w.
func (d *Displayer) GetVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	return d.encode(w, volume)
}

// GetListVolumes encodes resources as YAML, writing the result to w.
func (d *Displayer) GetListVolumes(ctx context.Context, w io.Writer, volumes []*output.Volume) error {
	return d.encode(w, volumes)
}

// GetPolicyGroup encodes resource as YAML, writing the result to w.
func (d *Displayer) GetPolicyGroup(ctx context.Context, w io.Writer, group *output.PolicyGroup) error {
	return d.encode(w, group)
}

// GetListPolicyGroups encodes resources as YAML, writing the result to w.
func (d *Displayer) GetListPolicyGroups(ctx context.Context, w io.Writer, groups []*output.PolicyGroup) error {
	return d.encode(w, groups)
}

// -----------------------------------------------------------------------------
// DESCRIBE
// -----------------------------------------------------------------------------

// DescribeCluster encodes a cluster as YAML, writing the result to w.
func (d *Displayer) DescribeCluster(ctx context.Context, w io.Writer, c *output.Cluster) error {
	return d.encode(w, c)
}

// DescribeNode encodes node as YAML, writing the result to w.
func (d *Displayer) DescribeNode(ctx context.Context, w io.Writer, node *output.NodeDescription) error {
	return d.encode(w, node)
}

// DescribeListNodes encodes nodes as YAML, writing the result to w.
func (d *Displayer) DescribeListNodes(ctx context.Context, w io.Writer, nodes []*output.NodeDescription) error {
	return d.encode(w, nodes)
}

// DescribeVolume encodes volume as YAML, writing the result to w
func (d *Displayer) DescribeVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	return d.encode(w, volume)
}

// DescribeListVolumes encodes volumes as YAML, writing the result to w
func (d *Displayer) DescribeListVolumes(ctx context.Context, w io.Writer, volumes []*output.Volume) error {
	return d.encode(w, volumes)
}

// DescribePolicyGroup encodes policy group as YAML, writing the result to w
func (d *Displayer) DescribePolicyGroup(ctx context.Context, w io.Writer, group *output.PolicyGroup) error {
	return d.encode(w, group)
}

// DescribeListPolicyGroups encodes policy groups as YAML, writing the result to w
func (d *Displayer) DescribeListPolicyGroups(ctx context.Context, w io.Writer, groups []*output.PolicyGroup) error {
	return d.encode(w, groups)
}

// DescribeUser encodes user as YAML, writing the result to w
func (d *Displayer) DescribeUser(ctx context.Context, w io.Writer, user *output.User) error {
	return d.encode(w, user)
}

// DescribeListUsers encodes users as YAML, writing the result to w
func (d *Displayer) DescribeListUsers(ctx context.Context, w io.Writer, users []*output.User) error {
	return d.encode(w, users)
}

// -----------------------------------------------------------------------------
// DELETE
// -----------------------------------------------------------------------------

// DeleteUser encodes the user deletion confirmation as YAML, writing the
// results to w.
func (d *Displayer) DeleteUser(ctx context.Context, w io.Writer, confirmation output.UserDeletion) error {
	return d.encode(w, confirmation)
}

// DeleteVolume encodes the deletion confirmation as YAML, writing the result
// to w.
func (d *Displayer) DeleteVolume(ctx context.Context, w io.Writer, confirmation output.VolumeDeletion) error {
	return d.encode(w, confirmation)
}

// DeleteVolumeAsync writes nothing to w.
func (d *Displayer) DeleteVolumeAsync(ctx context.Context, w io.Writer) error {
	return nil
}

// DeleteNamespace encodes the namespace deletion confirmation as YAML, writing
// the result to w
func (d *Displayer) DeleteNamespace(ctx context.Context, w io.Writer, confirmation output.NamespaceDeletion) error {
	return d.encode(w, confirmation)
}

// DeletePolicyGroup encodes the policy group deletion confirmation as YAML, writing
// the result to w
func (d *Displayer) DeletePolicyGroup(ctx context.Context, w io.Writer, confirmation output.PolicyGroupDeletion) error {
	return d.encode(w, confirmation)
}

// -----------------------------------------------------------------------------
// OTHER
// -----------------------------------------------------------------------------

// AttachVolume writes nothing in the writer
func (d *Displayer) AttachVolume(ctx context.Context, w io.Writer) error {
	return nil
}

// DetachVolume writes nothing to the writer
func (d *Displayer) DetachVolume(ctx context.Context, w io.Writer) error {
	return nil
}

// NewDisplayer initialises a Displayer which encodes StorageOS resources as
// YAML, using encoderIndent as the indentation string.
func NewDisplayer(encoderIndent string) *Displayer {
	return &Displayer{
		encoderIndent: encoderIndent,
	}
}
