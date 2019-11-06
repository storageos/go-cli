package output

import (
	"encoding/json"
	"io"

	"code.storageos.net/storageos/c2-cli/pkg/node"
)

type JSONDisplayer struct {
	encoderIndent string
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

func NewJSONDisplayer(encoderIndent string) *JSONDisplayer {
	return &JSONDisplayer{
		encoderIndent: encoderIndent,
	}
}
