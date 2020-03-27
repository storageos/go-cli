package output

import (
	"time"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

// Cluster defines a type that contains all the info needed to be outputted.
type Cluster struct {
	ID id.Cluster `json:"id" yaml:"id"`

	Licence *Licence `json:"licence" yaml:"licence"`

	DisableTelemetry      bool `json:"disableTelemetry" yaml:"disableTelemetry"`
	DisableCrashReporting bool `json:"disableCrashReporting" yaml:"disableCrashReporting"`
	DisableVersionCheck   bool `json:"disableVersionCheck" yaml:"disableVersionCheck"`

	LogLevel  cluster.LogLevel  `json:"logLevel" yaml:"logLevel"`
	LogFormat cluster.LogFormat `json:"logFormat" yaml:"logFormat"`

	CreatedAt time.Time       `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt" yaml:"updatedAt"`
	Version   version.Version `json:"version" yaml:"version"`
}

// Licence defines a type that contains all the info needed to be outputted.
type Licence struct {
	ClusterID            id.Cluster `json:"clusterID" yaml:"clusterID"`
	ExpiresAt            time.Time  `json:"expiresAt" yaml:"expiresAt"`
	ClusterCapacityBytes uint64     `json:"clusterCapacityBytes" yaml:"clusterCapacityBytes"`
	Kind                 string     `json:"kind" yaml:"kind"`
	CustomerName         string     `json:"customerName" yaml:"customerName"`
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
