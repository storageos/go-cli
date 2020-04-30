package output

import (
	"time"

	"code.storageos.net/storageos/c2-cli/licence"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

// Licence defines a type that contains all the info needed to be outputted.
type Licence struct {
	ClusterID            id.Cluster `json:"clusterID" yaml:"clusterID"`
	ExpiresAt            time.Time  `json:"expiresAt" yaml:"expiresAt"`
	ClusterCapacityBytes uint64     `json:"clusterCapacityBytes" yaml:"clusterCapacityBytes"`
	UsedBytes            uint64     `json:"usedBytes" yaml:"usedBytes"`
	Kind                 string     `json:"kind" yaml:"kind"`
	CustomerName         string     `json:"customerName" yaml:"customerName"`
}

// NewLicence returns a new licence object that contains all the info needed
// to be outputted.
func NewLicence(l *licence.Resource) *Licence {
	return &Licence{
		ClusterID:            l.ClusterID,
		ExpiresAt:            l.ExpiresAt,
		ClusterCapacityBytes: l.ClusterCapacityBytes,
		UsedBytes:            l.UsedBytes,
		Kind:                 l.Kind,
		CustomerName:         l.CustomerName,
	}
}
