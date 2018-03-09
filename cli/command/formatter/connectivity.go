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
func ConnectivityWrite(ctx Context, result []types.ConnectivityResult) error {
	render := func(format func(subContext subContext) error) error {
		for _, cr := range result {
			if err := format(&connectivityContext{result: cr}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&connectivityContext{}, render)
}

type connectivityContext struct {
	HeaderContext
	result types.ConnectivityResult
}

func (c *connectivityContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *connectivityContext) Name() string {
	c.AddHeader(nodeNameHeader)
	return c.result.Target.Name
}

func (c *connectivityContext) Status() string {
	c.AddHeader(connectivityStateHeader)
	if c.result.Pass {
		return "OK"
	}
	return "Fail"
}

func (c *connectivityContext) API() string {
	c.AddHeader(connectivityAPIHeader)
	if c.result.APIConnectivity {
		return "OK"
	}
	return "Fail"
}
