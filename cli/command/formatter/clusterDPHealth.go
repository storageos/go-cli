package formatter

import (
	"github.com/storageos/go-api/types"
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
func ClusterHealthDPWrite(ctx Context, clusterHealth *types.ClusterHealthDP) error {
	render := func(format func(subContext subContext) error) error {
		for _, status := range *clusterHealth {
			if err := format(&dpHealthContext{v: status}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&dpHealthContext{}, render)
}

type dpHealthContext struct {
	HeaderContext
	v types.DPHealthStatusWithID
}

func (d *dpHealthContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(d)
}

func (d *dpHealthContext) Node() string {
	d.AddHeader(clusterHealthDPNodeHeader)
	return d.v.ID
}

func (d *dpHealthContext) healthy() bool {
	return d.v.DirectFSClient.Status+d.v.DirectFSServer.Status+d.v.Director.Status+d.v.FSDriver.Status+d.v.FS.Status == "alivealivealivealivealive"
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
	return d.v.DirectFSClient.Status
}

func (d *dpHealthContext) DirectFSServer() string {
	d.AddHeader(clusterHealthDPDirectFSServerHeader)
	return d.v.DirectFSServer.Status
}

func (d *dpHealthContext) Director() string {
	d.AddHeader(clusterHealthDPDirectorHeader)
	return d.v.Director.Status
}

func (d *dpHealthContext) FSDriver() string {
	d.AddHeader(clusterHealthDPFSDriverHeader)
	return d.v.FSDriver.Status
}

func (d *dpHealthContext) FS() string {
	d.AddHeader(clusterHealthDPFSHeader)
	return d.v.FS.Status
}
