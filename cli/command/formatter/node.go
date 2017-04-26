package formatter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/storageos/go-api/types"
)

const (
	defaultNodeQuietFormat = "{{.Name}}"
	defaultNodeTableFormat = "table {{.Name}}\t{{.Address}}\t{{.Health}}\t{{.Scheduler}}\t{{.Labels}}"

	nodeNameHeader      = "NAME"
	nodeAddressHeader   = "ADDRESS"
	nodeHealthHeader    = "HEALTH"
	nodeSchedulerHeader = "SCHEDULER"
	nodeLabelHeader     = "LABEL"
)

// NewNodeFormat returns a format for use with a node Context
func NewNodeFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultNodeQuietFormat
		}
		return defaultNodeTableFormat
	case RawFormatKey:
		if quiet {
			return `name: {{.Name}}`
		}
		return `name: {{.Name}}\naddress: {{.Address}}\nhealth: {{.Health}}\nscheduler: {{.Scheduler}}\nlabels: {{.Labels}}\n`
	}
	return Format(source)
}

// NodeWrite writes formatted nodes using the Context
func NodeWrite(ctx Context, nodes []*types.Controller) error {
	render := func(format func(subContext subContext) error) error {
		for _, node := range nodes {
			if err := format(&nodeContext{v: *node}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&nodeContext{}, render)
}

type nodeContext struct {
	HeaderContext
	v types.Controller
}

func (c *nodeContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *nodeContext) Name() string {
	c.AddHeader(nodeNameHeader)
	return c.v.Name
}

func (c *nodeContext) Address() string {
	c.AddHeader(nodeAddressHeader)
	return c.v.Address
}

func (c *nodeContext) Health() string {
	c.AddHeader(nodeHealthHeader)
	return c.v.Health
}
func (c *nodeContext) Scheduler() string {
	c.AddHeader(nodeSchedulerHeader)
	return strconv.FormatBool(c.v.Scheduler)
}

func (c *nodeContext) Labels() string {
	c.AddHeader(labelsHeader)
	if c.v.Labels == nil {
		return ""
	}

	var joinLabels []string
	for k, v := range c.v.Labels {
		joinLabels = append(joinLabels, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(joinLabels, ",")
}

func (c *nodeContext) Label(name string) string {

	n := strings.Split(name, ".")
	r := strings.NewReplacer("-", " ", "_", " ")
	h := r.Replace(n[len(n)-1])

	c.AddHeader(h)

	if c.v.Labels == nil {
		return ""
	}
	return c.v.Labels[name]
}
