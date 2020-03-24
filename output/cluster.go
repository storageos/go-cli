package output

import (
	"time"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

// Cluster defines a type that contains all the info needed to be outputted.
type Cluster struct {
	ID id.Cluster `json:"id"`

	Licence *Licence `json:"licence"`

	DisableTelemetry      bool `json:"disableTelemetry"`
	DisableCrashReporting bool `json:"disableCrashReporting"`
	DisableVersionCheck   bool `json:"disableVersionCheck"`

	LogLevel  cluster.LogLevel  `json:"logLevel"`
	LogFormat cluster.LogFormat `json:"logFormat"`

	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Version   version.Version `json:"version"`
}

// Licence defines a type that contains all the info needed to be outputted.
type Licence struct {
	ClusterID            id.Cluster `json:"clusterId"`
	ExpiresAt            time.Time  `json:"expiresAt"`
	ClusterCapacityBytes uint64     `json:"clusterCapacityBytes"`
	Kind                 string     `json:"kind"`
	CustomerName         string     `json:"customerName"`
}

// NewCluster returns a new Cluster object that contains all the info needed
// to be outputted.
func NewCluster(c *cluster.Resource) *Cluster {
	return &Cluster{
		ID:                    c.ID,
		Licence:               newLicence(c.Licence),
		DisableTelemetry:      c.DisableTelemetry,
		DisableCrashReporting: c.DisableCrashReporting,
		DisableVersionCheck:   c.DisableVersionCheck,
		LogLevel:              c.LogLevel,
		LogFormat:             c.LogFormat,
		CreatedAt:             c.CreatedAt,
		UpdatedAt:             c.UpdatedAt,
		Version:               c.Version,
	}
}

func newLicence(l *cluster.Licence) *Licence {
	return &Licence{
		ClusterID:            l.ClusterID,
		ExpiresAt:            l.ExpiresAt,
		ClusterCapacityBytes: l.ClusterCapacityBytes,
		Kind:                 l.Kind,
		CustomerName:         l.CustomerName,
	}
}
