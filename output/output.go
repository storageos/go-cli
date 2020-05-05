// Package output provides a collection of Displayer types. Displayer types
// are used to abstract away the implementation details of how output is
// written for CLI app commands.
package output

import "time"

const unknownResourceName = "unknown"

// TimeHumanizer represents a generic type able to transform a time.Time object
// into a human-readable string.
type TimeHumanizer interface {
	TimeToHuman(t time.Time) string
}
