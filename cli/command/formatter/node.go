package formatter

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/storageos/go-api/types"
)

const (
	defaultNodeQuietFormat = "{{.Name}}"
	defaultNodeTableFormat = "table {{.Name}}\t{{.Address}}\t{{.Health}}\t{{.Scheduler}}\t{{.Volumes}}\t{{.Capacity}}\t{{.CapacityUsed}}\t{{.Version}}\t{{.Labels}}"

	nodeNameHeader          = "NAME"
	nodeAddressHeader       = "ADDRESS"
	nodeHealthHeader        = "HEALTH"
	nodeSchedulerHeader     = "SCHEDULER"
	nodeVolumesHeader       = "VOLUMES"
	nodeTotalCapacityHeader = "TOTAL"
	nodeCapacityUsedHeader  = "USED"
	nodeVersionUsedHeader   = "VERSION"
	nodeLabelHeader         = "LABEL"
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
		return `name: {{.Name}}\naddress: {{.Address}}\nhealth: {{.Health}}\nscheduler: {{.Scheduler}}\nvolumes: {{.Volumes}}\ncapacity: {{.Capacity}}\ncapacityUsed: {{.CapacityUsed}}\nversion: {{.Version}}\nlabels: {{.Labels}}\n`
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

	if c.v.HealthUpdatedAt.IsZero() {
		return strings.Title(c.v.Health)
	}

	return fmt.Sprintf("%s %s", strings.Title(c.v.Health), units.HumanDuration(time.Since(c.v.HealthUpdatedAt)))
}

func (c *nodeContext) Scheduler() string {
	c.AddHeader(nodeSchedulerHeader)
	return strconv.FormatBool(c.v.Scheduler)
}

func (c *nodeContext) Volumes() string {
	c.AddHeader(nodeVolumesHeader)
	return fmt.Sprintf("M: %d, R: %d", c.v.VolumeStats.MasterVolumeCount, c.v.VolumeStats.ReplicaVolumeCount)
}

func (c *nodeContext) Capacity() string {
	c.AddHeader(nodeTotalCapacityHeader)
	if c.v.CapacityStats.TotalCapacityBytes == 0 {
		return "-"
	}

	return units.BytesSize(float64(c.v.CapacityStats.TotalCapacityBytes))
}

func (c *nodeContext) CapacityUsed() string {
	c.AddHeader(nodeCapacityUsedHeader)
	if c.v.CapacityStats.TotalCapacityBytes == 0 {
		return "-"
	}
	return fmt.Sprintf("%.2f%%", float64(c.v.CapacityStats.TotalCapacityBytes-c.v.CapacityStats.AvailableCapacityBytes)*100/float64(c.v.CapacityStats.TotalCapacityBytes))
}

func (c *nodeContext) Version() string {
	c.AddHeader(nodeVersionUsedHeader)
	return fmt.Sprintf("%s (%s rev)", c.v.VersionInfo["storageos"].Version, c.v.VersionInfo["storageos"].Revision)
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
