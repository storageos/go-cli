package formatter

import (
	"fmt"
	"net/url"

	"github.com/storageos/go-cli/types"
)

const (
	defaultClusterHealthQuietFormat = "{{.Status}}"
	defaultClusterHealthTableFormat = "table {{.Node}}\t{{.Address}}\t{{.Status}}\t{{.KV}}\t{{.NATS}}\t{{.Scheduler}}\t{{.DirectFSClient}}\t{{.DirectFSServer}}\t{{.Director}}\t{{.FSDriver}}\t{{.FS}}"
	cpClusterHealthTableFormat      = "table {{.Node}}\t{{.Address}}\t{{.Status}}\t{{.KV}}\t{{.KVWrite}}\t{{.NATS}}\t{{.Scheduler}}"
	dpClusterHealthTableFormat      = "table {{.Node}}\t{{.Address}}\t{{.Status}}\t{{.DirectFSClient}}\t{{.DirectFSServer}}\t{{.Director}}\t{{.FSDriver}}\t{{.FS}}"

	clusterHealthNodeHeader           = "NODE"
	clusterHealthAddressHeader        = "ADDRESS"
	clusterHealthStatusHeader         = "STATUS"
	clusterHealthNATSHeader           = "NATS"
	clusterHealthKVHeader             = "KV"
	clusterHealthKVWriteHeader        = "KV_WRITE"
	clusterHealthSchedulerHeader      = "SCHEDULER"
	clusterHealthDirectFSClientHeader = "DFS_CLIENT"
	clusterHealthDirectFSServerHeader = "DFS_SERVER"
	clusterHealthDirectorHeader       = "DIRECTOR"
	clusterHealthFSDriverHeader       = "FS_DRIVER"
	clusterHealthFSHeader             = "FS"

	clusterHealthUnknown = "unknown"
)

// NewClusterHealthFormat returns a format for use with a cluster health Context
func NewClusterHealthFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultClusterHealthQuietFormat
		}
		return defaultClusterHealthTableFormat
	case CPHealthTableFormatKey:
		if quiet {
			return `{{.CPStatus}}`
		}
		return cpClusterHealthTableFormat
	case DPHealthTableFormatKey:
		if quiet {
			return `{{.DPStatus}}`
		}
		return dpClusterHealthTableFormat
	case RawFormatKey:
		if quiet {
			return `{{.Node}}: {{.Status}}`
		}
		return `node: {{.Node}}\nstatus: {{.Status}}\nkv: {{.KV}}\nkv_write: {{.KVWrite}}\nnats: {{.NATS}}\nscheduler: {{.Scheduler}}\ndfs_client: {{.DirectFSClient}}\ndfs_server: {{.DirectFSServer}}\ndirector: {{.Director}}\nfs_driver: {{.FSDriver}}\nfs: {{.FS}}\n`
	}
	return Format(source)
}

// ClusterHealthWrite writes formatted ClusterHealthhelements using the Context
func ClusterHealthWrite(ctx Context, nodes []*types.Node) error {
	if len(nodes) == 0 {
		return fmt.Errorf("No cluster nodes found")
	}
	render := func(format func(subContext subContext) error) error {
		for _, node := range nodes {
			if err := format(&clusterHealthContext{v: node}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&clusterHealthContext{}, render)
}

type clusterHealthContext struct {
	HeaderContext
	v *types.Node
}

func (c *clusterHealthContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *clusterHealthContext) Node() string {
	c.AddHeader(clusterHealthNodeHeader)
	return c.v.Name
}

func (c *clusterHealthContext) Address() string {
	c.AddHeader(clusterHealthAddressHeader)
	u, err := url.Parse(c.v.AdvertiseAddress)
	if err != nil {
		return c.v.AdvertiseAddress
	}
	return u.Hostname()
}

func (c *clusterHealthContext) healthy() bool {
	return c.cpHealthy() && c.dpHealthy()
}

func (c *clusterHealthContext) cpHealthy() bool {
	return c.v.Health.CP.NATS.Status+
		c.v.Health.CP.KV.Status+
		c.v.Health.CP.KVWrite.Status+
		c.v.Health.CP.Scheduler.Status == "alivealivealivealive"
}

func (c *clusterHealthContext) dpHealthy() bool {
	return c.v.Health.DP.DirectFSClient.Status+
		c.v.Health.DP.DirectFSServer.Status+
		c.v.Health.DP.Director.Status+
		c.v.Health.DP.FSDriver.Status+
		c.v.Health.DP.FS.Status == "alivealivealivealivealive"
}

func (c *clusterHealthContext) Status() string {
	c.AddHeader(clusterHealthStatusHeader)
	if c.v.Health.CP == nil || c.v.Health.DP == nil {
		return "Unreachable"
	}
	if c.healthy() {
		return "Healthy"
	}
	return "Not Ready"
}

func (c *clusterHealthContext) CPStatus() string {
	c.AddHeader(clusterHealthStatusHeader)
	if c.v.Health.CP == nil {
		return "Unreachable"
	}
	if c.cpHealthy() {
		return "Healthy"
	}
	return "Not Ready"
}

func (c *clusterHealthContext) DPStatus() string {
	c.AddHeader(clusterHealthStatusHeader)
	if c.v.Health.DP == nil {
		return "Unreachable"
	}
	if c.dpHealthy() {
		return "Healthy"
	}
	return "Not Ready"
}

func (c *clusterHealthContext) NATS() string {
	c.AddHeader(clusterHealthNATSHeader)
	if c.v.Health.CP == nil {
		return clusterHealthUnknown
	}
	return c.v.Health.CP.NATS.Status
}

func (c *clusterHealthContext) KV() string {
	c.AddHeader(clusterHealthKVHeader)
	if c.v.Health.CP == nil {
		return clusterHealthUnknown
	}
	return c.v.Health.CP.KV.Status
}

func (c *clusterHealthContext) KVWrite() string {
	c.AddHeader(clusterHealthKVWriteHeader)
	if c.v.Health.CP == nil {
		return clusterHealthUnknown
	}
	return c.v.Health.CP.KVWrite.Status
}

func (c *clusterHealthContext) Scheduler() string {
	c.AddHeader(clusterHealthSchedulerHeader)
	if c.v.Health.CP == nil {
		return clusterHealthUnknown
	}
	return c.v.Health.CP.Scheduler.Status
}

func (c *clusterHealthContext) DirectFSClient() string {
	c.AddHeader(clusterHealthDirectFSClientHeader)
	if c.v.Health.DP == nil {
		return clusterHealthUnknown
	}
	return c.v.Health.DP.DirectFSClient.Status
}

func (c *clusterHealthContext) DirectFSServer() string {
	c.AddHeader(clusterHealthDirectFSServerHeader)
	if c.v.Health.DP == nil {
		return clusterHealthUnknown
	}

	return c.v.Health.DP.DirectFSServer.Status
}

func (c *clusterHealthContext) Director() string {
	c.AddHeader(clusterHealthDirectorHeader)
	if c.v.Health.DP == nil {
		return clusterHealthUnknown
	}
	return c.v.Health.DP.Director.Status
}

func (c *clusterHealthContext) FSDriver() string {
	c.AddHeader(clusterHealthFSDriverHeader)
	if c.v.Health.DP == nil {
		return clusterHealthUnknown
	}
	return c.v.Health.DP.FSDriver.Status
}

func (c *clusterHealthContext) FS() string {
	c.AddHeader(clusterHealthFSHeader)
	if c.v.Health.DP == nil {
		return clusterHealthUnknown
	}
	return c.v.Health.DP.FS.Status
}
