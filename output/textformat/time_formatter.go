package textformat

import (
	"time"

	"github.com/dustin/go-humanize"
)

// TimeFormatter is a type able to transform a time.Time into a human-readable
// string. It uses the go-humanize external package.
type TimeFormatter struct{}

// TimeToHuman transforms a time.Time object into a human-readable string.
func (t2 *TimeFormatter) TimeToHuman(t time.Time) string {
	return humanize.Time(t)
}

// NewTimeFormatter creates a new TimeFormatter object.
func NewTimeFormatter() *TimeFormatter {
	return &TimeFormatter{}
}
