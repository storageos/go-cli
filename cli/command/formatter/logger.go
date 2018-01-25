package formatter

import (
	"fmt"
	"strings"

	"github.com/storageos/go-cli/api/types"
)

const (
	defaultLoggerQuietFormat = "{{.Level}}"
	defaultLoggerTableFormat = "table {{.Node}}\t{{.Level}}\t{{.Filter}}"

	loggerNodeHeader       = "NODE"
	loggerLevelHeader      = "LEVEL"
	loggerFilterHeader     = "FILTER"
	loggerCategoriesHeader = "CATEGORIES"
)

// NewLoggerFormat returns a format for use with a logger Context
func NewLoggerFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultLoggerQuietFormat
		}
		return defaultLoggerTableFormat
	case RawFormatKey:
		if quiet {
			return `name: {{.Level}}`
		}
		return `node: {{.Node}}\nlevel: {{.Level}}\nformat: {{.Filter}}\ncategories: {{.Categories}}\n`
	}
	return Format(source)
}

// LoggerWrite writes formatted loggers using the Context
func LoggerWrite(ctx Context, loggers []*types.Logger) error {
	render := func(format func(subContext subContext) error) error {
		for _, logger := range loggers {
			if err := format(&loggerContext{v: *logger}); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&loggerContext{}, render)
}

type loggerContext struct {
	HeaderContext
	v types.Logger
}

func (c *loggerContext) MarshalJSON() ([]byte, error) {
	return marshalJSON(c)
}

func (c *loggerContext) Node() string {
	c.AddHeader(loggerNodeHeader)
	return c.v.Node
}

func (c *loggerContext) Level() string {
	c.AddHeader(loggerLevelHeader)
	return c.v.Level
}

func (c *loggerContext) Filter() string {
	c.AddHeader(loggerFilterHeader)
	return c.v.Filter
}

func (c *loggerContext) Categories() string {
	c.AddHeader(loggerCategoriesHeader)

	var out []string
	for cat, level := range c.v.Categories {
		out = append(out, fmt.Sprintf("%s=%s", cat, level))
	}
	return strings.Join(out, ",")
}
