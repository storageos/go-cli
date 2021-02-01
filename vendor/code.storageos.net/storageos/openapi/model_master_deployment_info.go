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
// MasterDeploymentInfo struct for MasterDeploymentInfo
type MasterDeploymentInfo struct {
	// A unique identifier for a volume deployment. The format of this type is undefined and may change but the defined properties will not change. 
	Id string `json:"id,omitempty"`
	// A unique identifier for a node. The format of this type is undefined and may change but the defined properties will not change. 
	NodeID string `json:"nodeID,omitempty"`
	// Indicates if the volume deployment is eligible for promotion 
	Promotable bool `json:"promotable,omitempty"`
	Health MasterHealth `json:"health,omitempty"`
}
