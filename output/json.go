package output

import (
	"encoding/json"
	"io"

	"code.storageos.net/storageos/c2-cli/pkg/cluster"
	"code.storageos.net/storageos/c2-cli/pkg/node"
	"code.storageos.net/storageos/c2-cli/pkg/volume"
)

type JSONDisplayer struct {
	encoderIndent string
}

func (d *JSONDisplayer) WriteGetCluster(w io.Writer, resource *cluster.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resource)
}

func (d *JSONDisplayer) WriteGetNode(w io.Writer, resource *node.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resource)
}

func (d *JSONDisplayer) WriteGetNodeList(w io.Writer, resources []*node.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resources)
}

func (d *JSONDisplayer) WriteGetVolume(w io.Writer, resource *volume.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resource)
}

func (d *JSONDisplayer) WriteGetVolumeList(w io.Writer, resources []*volume.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resources)
}

func (d *JSONDisplayer) WriteDescribeCluster(w io.Writer, resource *cluster.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resource)
}

func (d *JSONDisplayer) WriteDescribeNode(w io.Writer, resource *node.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resource)
}

func (d *JSONDisplayer) WriteDescribeNodeList(w io.Writer, resources []*node.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resources)
}

func (d *JSONDisplayer) WriteDescribeVolume(w io.Writer, resource *volume.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resource)
}

func NewJSONDisplayer(encoderIndent string) *JSONDisplayer {
	return &JSONDisplayer{
		encoderIndent: encoderIndent,
	}
}
