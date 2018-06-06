package formatter

import (
	"fmt"
	"strconv"
	"strings"

	units "github.com/docker/go-units"
	"github.com/storageos/go-api/types"
)

const (
	defaultPoolQuietFormat = "{{.Name}}"
	defaultPoolTableFormat = "table {{.Name}}\t{{.Default}}\t{{.NodeSelector}}\t{{.DeviceSelector}}\t{{.Nodes}}\t{{.Total}}\t{{.CapacityUsed}}"

	poolNameHeader           = "NAME"
	poolDefaultHeader        = "DEFAULT"
	poolNodeSelectorHeader   = "NODE_SELECTOR"
	poolDeviceSelectorHeader = "DEVICE_SELECTOR"
	poolNodesHeader          = "NODES"
	poolCapacityUsedHeader   = "USED"
	poolTotalHeader          = "TOTAL"
)

// NewPoolFormat returns a format for use with a pool Context
func NewPoolFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultPoolQuietFormat
		}
		return defaultPoolTableFormat
	case RawFormatKey:
		if quiet {
			return `name: {{.Name}}`
		}
		return `name: {{.Name}}\n node selector: {{.NodeSelector}}\n`
	}
	return Format(source)
}

// PoolWrite writes formatted pools using the Context
func PoolWrite(ctx Context, pools []*types.Pool) error {
	render := func(format func(subContext subContext) error) error {
		for _, pool := range pools {
			if err := format(&poolContext{v: *pool}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&poolContext{}, render)
}

type poolContext struct {
	HeaderContext
	v types.Pool
}

func (c *poolContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *poolContext) Name() string {
	c.AddHeader(poolNameHeader)
	return c.v.Name
}

func (c *poolContext) Default() string {
	c.AddHeader(poolDefaultHeader)
	return strconv.FormatBool(c.v.Default)
}

func (c *poolContext) NodeSelector() string {
	c.AddHeader(poolNodeSelectorHeader)
	return c.v.NodeSelector
}

func (c *poolContext) DeviceSelector() string {
	c.AddHeader(poolDeviceSelectorHeader)
	return c.v.DeviceSelector
}

func (c *poolContext) Nodes() string {
	c.AddHeader(poolNodesHeader)
	return strconv.Itoa(len(c.v.Nodes))
}

func (c *poolContext) CapacityUsed() string {
	c.AddHeader(poolCapacityUsedHeader)
	if c.v.CapacityStats.TotalCapacityBytes == 0 {
		return "-"
	}
	return fmt.Sprintf("%.2f%%", float64(c.v.CapacityStats.TotalCapacityBytes-c.v.CapacityStats.AvailableCapacityBytes)*100/float64(c.v.CapacityStats.TotalCapacityBytes))
}

func (c *poolContext) Total() string {
	c.AddHeader(poolTotalHeader)
	return units.HumanSize(float64(c.v.CapacityStats.TotalCapacityBytes))
}

func (c *poolContext) Labels() string {
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

func (c *poolContext) Label(name string) string {

	n := strings.Split(name, ".")
	r := strings.NewReplacer("-", " ", "_", " ")
	h := r.Replace(n[len(n)-1])

	c.AddHeader(h)

	if c.v.Labels == nil {
		return ""
	}
	return c.v.Labels[name]
}
