package output

import (
	"fmt"
	"strings"
)

// Format is the type of a generic output format
type Format uint8

// Different output format to print results.
const (
	Unknown Format = iota
	JSON
	YAML
	Text
	// Keep update the `ValidFormats` list on adding new
)

var (
	// ValidFormats is the list of formats accepted. This is used in error
	// returned and in flag option description.
	ValidFormats = []string{
		"json",
		"yaml",
		"text",
	}

	// ErrInvalidFormat is an error that is returned when the string passed to
	// FormatFromString cannot be parsed in a valid output format.
	ErrInvalidFormat = fmt.Errorf("invalid output format string. Use one of %v", ValidFormats)
)

// FormatFromString parses a string to understand which output format has been
// selected. Case insensitive.
func FormatFromString(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "json":
		return JSON, nil
	case "yaml":
		return YAML, nil
	case "text":
		return Text, nil
	// Keep update the `ValidFormats` list on adding new

	default:
		return Unknown, ErrInvalidFormat
	}
}

// String returns the string representation of the output Format.
func (t Format) String() string {
	switch t {
	case JSON:
		return "json"
	case YAML:
		return "yaml"
	case Text:
		return "text"
	// Keep update the `ValidFormats` list on adding new

	default:
		return "unknown"
	}
}
