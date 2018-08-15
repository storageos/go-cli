package formatter

import (
	"fmt"

	"github.com/storageos/go-api/types"
)

const (
	connectivityTableFormat = "table {{.Source}}\t{{.Name}}\t{{.Address}}\t{{.Latency}}\t{{.Status}}\t{{.Error}}\t"
	connectivityRawFormat   = `source: {{.Source}}\nname: {{.Name}}\naddress: {{.Address}}\nlatency: {{.Latency}}\nstatus: {{.Status}}\nmessage: {{.Error}}\n`

	connectivityNameHeader    = "NAME"
	connectivityAddressHeader = "ADDRESS"
	connectivitySourceHeader  = "SOURCE"
	connectivityLatencyHeader = "LATENCY"
	connectivityStateHeader   = "STATUS"
	connectivityErrorHeader   = "MESSAGE"
)

// NewConnectivityFormat returns a format for use with a connectivity Context
func NewConnectivityFormat(source string, quiet bool) Format {

	// Quiet should return OK/ERROR summary
	if quiet {
		return "{{.OK}}"
	}

	switch source {
	case TableFormatKey:
		return connectivityTableFormat

	case RawFormatKey:
		return connectivityRawFormat
	}

	return Format(source)
}

// ConnectivityWriteSummary writes a formatted connectivity summary using the Context
func ConnectivityWriteSummary(ctx Context, ok bool) error {
	render := func(format func(subContext subContext) error) error {
		if err := format(&connectivitySummaryContext{ok: ok}); err != nil {
			return err
		}
		return nil
	}
	return ctx.Write(&connectivitySummaryContext{}, render)
}

// ConnectivityWrite writes formatted connectivity results using the Context
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

type connectivitySummaryContext struct {
	HeaderContext
	ok bool
}

func (c *connectivitySummaryContext) OK() string {
	if c.ok {
		return "OK"
	}
	return "ERROR"
}

type connectivityContext struct {
	HeaderContext
	result types.ConnectivityResult
}

func (c connectivityContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c connectivityContext) Name() string {
	c.AddHeader(connectivityNameHeader)
	return c.result.Label
}

func (c connectivityContext) Address() string {
	c.AddHeader(connectivityAddressHeader)
	return c.result.Address
}

func (c connectivityContext) Source() string {
	c.AddHeader(connectivitySourceHeader)
	return c.result.Source
}

func (c connectivityContext) Latency() string {
	c.AddHeader(connectivityLatencyHeader)
	return fmt.Sprint(c.result.LatencyNS)
}

func (c connectivityContext) Error() string {
	c.AddHeader(connectivityErrorHeader)
	return c.result.Error
}

func (c connectivityContext) Status() string {
	c.AddHeader(connectivityStateHeader)
	if c.result.Passes() {
		return "OK"
	}
	return "ERROR"
}
