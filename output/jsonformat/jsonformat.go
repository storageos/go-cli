// Package jsonformat implements a JSON format output mechanism for StorageOS
// API resources.
package jsonformat

import (
	"context"
	"encoding/json"
	"io"

	"code.storageos.net/storageos/c2-cli/output"
)

// DefaultEncodingIndent is the encoding indent string which consumers of the
// output package can default to when initialising Displayer types.
const DefaultEncodingIndent = "\t"

// Displayer is a type which encodes StorageOS resources to JSON and writes the
// result to io.Writers.
type Displayer struct {
	encoderIndent string
}

func (d *Displayer) encode(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(v)
}

// -----------------------------------------------------------------------------
// CREATE
// -----------------------------------------------------------------------------

// CreateUser encodes user as JSON, writing the result to w.
func (d *Displayer) CreateUser(ctx context.Context, w io.Writer, user *output.User) error {
	return d.encode(w, user)
}

// CreateVolume encodes volume as JSON, writing the result to w.
func (d *Displayer) CreateVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	return d.encode(w, volume)
}

// CreateVolumeAsync writes nothing to w.
func (d *Displayer) CreateVolumeAsync(ctx context.Context, w io.Writer) error {
	return nil
}

// CreateNamespace encodes namespace as JSON, writing the result to w.
func (d *Displayer) CreateNamespace(ctx context.Context, w io.Writer, namespace *output.Namespace) error {
	return d.encode(w, namespace)
}

// CreatePolicyGroup encodes group as JSON, writing the result to w.
func (d *Displayer) CreatePolicyGroup(ctx context.Context, w io.Writer, group *output.PolicyGroup) error {
	return d.encode(w, group)
}

// -----------------------------------------------------------------------------
// UPDATE
// -----------------------------------------------------------------------------

// UpdateLicence encodes licence as JSON, writing the result to w.
func (d *Displayer) UpdateLicence(ctx context.Context, w io.Writer, licence *output.Licence) error {
	return d.encode(w, licence)
}

// UpdateVolume encodes licence as JSON, writing the result to w.
func (d *Displayer) UpdateVolume(ctx context.Context, w io.Writer, updatedVol output.VolumeUpdate) error {
	return d.encode(w, updatedVol)
}

// SetReplicas does not encode anything to w
func (d *Displayer) SetReplicas(ctx context.Context, w io.Writer, new uint64) error {
	return nil
}

// -----------------------------------------------------------------------------
// GET
// -----------------------------------------------------------------------------

// GetCluster encodes resource as JSON, writing the result to w.
func (d *Displayer) GetCluster(ctx context.Context, w io.Writer, resource *output.Cluster) error {
	return d.encode(w, resource)
}

// GetLicence encodes resource as JSON, writing the result to w.
func (d *Displayer) GetLicence(ctx context.Context, w io.Writer, resource *output.Licence) error {
	return d.encode(w, resource)
}

// GetDiagnostics encodes outputPath as JSON, writing the result to w.
func (d *Displayer) GetDiagnostics(ctx context.Context, w io.Writer, outputPath string) error {
	o := struct {
		OutputPath string `json:"outputPath"`
	}{
		OutputPath: outputPath,
	}

	return d.encode(w, o)
}

// GetUser encodes resources as JSON, writing the result to w.
func (d *Displayer) GetUser(ctx context.Context, w io.Writer, user *output.User) error {
	return d.encode(w, user)
}

// GetUsers encodes resources as JSON, writing the result to w.
func (d *Displayer) GetUsers(ctx context.Context, w io.Writer, users []*output.User) error {
	return d.encode(w, users)
}

// GetNode encodes resource as JSON, writing the result to w.
func (d *Displayer) GetNode(ctx context.Context, w io.Writer, resource *output.Node) error {
	return d.encode(w, resource)
}

// GetListNodes encodes resources as JSON, writing the result to w.
func (d *Displayer) GetListNodes(ctx context.Context, w io.Writer, resources []*output.Node) error {
	return d.encode(w, resources)
}

// GetNamespace encodes resource as JSON, writing the result to w.
func (d *Displayer) GetNamespace(ctx context.Context, w io.Writer, resource *output.Namespace) error {
	return d.encode(w, resource)
}

// GetListNamespaces encodes resources as JSON, writing the result to w.
func (d *Displayer) GetListNamespaces(ctx context.Context, w io.Writer, resources []*output.Namespace) error {
	return d.encode(w, resources)
}

// GetVolume encodes resource as JSON, writing the result to w.
func (d *Displayer) GetVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	return d.encode(w, volume)
}

// GetListVolumes encodes resources as JSON, writing the result to w.
func (d *Displayer) GetListVolumes(ctx context.Context, w io.Writer, volumes []*output.Volume) error {
	return d.encode(w, volumes)
}

// GetPolicyGroup encodes resource as JSON, writing the result to w.
func (d *Displayer) GetPolicyGroup(ctx context.Context, w io.Writer, group *output.PolicyGroup) error {
	return d.encode(w, group)
}

// GetListPolicyGroups encodes resources as JSON, writing the result to w.
func (d *Displayer) GetListPolicyGroups(ctx context.Context, w io.Writer, groups []*output.PolicyGroup) error {
	return d.encode(w, groups)
}

// -----------------------------------------------------------------------------
// DESCRIBE
// -----------------------------------------------------------------------------

// DescribeCluster encodes a cluster as JSON, writing the result to w.
func (d *Displayer) DescribeCluster(ctx context.Context, w io.Writer, c *output.Cluster) error {
	return d.encode(w, c)
}

// DescribeLicence encodes a licence as JSON, writing the result to w.
func (d *Displayer) DescribeLicence(ctx context.Context, w io.Writer, l *output.Licence) error {
	return d.encode(w, l)
}

// DescribeNamespace encodes a namespace as JSON, writing the result to w.
func (d *Displayer) DescribeNamespace(ctx context.Context, w io.Writer, namespace *output.Namespace) error {
	return d.encode(w, namespace)
}

// DescribeListNamespaces encodes a list of namespaces as JSON, writing the result to w.
func (d *Displayer) DescribeListNamespaces(ctx context.Context, w io.Writer, namespaces []*output.Namespace) error {
	return d.encode(w, namespaces)
}

// DescribeNode encodes node as JSON, writing the result to w.
func (d *Displayer) DescribeNode(ctx context.Context, w io.Writer, node *output.NodeDescription) error {
	return d.encode(w, node)
}

// DescribeListNodes encodes nodes as JSON, writing the result to w.
func (d *Displayer) DescribeListNodes(ctx context.Context, w io.Writer, nodes []*output.NodeDescription) error {
	return d.encode(w, nodes)
}

// DescribeVolume encodes volume as JSON, writing the result to w
func (d *Displayer) DescribeVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	return d.encode(w, volume)
}

// DescribeListVolumes encodes volumes as JSON, writing the result to w
func (d *Displayer) DescribeListVolumes(ctx context.Context, w io.Writer, volumes []*output.Volume) error {
	return d.encode(w, volumes)
}

// DescribePolicyGroup encodes policy group as JSON, writing the result to w
func (d *Displayer) DescribePolicyGroup(ctx context.Context, w io.Writer, group *output.PolicyGroup) error {
	return d.encode(w, group)
}

// DescribeListPolicyGroups encodes policy groups as JSON, writing the result to w
func (d *Displayer) DescribeListPolicyGroups(ctx context.Context, w io.Writer, groups []*output.PolicyGroup) error {
	return d.encode(w, groups)
}

// DescribeUser encodes user as JSON, writing the result to w
func (d *Displayer) DescribeUser(ctx context.Context, w io.Writer, user *output.User) error {
	return d.encode(w, user)
}

// DescribeListUsers encodes users as JSON, writing the result to w
func (d *Displayer) DescribeListUsers(ctx context.Context, w io.Writer, users []*output.User) error {
	return d.encode(w, users)
}

// -----------------------------------------------------------------------------
// DELETE
// -----------------------------------------------------------------------------

// DeleteUser encodes the user deletion confirmation as JSON, writing the
// results to w.
func (d *Displayer) DeleteUser(ctx context.Context, w io.Writer, confirmation output.UserDeletion) error {
	return d.encode(w, confirmation)
}

// DeleteVolume encodes the volume deletion confirmation as JSON, writing the
// result to w.
func (d *Displayer) DeleteVolume(ctx context.Context, w io.Writer, confirmation output.VolumeDeletion) error {
	return d.encode(w, confirmation)
}

// DeleteVolumeAsync writes nothing to w.
func (d *Displayer) DeleteVolumeAsync(ctx context.Context, w io.Writer, target output.VolumeDeletion) error {
	return nil
}

// DeleteNamespace encodes the namespace deletion confirmation as JSON, writing
// the result to w
func (d *Displayer) DeleteNamespace(ctx context.Context, w io.Writer, confirmation output.NamespaceDeletion) error {
	return d.encode(w, confirmation)
}

// DeletePolicyGroup encodes the policy group deletion confirmation as JSON, writing
// the result to w
func (d *Displayer) DeletePolicyGroup(ctx context.Context, w io.Writer, confirmation output.PolicyGroupDeletion) error {
	return d.encode(w, confirmation)
}

// DeleteNode encodes the node deletion confirmation as JSON, writing
// the result to w
func (d *Displayer) DeleteNode(ctx context.Context, w io.Writer, confirmation output.NodeDeletion) error {
	return d.encode(w, confirmation)
}

// DeleteNodeAsync writes nothing to w.
func (d *Displayer) DeleteNodeAsync(ctx context.Context, w io.Writer, target output.NodeDeletion) error {
	return nil
}

// -----------------------------------------------------------------------------
// OTHER
// -----------------------------------------------------------------------------

// AttachVolume writes nothing to the writer.
func (d *Displayer) AttachVolume(ctx context.Context, w io.Writer) error {
	return nil
}

// DetachVolume writes nothing to the writer.
func (d *Displayer) DetachVolume(ctx context.Context, w io.Writer) error {
	return nil
}

// NewDisplayer initialises a Displayer which encodes StorageOS resources as
// JSON, using encoderIndent as the indentation string.
func NewDisplayer(encoderIndent string) *Displayer {
	return &Displayer{
		encoderIndent: encoderIndent,
	}
}
