package cluster

import (
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/entity"
	"code.storageos.net/storageos/c2-cli/pkg/id"
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

// Resource encapsulate a StorageOS cluster as a data type.
type Resource struct {
	ID id.Cluster `json:"id"`

	Licence *Licence `json:"licence,omitempty"`

	DisableTelemetry      bool `json:"disableTelemetry"`
	DisableCrashReporting bool `json:"disableCrashReporting"`
	DisableVersionCheck   bool `json:"disableVersionCheck"`

	LogLevel  LogLevel  `json:"logLevel"`
	LogFormat LogFormat `json:"logFormat"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Version   entity.Version `json:"version"`
}

// Licence encapsulates a StorageOS cluster product licence and the features
// included with it.
type Licence struct{}
