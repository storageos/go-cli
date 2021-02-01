/*
 * StorageOS API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 2.4.0-alpha
 * Contact: info@storageos.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// SetFailureModeRequest struct for SetFailureModeRequest
type SetFailureModeRequest struct {
	// The minimum number of replicas required to be online and receiving writes in order for the volume to remain read-writable. This value replaces any previously set failure threshold or intent-based failure mode.
	FailureThreshold uint64 `json:"failureThreshold,omitempty"`
	// An opaque representation of an entity version at the time it was obtained from the API. All operations that mutate the entity must include this version field in the request unchanged. The format of this type is undefined and may change but the defined properties will not change.
	Version string            `json:"version,omitempty"`
	Mode    FailureModeIntent `json:"mode,omitempty"`
}
