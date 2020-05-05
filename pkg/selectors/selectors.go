// Package selectors implements a mechanism for filtering StorageOS API
// resources by the state of their label set.
package selectors

import (
	"errors"
	"fmt"
	"strings"

	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/volume"
)

var (
	// ErrInvalidSelectorFormat is returned if a selector does not conform to
	// a valid format for a label key or label pair.
	ErrInvalidSelectorFormat = errors.New("invalid selector (must match key or key=value formats)")
)

// selector is a function which returns a boolean indicating if labelSet
// matches some arbitrary constraint which defines the selector.
type selector func(labelSet labels.Set) bool

// Set holds a collection of selectors, implementing functionality to filter
// lists of resource types using them.
//
// FilterXXX methods on a Set must return the parameter list as-is when no
// the Set does not have any selectors.
type Set struct {
	selectors []selector
}

// FilterNodes returns the subset of node resources in nodes which match
// on the selectors which s is configured with.
func (s *Set) FilterNodes(nodes []*node.Resource) []*node.Resource {
	if len(s.selectors) == 0 {
		return nodes
	}

	filtered := make([]*node.Resource, 0)

NextNode:
	for _, n := range nodes {
		for _, selects := range s.selectors {
			if !selects(n.Labels) {
				continue NextNode
			}
		}
		filtered = append(filtered, n)
	}

	return filtered
}

// FilterNamespaces returns the subset of namespace resources in namespaces
// which match on the selectors which s is configured with.
func (s *Set) FilterNamespaces(namespaces []*namespace.Resource) []*namespace.Resource {
	if len(s.selectors) == 0 {
		return namespaces
	}

	filtered := make([]*namespace.Resource, 0)

NextNamespace:
	for _, ns := range namespaces {
		for _, selects := range s.selectors {
			if !selects(ns.Labels) {
				continue NextNamespace
			}
		}
		filtered = append(filtered, ns)
	}

	return filtered
}

// FilterVolumes returns the subset of volume resources in volumes which match
// on the selectors which s is configured with.
func (s *Set) FilterVolumes(volumes []*volume.Resource) []*volume.Resource {
	if len(s.selectors) == 0 {
		return volumes
	}

	filtered := make([]*volume.Resource, 0)

NextVolume:
	for _, v := range volumes {
		for _, selected := range s.selectors {
			if !selected(v.Labels) {
				continue NextVolume
			}
		}
		filtered = append(filtered, v)
	}

	return filtered
}

// NewSetFromStrings constructs a new selector Set using each item in
// selectors. If any of the provided selector strings are not valid then an
// error is returned.
//
// If given no selectors or the empty string then the returned selector
// filtering are no-ops.
func NewSetFromStrings(selectors ...string) (*Set, error) {
	selectFns := make([]selector, 0, len(selectors))
	for _, selector := range selectors {
		parts := strings.Split(selector, "=")

		fn, err := newSelector(parts)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", err, selector)
		}

		selectFns = append(selectFns, fn)
	}
	return &Set{
		selectors: selectFns,
	}, nil
}

func newSelector(parts []string) (selector, error) {
	switch len(parts) {
	case 1:
		if parts[0] == "" {
			return nil, ErrInvalidSelectorFormat
		}

		return func(labelSet labels.Set) bool {
			_, exists := labelSet[parts[0]]
			return exists
		}, nil
	case 2:
		for parts[0] == "" || parts[1] == "" {
			return nil, ErrInvalidSelectorFormat
		}

		return func(labelSet labels.Set) bool {
			value, exists := labelSet[parts[0]]
			if !exists {
				return false
			}
			return value == parts[1]
		}, nil
	default:
		return nil, ErrInvalidSelectorFormat
	}
}
