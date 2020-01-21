package cluster

import (
	"fmt"
	"time"

	"github.com/alecthomas/units"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

// LogLevel is a typed wrapper around a cluster's log level configuration.
type LogLevel string

// LogLevelFromString wraps level as a LogLevel.
func LogLevelFromString(level string) LogLevel {
	return LogLevel(level)
}

// String returns the string representation of l.
func (l LogLevel) String() string {
	return string(l)
}

// LogFormat is a typed wrapper around a cluster's log entry format
// configuration.
type LogFormat string

// LogFormatFromString wraps format as a LogFormat.
func LogFormatFromString(format string) LogFormat {
	return LogFormat(format)
}

// String returns the string representation of f.
func (f LogFormat) String() string {
	return string(f)
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

func (l *Licence) String() string {
	return fmt.Sprintf(`Cluster ID: %v
Expires at: %v
Cluster capacity: %v
Kind: %v
Customer name: %v
`,
		l.ClusterID,
		l.ExpiresAt.Format(time.RFC3339),
		units.Base2Bytes(l.ClusterCapacityBytes).String(),
		l.Kind,
		l.CustomerName)
}
