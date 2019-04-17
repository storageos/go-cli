package formatter

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/docker/go-units"
	"github.com/storageos/go-cli/cli/command/inspect"
	"github.com/storageos/go-cli/pkg/templates"
)

// GiB is 1024 * 1024 * 1024
const GiB = 1024 * 1024 * 1024

func writeLabels(labels map[string]string) string {
	var joinLabels []string
	for k, v := range labels {
		joinLabels = append(joinLabels, fmt.Sprintf("%s=%s", k, v))
	}

	sort.SliceStable(joinLabels, func(i, j int) bool {
		fst := strings.Split(joinLabels[i], "=")[0]
		snd := strings.Split(joinLabels[j], "=")[0]
		return fst < snd
	})

	return strings.Join(joinLabels, ",")
}

// bytesSize returns a human-readable size in bytes, kibibytes,
// mebibytes, gibibytes, or tebibytes (eg. "44kiB", "17MiB").
// Ref: https://en.wikipedia.org/wiki/Binary_prefix.
//
// it should be used by all size display including pool, node and volume
func bytesSize(size uint64) string {
	return units.BytesSize(float64(size))
}

// TryFormatSpec may print in using format as a template.
//
// If format appears to be a valid template, it is used to print in and the
// process exits with a return code of 0. If the format is "help" the format
// help text is printed for in and the process exits with 0, see
// templates.HelpText. If format appears to be a template, but does not compile,
// an error is printed and the process exits with a return code of 1.
//
// If the format does not resemble a template, this function does nothing.
func TryFormatSpec(format string, in interface{}) {
	if strings.ToLower(format) == "help" {
		fmt.Println(templates.HelpText(in))
		os.Exit(0)
	}

	if !strings.Contains(format, "{{") {
		return
	}

	printer, err := inspect.NewTemplateInspectorFromString(os.Stdout, format)
	if err != nil {
		fmt.Printf("There was an error compiling the provided template:\n\t%v\n", err)
		os.Exit(1)
	}

	if err := printer.Inspect(in, nil); err != nil {
		logrus.WithError(err).Debug("template error")
	}
	printer.Flush()
	os.Exit(0)
}

// TableMatcher is a matcher used for TryFormatUnlessMatches, which matches any format starting
// with a table function.
var TableMatcher = regexp.MustCompile(`table.*`)

// FormatMatcher is an interface that describes a type that can match on a format string
type FormatMatcher interface {
	MatchString(string) bool
}

type exactMatcher struct{ s string }

func (f *exactMatcher) MatchString(s string) bool { return s == f.s }

// NewExactMatcher returns a new format matcher that will only match if the provided strings are byte equivalent
func NewExactMatcher(s string) FormatMatcher {
	return &exactMatcher{s}
}

// TryFormatUnlessMatches calls TryFormatSpec unless format satisfies one of the FormatMatcher types return true
//
// This works around some of the formatters using templates that are
// incompatible with the results (using a subformatter).
func TryFormatUnlessMatches(format string, in interface{}, notIfMatch ...FormatMatcher) {
	for _, r := range notIfMatch {
		if r.MatchString(format) {
			return
		}
	}

	TryFormatSpec(format, in)
}
