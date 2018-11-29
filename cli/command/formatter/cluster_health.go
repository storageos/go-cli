package formatter

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/storageos/go-cli/types"
)

const (
	defaultClusterHealthQuietFormat  = "{{.Status}}"
	defaultClusterHealthTableFormat  = "table {{.Node}}\t{{.Address}}\t{{.CPStatus}}\t{{.DPStatus}}"
	detailedClusterHealthTableFormat = "table {{.Node}}\t{{.Address}}\t{{.Status}}\t{{.KV}}\t{{.NATS}}\t{{.DirectFSClient}}\t{{.Director}}\t{{.FSDriver}}\t{{.FS}}"
	cpClusterHealthTableFormat       = "table {{.Node}}\t{{.Address}}\t{{.Status}}\t{{.KV}}\t{{.KVWrite}}\t{{.NATS}}"
	dpClusterHealthTableFormat       = "table {{.Node}}\t{{.Address}}\t{{.Status}}\t{{.DirectFSClient}}\t{{.Director}}\t{{.FSDriver}}\t{{.FS}}"

	clusterHealthNodeHeader           = "NODE"
	clusterHealthAddressHeader        = "ADDRESS"
	clusterHealthCPStatusHeader       = "CP_STATUS"
	clusterHealthDPStatusHeader       = "DP_STATUS"
	clusterHealthStatusHeader         = "STATUS"
	clusterHealthNATSHeader           = "NATS"
	clusterHealthKVHeader             = "KV"
	clusterHealthKVWriteHeader        = "KV_WRITE"
	clusterHealthDirectFSClientHeader = "DFS_CLIENT"
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
	case DetailedTableFormatKey:
		return detailedClusterHealthTableFormat
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
		return `node: {{.Node}}\nstatus: {{.Status}}\nkv: {{.KV}}\nkv_write: {{.KVWrite}}\nnats: {{.NATS}}\ndfs_client: {{.DirectFSClient}}\ndirector: {{.Director}}\nfs_driver: {{.FSDriver}}\nfs: {{.FS}}\n`
	}
	return Format(source)
}

// ClusterHealthWrite writes formatted ClusterHealthhelements using the Context
func ClusterHealthWrite(ctx Context, nodes []*types.Node) error {
	// Try handle a custom format, excluding the predefined templates
	TryFormatUnless(
		string(ctx.Format),
		nodes,
		defaultClusterHealthQuietFormat,
		defaultClusterHealthTableFormat,
		detailedClusterHealthTableFormat,
		cpClusterHealthTableFormat,
		dpClusterHealthTableFormat,
		`{{.CPStatus}}`,
		`{{.DPStatus}}`,
		`{{.Node}}: {{.Status}}`,
		`node: {{.Node}}\nstatus: {{.Status}}\nkv: {{.KV}}\nkv_write: {{.KVWrite}}\nnats: {{.NATS}}\ndfs_client: {{.DirectFSClient}}\ndirector: {{.Director}}\nfs_driver: {{.FSDriver}}\nfs: {{.FS}}\n`,
	)

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
	if err == nil && u.Host != "" {
		return u.Host
	}
	return c.v.AdvertiseAddress
}

func (c *clusterHealthContext) healthy() bool {
	return c.cpHealthy() && c.dpHealthy()
}

func (c *clusterHealthContext) cpHealthy() bool {
	return c.v.Health.CP.NATS.Status+
		c.v.Health.CP.KV.Status+
		c.v.Health.CP.KVWrite.Status == "alivealivealive"
}

func (c *clusterHealthContext) cpDegraded() string {
	degraded := []string{}
	if c.v.Health.CP.KV.Status != "alive" {
		degraded = append(degraded, clusterHealthKVHeader)
	}
	if c.v.Health.CP.KVWrite.Status != "alive" {
		degraded = append(degraded, clusterHealthKVWriteHeader)
	}
	if c.v.Health.CP.NATS.Status != "alive" {
		degraded = append(degraded, clusterHealthNATSHeader)
	}
	return "Degraded (" + strings.Join(degraded, ", ") + ")"
}

func (c *clusterHealthContext) dpHealthy() bool {
	return c.v.Health.DP.DirectFSClient.Status+
		c.v.Health.DP.Director.Status+
		c.v.Health.DP.FSDriver.Status+
		c.v.Health.DP.FS.Status == "alivealivealivealive"
}

func (c *clusterHealthContext) dpDegraded() string {
	degraded := []string{}
	if c.v.Health.DP.DirectFSClient.Status != "alive" {
		degraded = append(degraded, clusterHealthDirectFSClientHeader)
	}
	if c.v.Health.DP.Director.Status != "alive" {
		degraded = append(degraded, clusterHealthDirectorHeader)
	}
	if c.v.Health.DP.FSDriver.Status != "alive" {
		degraded = append(degraded, clusterHealthFSDriverHeader)
	}
	if c.v.Health.DP.FS.Status != "alive" {
		degraded = append(degraded, clusterHealthFSHeader)
	}
	return "Degraded (" + strings.Join(degraded, ", ") + ")"
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
	c.AddHeader(clusterHealthCPStatusHeader)
	if c.v.Health.CP == nil {
		return "Unreachable"
	}
	if c.cpHealthy() {
		return "Healthy"
	}
	return c.cpDegraded()
}

func (c *clusterHealthContext) DPStatus() string {
	c.AddHeader(clusterHealthDPStatusHeader)
	if c.v.Health.DP == nil {
		return "Unreachable"
	}
	if c.dpHealthy() {
		return "Healthy"
	}
	return c.dpDegraded()
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

func (c *clusterHealthContext) DirectFSClient() string {
	c.AddHeader(clusterHealthDirectFSClientHeader)
	if c.v.Health.DP == nil {
		return clusterHealthUnknown
	}
	return c.v.Health.DP.DirectFSClient.Status
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
