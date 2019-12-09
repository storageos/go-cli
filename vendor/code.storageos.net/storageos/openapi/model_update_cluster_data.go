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
// UpdateClusterData struct for UpdateClusterData
type UpdateClusterData struct {
	// A StorageOS product licence key, used to register a cluster. The format of this type is opaque and may change. If this field is empty we assume that the called does not want to update the cluster's licence, and do not perform any operation. 
	LicenceKey string `json:"licenceKey,omitempty"`
	// Disables collection of telemetry data across the cluster. 
	DisableTelemetry bool `json:"disableTelemetry,omitempty"`
	// Disables collection of reports for any fatal crashes across the cluster. 
	DisableCrashReporting bool `json:"disableCrashReporting,omitempty"`
	// Disables the mechanism responsible for checking if there is an updated version of StorageOS available for installation. 
	DisableVersionCheck bool `json:"disableVersionCheck,omitempty"`
	LogLevel LogLevel `json:"logLevel,omitempty"`
	LogFormat LogFormat `json:"logFormat,omitempty"`
	// An opaque representation of an entity version at the time it was obtained from the API. All operations that mutate the entity must include this version field in the request unchanged. The format of this type is undefined and may change but the defined properties will not change. 
	Version string `json:"version,omitempty"`
}
