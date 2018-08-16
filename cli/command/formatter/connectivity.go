package formatter

import (
	"fmt"

	"github.com/storageos/go-api/types"
)

const (
	connectivityTableFormat      = "table {{.Source}}\t{{.Name}}\t{{.Address}}\t{{.Latency}}\t{{.Status}}\t{{.Error}}"
	connectivityTableQuietFormat = "table {{.Test}}\t{{.Status}}"
	connectivityRawFormat        = `source: {{.Source}}\nname: {{.Name}}\naddress: {{.Address}}\nlatency: {{.Latency}}\nstatus: {{.Status}}\nmessage: {{.Error}}\n`
	connectivityRawQuietFormat   = "{{.Test}}: {{.Status}}"
	connectivitySummaryFormat    = "{{.Status}}"

	connectivityNameHeader    = "NAME"
	connectivityAddressHeader = "ADDRESS"
	connectivitySourceHeader  = "SOURCE"
	connectivityLatencyHeader = "LATENCY"
	connectivityStateHeader   = "STATUS"
	connectivityErrorHeader   = "MESSAGE"
	connectivityTestHeader    = "SOURCE->ADDRESS"
)

// NewConnectivityFormat returns a format for use with a connectivity Context
func NewConnectivityFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return connectivityTableQuietFormat
		}
		return connectivityTableFormat
	case RawFormatKey:
		if quiet {
			return connectivityRawQuietFormat
		}
		return connectivityRawFormat
	case SummaryFormatKey:
		return connectivitySummaryFormat
	}
	return Format(source)
}

// ConnectivityWrite writes formatted connectivity results using the Context
func ConnectivityWrite(ctx Context, results types.ConnectivityResults) error {
	render := func(format func(subContext subContext) error) error {
		switch ctx.Trunc {
		case true:
			if err := format(&connectivityContext{ok: results.IsOK()}); err != nil {
				return err
			}
		case false:
			for _, cr := range results {
				if err := format(&connectivityContext{result: cr}); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return ctx.Write(&connectivityContext{}, render)
}

type connectivityContext struct {
	HeaderContext
	result types.ConnectivityResult
	ok     bool
}

func (c *connectivityContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *connectivityContext) Name() string {
	c.AddHeader(connectivityNameHeader)
	return c.result.Label
}

func (c *connectivityContext) Address() string {
	c.AddHeader(connectivityAddressHeader)
	return c.result.Address
}

func (c *connectivityContext) Source() string {
	c.AddHeader(connectivitySourceHeader)
	return c.result.Source
}

func (c *connectivityContext) Latency() string {
	c.AddHeader(connectivityLatencyHeader)
	return fmt.Sprint(c.result.LatencyNS)
}

func (c *connectivityContext) Error() string {
	c.AddHeader(connectivityErrorHeader)
	return c.result.Error
}

func (c *connectivityContext) Test() string {
	c.AddHeader(connectivityTestHeader)
	return c.result.Source + "->" + c.result.Address
}

func (c *connectivityContext) Status() string {
	c.AddHeader(connectivityStateHeader)

	// Return result status if single result
	if c.result.Address != "" {
		if c.result.IsOK() {
			return "OK"
		}
		return "ERROR"
	}

	// If we only have a summary, use that
	if c.ok {
		return "OK"
	}
	return "ERROR"
}
