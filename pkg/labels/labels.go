package labels

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidLabelFormat = errors.New("invalid label format")

type Set map[string]string

// SetFromPairs constructs a label set from labelPairs, returning an
// error if any of the provided items is not a key=value pair.
func SetFromPairs(labelPairs []string) (Set, error) {
	set := map[string]string{}

	for _, pair := range labelPairs {
		parts := strings.Split(pair, "=")
		switch len(parts) {
		case 2:
			set[parts[0]] = parts[1]
		default:
			return nil, fmt.Errorf("%w: %s", ErrInvalidLabelFormat, pair)
		}
	}

	return set, nil
}
