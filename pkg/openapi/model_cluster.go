/*
 * StorageOS API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 2.0.0
 * Contact: info@storageos.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"time"
)

// Cluster struct for Cluster
type Cluster struct {
	// A unique identifier for a cluster. The format of this type is undefined and may change but the defined properties will not change.
	Id      string  `json:"id,omitempty"`
	Licence Licence `json:"licence,omitempty"`
	// Disables collection of telemetry data across the cluster.
	DisableTelemetry bool `json:"disableTelemetry,omitempty"`
	// Disables collection of reports for any fatal crashes across the cluster.
	DisableCrashReporting bool `json:"disableCrashReporting,omitempty"`
	// Disables the mechanism responsible for checking if there is an updated version of StorageOS available for installation.
	DisableVersionCheck bool      `json:"disableVersionCheck,omitempty"`
	LogLevel            LogLevel  `json:"logLevel,omitempty"`
	LogFormat           LogFormat `json:"logFormat,omitempty"`
	// The time the entity was created. This timestamp is set by the node that created the entity, and may not be correct if the node's local clock was skewed. This value is for the user's informative purposes only, and correctness is not required. String format is RFC3339.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// The time the entity was last updated. This timestamp is set by the node that last updated the entity, and may not be correct if the node's local clock was skewed. This value is for the user's informative purposes only, and correctness is not required. String format is RFC3339.
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// An opaque representation of an entity version at the time it was obtained from the API. All operations that mutate the entity must include this version field in the request unchanged. The format of this type is undefined and may change but the defined properties will not change.
	Version string `json:"version,omitempty"`
}
