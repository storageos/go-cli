// Package labels exports an encapsulation of a label set which may be
// associated with certain StorageOS API resources.
package labels

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	// ErrInvalidLabelFormat indicates that a provided label is not using the
	// correct format.
	ErrInvalidLabelFormat = errors.New("invalid label (must match key=value format)")
	// ErrLabelKeyConflict is an error indicating that a set of label pairs
	// used to construct a Set contains a duplicated label key.
	ErrLabelKeyConflict = errors.New("conflict for provided label key")
)

// NewErrInvalidLabelFormatWithDetails wraps an ErrInvalidLabelFormat with
// message details.
func NewErrInvalidLabelFormatWithDetails(details string) error {
	return fmt.Errorf("%w: %s", ErrInvalidLabelFormat, details)
}

// NewErrLabelKeyConflictWithDetails wraps an ErrLabelKeyConflict with
// message details.
func NewErrLabelKeyConflictWithDetails(details string) error {
	return fmt.Errorf("%w: %s", ErrLabelKeyConflict, details)
}

// Set provides a typed wrapper for a label map.
type Set map[string]string

func (s Set) String() string {
	if len(s) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.Grow(len(s) * 40) // Just estimating for performance

	// extract keys and order them in place to make the print deterministic
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// print keys/values ordered alphabetically by key
	for _, k := range keys {
		sb.WriteString(k + "=" + s[k] + ",")
	}
	return strings.TrimRight(sb.String(), ",")
}

// UpsertExisting sets the keys in toEdit to the values found in new, adding any
// of the missing values
func UpsertExisting(toEdit, new Set) Set {

	for k, v := range new {
		toEdit[k] = v
	}

	return toEdit
}

// RemoveExisting deletes the keys contained in toDelete
func RemoveExisting(toEdit, toDelete Set) Set {

	for k := range toDelete {
		delete(toEdit, k)
	}

	return toEdit
}

// FromString parses a given string to try and extract a label set.
//
// it expects the labels in the given string to be comma separated
func FromString(str string) (Set, error) {

	if str == "" {
		// don't error if the string given is empty
		return map[string]string{}, nil
	}

	splitLabels := strings.Split(str, ",")
	return NewSetFromPairs(splitLabels)
}

// NewSetFromPairs constructs a label set from pairs, returning an
// error if any of the provided items is not a key=value pair.
func NewSetFromPairs(pairs []string) (Set, error) {

	set := map[string]string{}

	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		switch len(parts) {

		case 2:
			// Duplicate pair given, return a conflict error
			if _, exists := set[parts[0]]; exists {
				return nil, NewErrLabelKeyConflictWithDetails(parts[0])
			}

			// If either the key or the label is empty then return a format error
			if parts[0] == "" || parts[1] == "" {
				return nil, NewErrInvalidLabelFormatWithDetails(pair)
			}

			set[parts[0]] = parts[1]

		default:
			return nil, NewErrInvalidLabelFormatWithDetails(pair)
		}
	}

	return set, nil
}
