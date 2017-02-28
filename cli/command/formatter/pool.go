package formatter

import (
	"fmt"
	"strings"

	"github.com/storageos/go-api/types"
)

const (
	defaultPoolQuietFormat = "{{.Name}}"
	defaultPoolTableFormat = "table {{.Driver}}\t{{.Name}}"

	poolNameHeader = "POOL NAME"
	// Status header ?
)

// NewPoolFormat returns a format for use with a pool Context
func NewPoolFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultPoolQuietFormat
		}
		return defaultPoolTableFormat
	case RawFormatKey:
		if quiet {
			return `name: {{.Name}}`
		}
		return `name: {{.Name}}\ndriver: {{.Driver}}\n`
	}
	return Format(source)
}

// PoolWrite writes formatted pools using the Context
func PoolWrite(ctx Context, pools []*types.Pool) error {
	render := func(format func(subContext subContext) error) error {
		for _, pool := range pools {
			if err := format(&poolContext{v: *pool}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&poolContext{}, render)
}

type poolContext struct {
	HeaderContext
	v types.Pool
}

func (c *poolContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *poolContext) Name() string {
	c.AddHeader(poolNameHeader)
	return c.v.Name
}

func (c *poolContext) Labels() string {
	c.AddHeader(labelsHeader)
	if c.v.Labels == nil {
		return ""
	}

	var joinLabels []string
	for k, v := range c.v.Labels {
		joinLabels = append(joinLabels, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(joinLabels, ",")
}

func (c *poolContext) Label(name string) string {

	n := strings.Split(name, ".")
	r := strings.NewReplacer("-", " ", "_", " ")
	h := r.Replace(n[len(n)-1])

	c.AddHeader(h)

	if c.v.Labels == nil {
		return ""
	}
	return c.v.Labels[name]
}
