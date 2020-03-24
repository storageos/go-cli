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

// CreateNamespaceData struct for CreateNamespaceData
type CreateNamespaceData struct {
	// The name of the namespace shown in the CLI and UI
	Name string `json:"name,omitempty"`
	// A set of arbitrary key value labels to apply to the entity.
	Labels map[string]string `json:"labels,omitempty"`
}
