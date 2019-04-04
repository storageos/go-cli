package types

import (
	"errors"
	"sort"
	"strconv"
	"time"

	apiTypes "github.com/storageos/go-api/types"
)

// ClusterCreateOps - optional fields when creating cluster
type ClusterCreateOps struct {
	AccountID string
	// optional value when to expire cluster
	TTL  int64
	Name string
	Size int
}

// Cluster is a representation of a storageos cluster as used by a
// storageos discovery service.
type Cluster struct {
	// cluster ID used for joining or getting cluster status
	ID string `json:"id,omitempty"`

	// cluster size, defaults to 3
	Size int `json:"size,omitempty"`

	Name string `json:"name,omitempty"`

	// optional account ID
	AccountID string `json:"accountID,omitempty"`

	// nodes participating in cluster
	Nodes []*Node `json:"nodes,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// NodeHealth is a containter type for holding Node health information
type NodeHealth struct {
	CP *apiTypes.CPHealthStatus
	DP *apiTypes.DPHealthStatus
}

// Node is an encapsulation of a storageos cluster node and its health state.
type Node struct {
	ID               string `json:"id,omitempty"` // node/controller UUID
	Name             string `json:"name,omitempty"`
	AdvertiseAddress string `json:"advertiseAddress,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`

	Health NodeHealth
}

type nodeSortBy int

// Pre-defined sorting methods
const (
	ByNodeName nodeSortBy = iota
)

// SortAPINodes sorts the set of nodes by the provided scheme
func SortAPINodes(by nodeSortBy, nodes []*apiTypes.Node) error {
	lessfunc, err := apiNodeSortFunc(by, nodes)
	if err != nil {
		return err
	}
	sort.Slice(nodes, lessfunc)
	return nil
}

// SortCLINodes sorts the set of nodes by the provided scheme
func SortCLINodes(by nodeSortBy, nodes []*Node) error {
	lessfunc, err := cliNodeSortFunc(by, nodes)
	if err != nil {
		return err
	}
	sort.Slice(nodes, lessfunc)
	return nil
}

func apiNodeSortFunc(sortBy nodeSortBy, nodes []*apiTypes.Node) (func(i, j int) bool, error) {
	switch sortBy {
	case ByNodeName:
		return func(i, j int) bool {
			return HumanisedStringLess(nodes[i].Name, nodes[j].Name)
		}, nil

	default:
		return nil, errors.New("sort method not implemented")
	}
}

func cliNodeSortFunc(sortBy nodeSortBy, nodes []*Node) (func(i, j int) bool, error) {
	switch sortBy {
	case ByNodeName:
		return func(i, j int) bool {
			return HumanisedStringLess(nodes[i].Name, nodes[j].Name)
		}, nil

	default:
		return nil, errors.New("sort method not implemented")
	}
}

// HumanisedStringLess is a string compare function, useable for sorting that
// attempts to detect expected humanised sorting e.g. hostnames with numeric
// postfixes.
//
// This function (for now) is quite basic, but could support more edge-cases as
// they arise.
func HumanisedStringLess(i, j string) bool {
	name1, name2 := trimCommonPrefix(i, j)

	// Are the postfixes both numerical, if so sort as integers
	n1, err1 := strconv.Atoi(name1)
	n2, err2 := strconv.Atoi(name2)
	if err1 == nil && err2 == nil {
		return n1 < n2
	}

	// Postfixes don't appear to be numerical, sort them lexicographically
	return name1 < name2
}

func trimCommonPrefix(a, b string) (string, string) {
	if a == b {
		return "", ""
	}

	for i, r := range a {
		if r != []rune(b)[i] {
			a, b = a[i:], b[i:]
			break
		}
	}

	return a, b
}
