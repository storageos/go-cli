package formatter

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

const defaultTimestampFormat = time.RFC3339

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 36
	gray    = 37
)

var (
	baseTimestamp time.Time
)

func init() {
	baseTimestamp = time.Now()
}

// TextFormatter formats logs into text
type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Disable timestamp logging. Useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	sync.Once
}

func (f *TextFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}
}

// Format renders a single log entry
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	prefixFieldClashes(entry.Data)

	f.Do(func() { f.init(entry) })

	isColored := (f.ForceColors || f.isTerminal) && !f.DisableColors

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	if isColored {
		if err := f.printColored(b, entry, keys, timestampFormat); err != nil {
			return nil, err
		}
	} else {
		if !f.DisableTimestamp {
			if err := f.appendKeyValue(b, "time", entry.Time.Format(timestampFormat)); err != nil {
				return nil, err
			}
		}
		if err := f.appendKeyValue(b, "level", entry.Level.String()); err != nil {
			return nil, err
		}

		// Add StorageOS custom fields
		if host, ok := entry.Data["host"]; ok {
			if err := f.appendKeyValue(b, "host", host); err != nil {
				return nil, err
			}
			delete(entry.Data, "host")
		}
		if module, ok := entry.Data["module"]; ok {
			if err := f.appendKeyValue(b, "module", module); err != nil {
				return nil, err
			}
			delete(entry.Data, "module")
		}
		if category, ok := entry.Data["category"]; ok {
			if err := f.appendKeyValue(b, "category", category); err != nil {
				return nil, err
			}
			delete(entry.Data, "category")
		}

		if entry.Message != "" {
			if err := f.appendKeyValue(b, "msg", entry.Message); err != nil {
				return nil, err
			}
		}
		for _, key := range keys {
			if err := f.appendKeyValue(b, key, entry.Data[key]); err != nil {
				return nil, err
			}
		}
	}

	if err := b.WriteByte('\n'); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat string) error {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	levelText := strings.ToUpper(entry.Level.String())[0:4]

	if f.DisableTimestamp {
		err := printColoredWithoutTimestamp(b, levelColor, levelText, entry.Message)
		if err != nil {
			return err
		}
	} else if !f.FullTimestamp {
		err := printColoredWithDuration(b, levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message)
		if err != nil {
			return err
		}
	} else {
		err := printColoredWithTimestamp(b, levelColor, levelText, entry.Time.Format(timestampFormat), entry.Message)
		if err != nil {
			return err
		}
	}

	for _, k := range keys {
		v := entry.Data[k]
		if err := printColoredKey(b, levelColor, k); err != nil {
			return err
		}

		if err := f.appendValue(b, v); err != nil {
			return err
		}
	}
	return nil
}

func printColoredWithoutTimestamp(b *bytes.Buffer, levelColor int, levelText string, message string) error {
	_, err := fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m %-44s ", levelColor, levelText, message)
	return err
}

func printColoredWithDuration(b *bytes.Buffer, levelColor int, levelText string, duration int, message string) error {
	_, err := fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d] %-44s ", levelColor, levelText, duration, message)
	return err
}

func printColoredWithTimestamp(b *bytes.Buffer, levelColor int, levelText string, timestamp string, message string) error {
	_, err := fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s] %-44s ", levelColor, levelText, timestamp, message)
	return err
}

func printColoredKey(b *bytes.Buffer, levelColor int, key string) error {
	_, err := fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, key)
	return err
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) error {
	if b.Len() > 0 {
		if err := b.WriteByte(' '); err != nil {
			return err
		}
	}

	if _, err := b.WriteString(key); err != nil {
		return err
	}

	if err := b.WriteByte('='); err != nil {
		return err
	}

	return f.appendValue(b, value)
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) error {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		if _, err := b.WriteString(stringVal); err != nil {
			return err
		}
	} else {
		if _, err := b.WriteString(fmt.Sprintf("%q", stringVal)); err != nil {
			return err
		}
	}
	return nil
}

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

// This is to not silently overwrite `time`, `msg` and `level` fields when
// dumping it. If this code wasn't there doing:
//
//  logrus.WithField("level", 1).Info("hello")
//
// Would just silently drop the user provided level. Instead with this code
// it'll logged as:
//
//  {"level": "info", "fields.level": 1, "msg": "hello", "time": "..."}
//
// It's not exported because it's still using Data in an opinionated way. It's to
// avoid code duplication between the two default formatters.
func prefixFieldClashes(data logrus.Fields) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["msg"]; ok {
		data["fields.msg"] = m
	}

	if l, ok := data["level"]; ok {
		data["fields.level"] = l
	}
}
