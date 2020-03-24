package output

import (
	"time"

	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

// Namespace defines a type that contains all the info we need to output a
// namespace.
type Namespace struct {
	ID     id.Namespace `json:"id"`
	Name   string       `json:"name"`
	Labels labels.Set   `json:"labels"`

	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Version   version.Version `json:"version"`
}

// NewNamespace returns a new Namespace object that contains all the info needed
// to be outputted.
func NewNamespace(n *namespace.Resource) *Namespace {
	return &Namespace{
		ID:        n.ID,
		Name:      n.Name,
		Labels:    n.Labels,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Version:   n.Version,
	}
}

// NewNamespaces returns a list of Namespace objects that contains all the info
// needed to be outputted.
func NewNamespaces(ns []*namespace.Resource) []*Namespace {
	namespaces := make([]*Namespace, 0, len(ns))
	for _, n := range ns {
		namespaces = append(namespaces, NewNamespace(n))
	}
	return namespaces
}
