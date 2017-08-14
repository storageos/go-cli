package formatter

import (
	"fmt"

	"github.com/storageos/go-cli/types"
)

const (
	defaultHealthDPTableFormat = "table {{.Node}}\t{{.Status}}\t{{.DirectFSClient}}\t{{.DirectFSServer}}\t{{.Director}}\t{{.FSDriver}}\t{{.FS}}"

	clusterHealthDPNodeHeader           = "NODE"
	clusterHealthDPStatusHeader         = "STATUS"
	clusterHealthDPDirectFSClientHeader = "DFS_CLIENT"
	clusterHealthDPDirectFSServerHeader = "DFS_SERVER"
	clusterHealthDPDirectorHeader       = "DIRECTOR"
	clusterHealthDPFSDriverHeader       = "FS_DRIVER"
	clusterHealthDPFSHeader             = "FS"
)

// NewHealthDPFormat returns a format for use with a dpHealth Context
func NewHealthDPFormat(source string) Format {
	switch source {
	case TableFormatKey:
		return defaultHealthDPTableFormat
	case RawFormatKey:
		return `node: {{.Node}}\nstatus: {{.Status}}\ndfs_client: {{.DirectFSClient}}\ndfs_server: {{.DirectFSServer}}\ndirector: {{.Director}}\nfs_driver: {{.FSDriver}}\nfs: {{.FS}}\n`
	}
	return Format(source)
}

// ClusterHealthDPWrite writes formatted ClusterHealthDP elements using the Context
func ClusterHealthDPWrite(ctx Context, nodes []*types.Node) error {
	if len(nodes) == 0 {
		return fmt.Errorf("No cluster nodes found")
	}
	render := func(format func(subContext subContext) error) error {
		for _, node := range nodes {
			if err := format(&dpHealthContext{v: node}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&dpHealthContext{}, render)
}

type dpHealthContext struct {
	HeaderContext
	v *types.Node
}

func (d *dpHealthContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(d)
}

func (d *dpHealthContext) Node() string {
	d.AddHeader(clusterHealthDPNodeHeader)
	return d.v.Name
}

func (d *dpHealthContext) healthy() bool {
	return d.v.Health.DP.DirectFSClient.Status+
		d.v.Health.DP.DirectFSServer.Status+
		d.v.Health.DP.Director.Status+
		d.v.Health.DP.FSDriver.Status+
		d.v.Health.DP.FS.Status == "alivealivealivealivealive"
}

func (d *dpHealthContext) Status() string {
	d.AddHeader(clusterHealthDPStatusHeader)
	if d.healthy() {
		return "Healthy"
	}
	return "Not Ready"
}

func (d *dpHealthContext) DirectFSClient() string {
	d.AddHeader(clusterHealthDPDirectFSClientHeader)
	return d.v.Health.DP.DirectFSClient.Status
}

func (d *dpHealthContext) DirectFSServer() string {
	d.AddHeader(clusterHealthDPDirectFSServerHeader)
	return d.v.Health.DP.DirectFSServer.Status
}

func (d *dpHealthContext) Director() string {
	d.AddHeader(clusterHealthDPDirectorHeader)
	return d.v.Health.DP.Director.Status
}

func (d *dpHealthContext) FSDriver() string {
	d.AddHeader(clusterHealthDPFSDriverHeader)
	return d.v.Health.DP.FSDriver.Status
}

func (d *dpHealthContext) FS() string {
	d.AddHeader(clusterHealthDPFSHeader)
	return d.v.Health.DP.FS.Status
}
