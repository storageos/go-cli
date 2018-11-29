package formatter

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultLoggerStreamQuietFormat = "{{.Msg}}"
	defaultLoggerStreamTableFormat = "table {{.Time}}\t{{.Level}}\t{{.Category}}\t{{.Msg}}\t{{.Error}}"

	timeFormat = time.RFC3339
)

// NewLogStreamFormat returns a format for use with a stream of log messages
func NewLogStreamFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		if quiet {
			return defaultLoggerQuietFormat
		}
		return defaultLoggerStreamTableFormat
	case RawFormatKey:
		return `{{.Text}}`
	}
	return Format(source)
}

// LogStreamWrite writes a formatted stream of log messages
func LogStreamWrite(ctx Context, msg []byte) error {
	render := func(format func(subContext subContext) error) error {

		entry, err := marshalEntry(msg)
		if err != nil {
			return err
		}
		return format(&logStreamContext{entry: entry})
	}
	return ctx.Write(&loggerContext{}, render)
}

// marshalEntry marshals a json message into a logrus log entry
func marshalEntry(msg []byte) (*logrus.Entry, error) {

	var fields map[string]interface{}
	if err := json.Unmarshal(msg, &fields); err != nil {
		return nil, err
	}

	entry := &logrus.Entry{
		Data: logrus.Fields{},
	}

	for k, v := range fields {
		switch k {
		case "time":
			t, err := time.Parse(timeFormat, str(fields["time"]))
			if err != nil {
				t = time.Now()
			}
			entry.Time = t
		case "level":
			level, err := logrus.ParseLevel(str(fields["level"]))
			if err != nil {
				level = logrus.InfoLevel
			}
			entry.Level = level
		case "msg":
			entry.Message = str(v)
		default:
			entry.Data[k] = v
		}

	}
	return entry, nil
}

type logStreamContext struct {
	HeaderContext
	entry *logrus.Entry
}

// Text uses the default logrus text formatter
func (c *logStreamContext) Text() string {

	f := TextFormatter{}

	out, err := f.Format(c.entry)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (c *logStreamContext) Msg() string {
	return c.entry.Message
}

func (c *logStreamContext) Time() string {
	return c.entry.Time.Format(timeFormat)
}

func (c *logStreamContext) Level() string {
	return c.entry.Level.String()
}

func (c *logStreamContext) Host() string {
	if v, ok := c.entry.Data["host"]; ok {
		return str(v)
	}
	return ""
}

func (c *logStreamContext) Module() string {
	if v, ok := c.entry.Data["module"]; ok {
		return str(v)
	}
	return ""
}

func (c *logStreamContext) Category() string {
	if v, ok := c.entry.Data["category"]; ok {
		return str(v)
	}
	return ""
}

func (c *logStreamContext) Action() string {
	if v, ok := c.entry.Data["action"]; ok {
		return str(v)
	}
	return ""
}

func (c *logStreamContext) Error() string {
	if v, ok := c.entry.Data["error"]; ok {
		return str(v)
	}
	return ""
}

func str(v interface{}) string {
	switch vv := v.(type) {
	case string:
		return vv
	default:
		return ""
	}
}
