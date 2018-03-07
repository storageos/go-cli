package formatter

import (
	"github.com/storageos/go-api/types"
)

const (
	defaultConnectivityTableFormat = "table {{.Name}}\t{{.Status}}\t{{.API}}"

	connectivityNodeNameHeader = "NAME"
	connectivityStateHeader    = "STATUS"
	connectivityAPIHeader      = "API"
)

// NewConnectivityFormat returns a format for use with a node conectivity Context
func NewConnectivityFormat(source string) Format {
	switch source {
	case TableFormatKey:
		return defaultConnectivityTableFormat

	case RawFormatKey:
		return `name: {{.Name}}\nstatus: {{.Status}}\napi: {{.API}}\n`
	}

	return Format(source)
}

// ConnectivityWrite writes formatted node conectivities using the Context
func ConnectivityWrite(ctx Context, result *types.NodeConnectivity) error {
	render := func(format func(subContext subContext) error) error {
		for node, status := range result.Nodes {
			if err := format(&connectivityContext{nodeName: node, status: status}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&connectivityContext{}, render)
}

type connectivityContext struct {
	HeaderContext
	status   types.ConnectivityStatus
	nodeName string
}

func (c *connectivityContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *connectivityContext) Name() string {
	c.AddHeader(nodeNameHeader)
	return c.nodeName
}

func (c *connectivityContext) Status() string {
	c.AddHeader(connectivityStateHeader)
	if c.status.CanConnect {
		return "OK"
	}
	return "Fail"
}

func (c *connectivityContext) API() string {
	c.AddHeader(connectivityAPIHeader)
	for _, name := range c.status.FailedTests {
		if name == "API" {
			return "Fail"
		}
	}
	return "OK"
}
