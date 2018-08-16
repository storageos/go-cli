package formatter

import (
	"fmt"
	"sort"
	"strings"
)

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
