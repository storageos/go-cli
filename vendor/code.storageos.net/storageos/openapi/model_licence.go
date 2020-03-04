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
// Licence A representation of a cluster's licence properties 
type Licence struct {
	// A unique identifier for a cluster. The format of this type is undefined and may change but the defined properties will not change. 
	ClusterID string `json:"clusterID,omitempty"`
	// The time after which a licence will no longer be valid This timestamp is set when the licence is created. String format is RFC3339. 
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
	// The allowed provisioning capacity in bytes This value if for the cluster, if provisioning a volume brings the cluster's total provisioned capacity above it the request will fail 
	ClusterCapacityBytes uint64 `json:"clusterCapacityBytes,omitempty"`
	// Denotes which category the licence belongs to 
	Kind string `json:"kind,omitempty"`
	// A user friendly reference to the customer 
	CustomerName string `json:"customerName,omitempty"`
}
