package formatter

import (
	"fmt"
	"strings"

	apiTypes "github.com/storageos/go-api/types"
)

const (
	defaultClusterHealthQuietFormat  = "{{.Status}}"
	defaultClusterHealthTableFormat  = "table {{.Node}}\t{{.CPStatus}}\t{{.DPStatus}}"
	detailedClusterHealthTableFormat = "table {{.Node}}\t{{.Status}}\t{{.KV}}\t{{.KVWrite}}\t{{.NATS}}\t{{.DirectFSInitiator}}\t{{.Director}}\t{{.Presentation}}\t{{.RDB}}"
	cpClusterHealthTableFormat       = "table {{.Node}}\t{{.Status}}\t{{.KV}}\t{{.KVWrite}}\t{{.NATS}}"
	dpClusterHealthTableFormat       = "table {{.Node}}\t{{.Status}}\t{{.DirectFSInitiator}}\t{{.Director}}\t{{.Presentation}}\t{{.RDB}}"

	clusterHealthNodeHeader              = "NODE"
	clusterHealthAddressHeader           = "ADDRESS"
	clusterHealthCPStatusHeader          = "CP_STATUS"
	clusterHealthDPStatusHeader          = "DP_STATUS"
	clusterHealthStatusHeader            = "STATUS"
	clusterHealthNATSHeader              = "NATS"
	clusterHealthKVHeader                = "KV"
	clusterHealthKVWriteHeader           = "KV_WRITE"
	clusterHealthDirectFSInitiatorHeader = "DFS_INITIATOR"
	clusterHealthDirectorHeader          = "DIRECTOR"
	clusterHealthPresentationHeader      = "PRESENTATION"
	clusterHealthRDBHeader               = "RDB"

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
		return `node: {{.Node}}\nstatus: {{.Status}}\nkv: {{.KV}}\nkv_write: {{.KVWrite}}\nnats: {{.NATS}}\ndfs_initiator: {{.DirectFSInitiator}}\ndirector: {{.Director}}\npresentation: {{.Presentation}}\nrdb: {{.RDB}}\n`
	}
	return Format(source)
}

// ClusterHealthWrite writes formatted ClusterHealthhelements using the Context
func ClusterHealthWrite(ctx Context, status []*apiTypes.ClusterHealthNode) error {
	TryFormatUnlessMatches(
		string(ctx.Format),
		status,
		TableMatcher,
		NewExactMatcher(defaultClusterHealthQuietFormat),
		NewExactMatcher(defaultClusterHealthTableFormat),
		NewExactMatcher(detailedClusterHealthTableFormat),
		NewExactMatcher(cpClusterHealthTableFormat),
		NewExactMatcher(dpClusterHealthTableFormat),
		NewExactMatcher(`{{.CPStatus}}`),
		NewExactMatcher(`{{.DPStatus}}`),
		NewExactMatcher(`{{.Node}}: {{.Status}}`),
		NewExactMatcher(`node: {{.Node}}\nstatus: {{.Status}}\nkv: {{.KV}}\nkv_write: {{.KVWrite}}\nnats: {{.NATS}}\ndfs_client: {{.DirectFSClient}}\ndirector: {{.Director}}\nfs_driver: {{.FSDriver}}\nfs: {{.FS}}\n`),
	)

	if len(status) == 0 {
		return fmt.Errorf("No cluster nodes found")
	}
	render := func(format func(subContext subContext) error) error {
		for _, s := range status {
			if err := format(&clusterHealthContext{v: s}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&clusterHealthContext{}, render)
}

type clusterHealthContext struct {
	HeaderContext
	v *apiTypes.ClusterHealthNode
}

func (c *clusterHealthContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *clusterHealthContext) Node() string {
	c.AddHeader(clusterHealthNodeHeader)
	return c.v.NodeName
}

func (c *clusterHealthContext) healthy() bool {
	return c.cpHealthy() && c.dpHealthy()
}

func (c *clusterHealthContext) cpHealthy() bool {
	return c.v.Submodules.KV.Status+
		c.v.Submodules.KVWrite.Status+
		c.v.Submodules.NATS.Status == "alivealivealive"
}

func (c *clusterHealthContext) cpDegraded() string {
	degraded := []string{}
	if c.v.Submodules.KV.Status != "alive" {
		degraded = append(degraded, clusterHealthKVHeader)
	}
	if c.v.Submodules.KVWrite.Status != "alive" {
		degraded = append(degraded, clusterHealthKVWriteHeader)
	}
	if c.v.Submodules.NATS.Status != "alive" {
		degraded = append(degraded, clusterHealthNATSHeader)
	}
	return "Degraded (" + strings.Join(degraded, ", ") + ")"
}

func (c *clusterHealthContext) dpHealthy() bool {
	return c.v.Submodules.DirectFSInitiator.Status+
		c.v.Submodules.Director.Status+
		c.v.Submodules.Presentation.Status+
		c.v.Submodules.RDB.Status == "alivealivealivealive"
}

func (c *clusterHealthContext) dpDegraded() string {
	degraded := []string{}
	if c.v.Submodules.DirectFSInitiator.Status != "alive" {
		degraded = append(degraded, clusterHealthDirectFSInitiatorHeader)
	}
	if c.v.Submodules.Director.Status != "alive" {
		degraded = append(degraded, clusterHealthDirectorHeader)
	}
	if c.v.Submodules.Presentation.Status != "alive" {
		degraded = append(degraded, clusterHealthPresentationHeader)
	}
	if c.v.Submodules.RDB.Status != "alive" {
		degraded = append(degraded, clusterHealthRDBHeader)
	}
	return "Degraded (" + strings.Join(degraded, ", ") + ")"
}

func (c *clusterHealthContext) Status() string {
	c.AddHeader(clusterHealthStatusHeader)
	if c.healthy() {
		return "Healthy"
	}
	return "Not Ready"
}

func (c *clusterHealthContext) CPStatus() string {
	c.AddHeader(clusterHealthCPStatusHeader)
	if c.cpHealthy() {
		return "Healthy"
	}
	return c.cpDegraded()
}

func (c *clusterHealthContext) DPStatus() string {
	c.AddHeader(clusterHealthDPStatusHeader)
	if c.dpHealthy() {
		return "Healthy"
	}
	return c.dpDegraded()
}

func (c *clusterHealthContext) NATS() string {
	c.AddHeader(clusterHealthNATSHeader)
	if s := c.v.Submodules.NATS.Status; s != "" {
		return s
	}
	return "unknown"
}

func (c *clusterHealthContext) KV() string {
	c.AddHeader(clusterHealthKVHeader)
	if s := c.v.Submodules.KV.Status; s != "" {
		return s
	}
	return "unknown"
}

func (c *clusterHealthContext) KVWrite() string {
	c.AddHeader(clusterHealthKVWriteHeader)
	if s := c.v.Submodules.KVWrite.Status; s != "" {
		return s
	}
	return "unknown"
}

func (c *clusterHealthContext) DirectFSInitiator() string {
	c.AddHeader(clusterHealthDirectFSInitiatorHeader)
	if s := c.v.Submodules.DirectFSInitiator.Status; s != "" {
		return s
	}
	return "unknown"
}

func (c *clusterHealthContext) Director() string {
	c.AddHeader(clusterHealthDirectorHeader)
	if s := c.v.Submodules.Director.Status; s != "" {
		return s
	}
	return "unknown"
}

func (c *clusterHealthContext) Presentation() string {
	c.AddHeader(clusterHealthPresentationHeader)
	if s := c.v.Submodules.Presentation.Status; s != "" {
		return s
	}
	return "unknown"
}

func (c *clusterHealthContext) RDB() string {
	c.AddHeader(clusterHealthRDBHeader)
	if s := c.v.Submodules.RDB.Status; s != "" {
		return s
	}
	return "unknown"
}
