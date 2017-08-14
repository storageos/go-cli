package formatter

import (
	"fmt"

	"github.com/storageos/go-cli/types"
)

const (
	defaultHealthCPTableFormat = "table {{.Node}}\t{{.Status}}\t{{.NATS}}\t{{.KV}}\t{{.KVWrite}}\t{{.Scheduler}}"

	clusterHealthCPNodeHeader      = "NODE"
	clusterHealthCPStatusHeader    = "STATUS"
	clusterHealthCPNATSHeader      = "NATS"
	clusterHealthCPKVHeader        = "KV"
	clusterHealthCPKVWriteHeader   = "KV_WRITE"
	clusterHealthCPSchedulerHeader = "SCHEDULER"
)

// NewHealthCPFormat returns a format for use with a cpHealth Context
func NewHealthCPFormat(source string) Format {
	switch source {
	case TableFormatKey:
		return defaultHealthCPTableFormat
	case RawFormatKey:
		return `node: {{.Node}}\nstatus: {{.Status}}\nnats: {{.NATS}}\nkv: {{.KV}}\nkv_write: {{.KVWrite}}\nscheduler: {{.Scheduler}}\n`
	}
	return Format(source)
}

// ClusterHealthCPWrite writes formatted ClusterHealthCP elements using the Context
func ClusterHealthCPWrite(ctx Context, nodes []*types.Node) error {
	if len(nodes) == 0 {
		return fmt.Errorf("No cluster nodes found")
	}
	render := func(format func(subContext subContext) error) error {
		for _, node := range nodes {
			if err := format(&cpHealthContext{v: node}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&cpHealthContext{}, render)
}

type cpHealthContext struct {
	HeaderContext
	v *types.Node
}

func (c *cpHealthContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *cpHealthContext) Node() string {
	c.AddHeader(clusterHealthCPNodeHeader)
	return c.v.Name
}

func (c *cpHealthContext) healthy() bool {
	return c.v.Health.CP.NATS.Status+
		c.v.Health.CP.KV.Status+
		c.v.Health.CP.KVWrite.Status+
		c.v.Health.CP.Scheduler.Status == "alivealivealivealive"
}

func (c *cpHealthContext) Status() string {
	c.AddHeader(clusterHealthCPStatusHeader)
	if c.healthy() {
		return "Healthy"
	}
	return "Not Ready"
}

func (c *cpHealthContext) NATS() string {
	c.AddHeader(clusterHealthCPNATSHeader)
	return c.v.Health.CP.NATS.Status
}

func (c *cpHealthContext) KV() string {
	c.AddHeader(clusterHealthCPKVHeader)
	return c.v.Health.CP.KV.Status
}

func (c *cpHealthContext) KVWrite() string {
	c.AddHeader(clusterHealthCPKVWriteHeader)
	return c.v.Health.CP.KVWrite.Status
}

func (c *cpHealthContext) Scheduler() string {
	c.AddHeader(clusterHealthCPSchedulerHeader)
	return c.v.Health.CP.Scheduler.Status
}
