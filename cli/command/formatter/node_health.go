package formatter

import (
	"github.com/storageos/go-api/types"
)

const (
	defaultNodeSubmodulesTableFormat = "table {{.Name}}\t{{.Status}}\t{{.UpdatedAt}}\t{{.ChangedAt}}\t{{.Message}}"

	nodeSubmodulesNameHeader      = "SUBMODULE"
	nodeSubmodulesStatusHeader    = "STATUS"
	nodeSubmodulesUpdatedAtHeader = "UPDATED_AT"
	nodeSubmodulesChangedAtHeader = "CHANGED_AT"
	nodeSubmodulesMessageHeader   = "MESSAGE"
)

// NewNodeHealthFormat returns a format for use with a node health Context
func NewNodeHealthFormat(source string) Format {
	switch source {
	case TableFormatKey:
		return defaultNodeSubmodulesTableFormat
	case RawFormatKey:
		return `submodule: {{.Submodule}}\nstatus: {{.Status}}\nupdated_at: {{.UpdatedAt}}\nchanged_at: {{.ChangedAt}}\nmessage: {{.Message}}\n`
	}
	return Format(source)
}

// NodeHealthWrite writes formatted NamedSubModuleStatus elements using the Context
func NodeHealthWrite(ctx Context, nodesHealth []types.NamedSubModuleStatus) error {
	render := func(format func(subContext subContext) error) error {
		for _, status := range nodesHealth {
			if err := format(&nodeHealthContext{v: status}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&nodeHealthContext{}, render)
}

type nodeHealthContext struct {
	HeaderContext
	v types.NamedSubModuleStatus
}

func (n *nodeHealthContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(n)
}

func (n *nodeHealthContext) Name() string {
	n.AddHeader(nodeSubmodulesNameHeader)
	return n.v.Name
}

func (n *nodeHealthContext) Status() string {
	n.AddHeader(nodeSubmodulesStatusHeader)
	return n.v.Status
}

func (n *nodeHealthContext) UpdatedAt() string {
	n.AddHeader(nodeSubmodulesUpdatedAtHeader)
	return n.v.UpdatedAt
}

func (n *nodeHealthContext) ChangedAt() string {
	n.AddHeader(nodeSubmodulesChangedAtHeader)
	return n.v.ChangedAt
}

func (n *nodeHealthContext) Message() string {
	n.AddHeader(nodeSubmodulesMessageHeader)
	return n.v.Message
}
