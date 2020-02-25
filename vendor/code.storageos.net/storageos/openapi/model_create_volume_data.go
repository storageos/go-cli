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
// CreateVolumeData struct for CreateVolumeData
type CreateVolumeData struct {
	// A unique identifier for a namespace. The format of this type is undefined and may change but the defined properties will not change.. 
	NamespaceID string `json:"namespaceID"`
	// A set of arbitrary key value labels to apply to the entity. 
	Labels map[string]string `json:"labels,omitempty"`
	// The name of the volume shown in the CLI and UI 
	Name string `json:"name"`
	FsType FsType `json:"fsType"`
	Description string `json:"description,omitempty"`
	// A volume's size in bytes 
	SizeBytes uint64 `json:"sizeBytes"`
}
