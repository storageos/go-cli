package formatter

import (
	"github.com/storageos/go-api/types"
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
func ClusterHealthCPWrite(ctx Context, clusterHealth *types.ClusterHealthCP) error {
	render := func(format func(subContext subContext) error) error {
		for _, status := range *clusterHealth {
			if err := format(&cpHealthContext{v: status}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&cpHealthContext{}, render)
}

type cpHealthContext struct {
	HeaderContext
	v types.CPHealthStatusWithID
}

func (c *cpHealthContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *cpHealthContext) Node() string {
	c.AddHeader(clusterHealthCPNodeHeader)
	return c.v.ID
}

func (c *cpHealthContext) healthy() bool {
	return c.v.NATS.Status+c.v.KV.Status+c.v.KVWrite.Status+c.v.Scheduler.Status == "alivealivealivealive"
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
	return c.v.NATS.Status
}

func (c *cpHealthContext) KV() string {
	c.AddHeader(clusterHealthCPKVHeader)
	return c.v.KV.Status
}

func (c *cpHealthContext) KVWrite() string {
	c.AddHeader(clusterHealthCPKVWriteHeader)
	return c.v.KVWrite.Status
}

func (c *cpHealthContext) Scheduler() string {
	c.AddHeader(clusterHealthCPSchedulerHeader)
	return c.v.Scheduler.Status
}
