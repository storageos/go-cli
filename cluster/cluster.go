package cluster

import (
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

type LogLevel string

func LogLevelFromString(level string) LogLevel {
	return LogLevel(level)
}

func (l LogLevel) String() string {
	return string(l)
}

type LogFormat string

func LogFormatFromString(format string) LogFormat {
	return LogFormat(format)
}

func (l LogFormat) String() string {
	return string(l)
}

// Resource encapsulate a StorageOS cluster api resource as a data type.
type Resource struct {
	ID id.Cluster `json:"id"`

	Licence *Licence `json:"licence"`

	DisableTelemetry      bool `json:"disableTelemetry"`
	DisableCrashReporting bool `json:"disableCrashReporting"`
	DisableVersionCheck   bool `json:"disableVersionCheck"`

	LogLevel  LogLevel  `json:"logLevel"`
	LogFormat LogFormat `json:"logFormat"`

	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Version   version.Version `json:"version"`
}

// Licence describes a StorageOS product licence and the features included with
// it.
type Licence struct {
	ClusterID            id.Cluster `json:"clusterId"`
	ExpiresAt            time.Time  `json:"expiresAt"`
	ClusterCapacityBytes uint64     `json:"clusterCapacityBytes"`
	Kind                 string     `json:"kind"`
	CustomerName         string     `json:"customerName"`
}
