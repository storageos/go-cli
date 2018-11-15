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

// SortableNodeType provides a container type to sort slices of the node type
// by name. This simplifies sorting as the type will manage the sort
// opperations on both types
type SortableNodeType struct {
	apiNode *[]*apiTypes.Node
	cliNode *[]*Node
}

func SortableAPIType(apiNode *[]*apiTypes.Node) *SortableNodeType {
	return &SortableNodeType{
		apiNode: apiNode,
	}
}

func SortableCLIType(cliNode *[]*Node) *SortableNodeType {
	return &SortableNodeType{
		cliNode: cliNode,
	}
}

func (s *SortableNodeType) SortByName() error {
	switch {
	case s.apiNode != nil && s.cliNode == nil:
		break
	case s.cliNode != nil && s.apiNode == nil:
		break
	default:
		return errors.New("more than one type of node slice used")
	}

	sort.Sort(s)
	return nil
}

func (s *SortableNodeType) Len() int {
	switch {
	case s.apiNode != nil:
		return len(*s.apiNode)
	case s.cliNode != nil:
		return len(*s.cliNode)
	default:
		return 0
	}
}

func (s *SortableNodeType) Swap(i, j int) {
	if s.apiNode != nil {
		apiSlice := *s.apiNode
		apiSlice[i], apiSlice[j] = apiSlice[j], apiSlice[i]
	}

	if s.cliNode != nil {
		cliSlice := *s.cliNode
		cliSlice[i], cliSlice[j] = cliSlice[j], cliSlice[i]
	}
}

func (s *SortableNodeType) Less(i, j int) bool {
	var name1, name2 string
	if s.apiNode != nil {
		apiSlice := *s.apiNode
		name1, name2 = trimCommonPrefix(apiSlice[i].Name, apiSlice[j].Name)
	}
	if s.cliNode != nil {
		cliSlice := *s.cliNode
		name1, name2 = trimCommonPrefix(cliSlice[i].Name, cliSlice[j].Name)
	}

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

//// NodeByName sorts node list by hostname.
//// In the case that the hostnames are formed of a common prefix, followed by a
//// trailing number they will be sorted by this number rather than just
//// lexicographically.
//type NodeByName []*Node
//
//func (n NodeByName) Len() int      { return len(n) }
//func (n NodeByName) Swap(i, j int) { n[i], n[j] = n[j], n[i] }
//func (n NodeByName) Less(i, j int) bool {
//	name1, name2 := trimCommonPrefix(n[i].Name, n[j].Name)
//
//	// Are the postfixes both numerical, if so sort as integers
//	n1, err1 := strconv.Atoi(name1)
//	n2, err2 := strconv.Atoi(name2)
//	if err1 == nil && err2 == nil {
//		return n1 < n2
//	}
//
//	// Postfixes don't appear to be numerical, sort them lexicographically
//	return name1 < name2
//}
//
//func trimCommonPrefix(a, b string) (string, string) {
//	if a == b {
//		return "", ""
//	}
//
//	for i, r := range a {
//		if r != []rune(b)[i] {
//			a, b = a[i:], b[i:]
//			break
//		}
//	}
//
//	return a, b
//}
