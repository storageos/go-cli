package textformat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/gosuri/uitable"

	"code.storageos.net/storageos/c2-cli/output"
)

// Displayer is a type which creates human-readable strings and writes them to
// io.Writers.
type Displayer struct {
	timeHumanizer output.TimeHumanizer
}

// NewDisplayer initialises a Displayer which prints human readable strings
// StorageOS to output CLI results.
func NewDisplayer(timeHumanizer output.TimeHumanizer) *Displayer {
	return &Displayer{
		timeHumanizer: timeHumanizer,
	}
}

func (d *Displayer) timeToHuman(t time.Time) string {
	humanized := d.timeHumanizer.TimeToHuman(t)
	rfc := t.Format(time.RFC3339)
	return fmt.Sprintf("%s (%s)", rfc, humanized)
}

const format pb.ProgressBarTemplate = `{{counters . }} {{bar . "[" "#" "#" "." "]"}} {{percent . }}  -  {{string . "suffix"}}`

func syncProgressBarString(current, max, secondsRemaining uint64) (string, error) {
	var maxInt64 uint64 = math.MaxInt64
	if current > maxInt64 || max > maxInt64 {
		return "", errors.New("invalid sync progress value received: int64 overflow")
	}

	if current > max {
		return "", errors.New("invalid sync progress value received: current < max")
	}

	bar := format.Start64(int64(max))
	bar.Set(pb.Bytes, true)
	bar.SetWidth(80)
	bar.SetCurrent(int64(max - current))

	etaString := time.Duration(secondsRemaining) * time.Second
	bar.Set("suffix", fmt.Sprintf("ETA: %s", etaString))
	return bar.String(), nil
}

func disableToHuman(b bool) string {
	if b {
		return "Disabled"
	}
	return "Enabled"
}

var (
	nodeHeaders        = []interface{}{"NAME", "HEALTH", "AGE", "LABELS"}
	namespaceHeaders   = []interface{}{"NAME", "AGE"}
	userHeaders        = []interface{}{"NAME", "ROLE", "AGE", "GROUPS"}
	volumeHeaders      = []interface{}{"NAMESPACE", "NAME", "SIZE", "LOCATION", "ATTACHED ON", "REPLICAS", "AGE"}
	policyGroupHeaders = []interface{}{"NAME", "USERS", "SPECS", "AGE"}
)

func createTable(headers []interface{}) (*uitable.Table, func(io.Writer) error) {
	table := uitable.New()
	table.MaxColWidth = 50
	table.Separator = "  "

	// header
	if headers != nil {
		table.AddRow(headers...)
	}

	return table, func(w io.Writer) error {
		_, err := fmt.Fprintln(w, table)
		return err
	}
}

// AsyncRequest writes nothing to w.
func (d *Displayer) AsyncRequest(ctx context.Context, w io.Writer) error {
	_, err := fmt.Fprintln(w, "Async request accepted.")
	return err
}
