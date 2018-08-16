package formatter

import (
	"time"

	units "github.com/docker/go-units"

	apiTypes "github.com/storageos/go-api/types"
	cliTypes "github.com/storageos/go-cli/types"
)

const (
	defaultNodeSubmodulesQuietFormat = "{{.Status}}"
	defaultNodeSubmodulesTableFormat = "table {{.Name}}\t{{.Type}}\t{{.Status}}\t{{.ChangedAt}}\t{{.UpdatedAt}}"
	cpNodeSubmodulesTableFormat      = "table {{if eq .T \"controlplane\"}}{{.Name}}\t{{.Type}}\t{{.Status}}\t{{.ChangedAt}}\t{{.UpdatedAt}}{{end}}"
	dpNodeSubmodulesTableFormat      = "table {{if eq .T \"dataplane\"}}{{.Name}}\t{{.Type}}\t{{.Status}}\t{{.ChangedAt}}\t{{.UpdatedAt}}{{end}}"
	cpNodeSubmodulesQuietFormat      = "{{if eq .T \"controlplane\"}}{{.Status}}{{end}}"
	dpNodeSubmodulesQuietFormat      = "{{if eq .T \"dataplane\"}}{{.Status}}{{end}}"
	defaultNodeRawFormat             = "submodule: {{.Name}}\ntype: {{.Type}}\nstatus: {{.Status}}\nupdated_at: {{.UpdatedAt}}\nchanged_at: {{.ChangedAt}}\n"

	nodeSubmodulesNameHeader      = "SUBMODULE"
	nodeSubmodulesTypeHeader      = "TYPE"
	nodeSubmodulesStatusHeader    = "STATUS"
	nodeSubmodulesUpdatedAtHeader = "UPDATED"
	nodeSubmodulesChangedAtHeader = "CHANGED"
	nodeSubmodulesMessageHeader   = "MESSAGE"
)

// NewNodeHealthFormat returns a format for use with a node health Context
func NewNodeHealthFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultClusterHealthQuietFormat
		}
		return defaultNodeSubmodulesTableFormat
	case CPHealthTableFormatKey:
		if quiet {
			return cpNodeSubmodulesQuietFormat
		}
		return cpNodeSubmodulesTableFormat
	case DPHealthTableFormatKey:
		if quiet {
			return dpNodeSubmodulesQuietFormat
		}
		return dpNodeSubmodulesTableFormat
	case RawFormatKey:
		return defaultNodeRawFormat
	}
	return Format(source)
}

// NodeHealthWrite writes formatted NamedSubModuleStatus elements using the Context
func NodeHealthWrite(ctx Context, node *cliTypes.Node) error {
	render := func(format func(subContext subContext) error) error {
		if node.Health.CP != nil {
			for _, submodule := range node.Health.CP.ToNamedSubmodules() {
				if err := format(&nodeHealthContext{t: "controlplane", v: submodule}); err != nil {
					return err
				}
			}
		}
		if node.Health.DP != nil {
			for _, submodule := range node.Health.DP.ToNamedSubmodules() {
				// Skip when there's not status data. This is to avoid printing
				// empty DP fields for now.
				if submodule.Status == "" {
					continue
				}
				if err := format(&nodeHealthContext{t: "dataplane", v: submodule}); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return ctx.Write(&nodeHealthContext{}, render)
}

type nodeHealthContext struct {
	HeaderContext
	t string
	v apiTypes.NamedSubModuleStatus
}

func (n *nodeHealthContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(n)
}

// T only returns the submodule type and avoids adding an output header
func (n *nodeHealthContext) T() string {
	return n.t
}

func (n *nodeHealthContext) Type() string {
	n.AddHeader(nodeSubmodulesTypeHeader)
	return n.t
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
	updatedTime, err := time.Parse(time.RFC3339, n.v.UpdatedAt)
	if err != nil {
		return n.v.UpdatedAt
	}
	return units.HumanDuration(time.Now().UTC().Sub(updatedTime)) + " ago"
}

func (n *nodeHealthContext) ChangedAt() string {
	n.AddHeader(nodeSubmodulesChangedAtHeader)
	changedTime, err := time.Parse(time.RFC3339, n.v.ChangedAt)
	if err != nil {
		return n.v.ChangedAt
	}
	return units.HumanDuration(time.Now().UTC().Sub(changedTime)) + " ago"
}

func (n *nodeHealthContext) Message() string {
	n.AddHeader(nodeSubmodulesMessageHeader)
	return n.v.Message
}
