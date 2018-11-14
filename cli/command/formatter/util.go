package formatter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/docker/go-units"
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
