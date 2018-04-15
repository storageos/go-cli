package formatter

import (
	"github.com/storageos/go-api/types"
)

const (
	defaultNodeConnectivityTableFormat = "table {{.Name}}\t{{.Status}}\t{{.API}}\t{{.Nats}}\t{{.Etcd}}"

	connectivityNodeNameHeader = "NAME"
	connectivityStateHeader    = "STATUS"
	connectivityAPIHeader      = "API"
	connectivityNatsHeader     = "NATS"
	connectivityEtcdHeader     = "ETCD"
)

// NewNodeConnectivityFormat returns a format for use with a node connectivity Context
func NewNodeConnectivityFormat(source string) Format {
	switch source {
	case TableFormatKey:
		return defaultNodeConnectivityTableFormat

	case RawFormatKey:
		return `name: {{.Name}}\nstatus: {{.Status}}\napi: {{.API}}\nnats: {{.Nats}}\netcd: {{.Etcd}}\n`
	}

	return Format(source)
}

// NodeConnectivityWrite writes formatted node connectivity results using the Context
func NodeConnectivityWrite(ctx Context, result []types.ConnectivityResult) error {
	render := func(format func(subContext subContext) error) error {
		for _, cr := range result {
			if err := format(&nodeConnectivityContext{result: cr}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&nodeConnectivityContext{}, render)
}

type nodeConnectivityContext struct {
	HeaderContext
	result types.ConnectivityResult
}

func (c *nodeConnectivityContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *nodeConnectivityContext) Name() string {
	c.AddHeader(nodeNameHeader)
	return c.result.Target.Name
}

func (c *nodeConnectivityContext) Status() string {
	c.AddHeader(connectivityStateHeader)
	if c.result.Passes() {
		return "OK"
	}
	return "Fail"
}

func (c *nodeConnectivityContext) API() string {
	c.AddHeader(connectivityAPIHeader)
	if c.result.Connectivity[types.APIConnectivity] {
		return "OK"
	}
	return "Fail"
}

func (c *nodeConnectivityContext) Nats() string {
	c.AddHeader(connectivityNatsHeader)
	if c.result.Connectivity[types.NatsConnectivity] {
		return "OK"
	}
	return "Fail"
}

func (c *nodeConnectivityContext) Etcd() string {
	c.AddHeader(connectivityEtcdHeader)
	if c.result.Connectivity[types.EtcdConnectivity] {
		return "OK"
	}
	return "Fail"
}
