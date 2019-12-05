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

// Node struct for Node
type Node struct {
	// A unique identifier for a node. The format of this type is undefined and may change but the defined properties will not change.
	Id string `json:"id,omitempty"`
	// The hostname of the node. This value is set by the node each time it joins the StorageOS cluster.
	Name   string `json:"name,omitempty"`
	Health Health `json:"health,omitempty"`
	// Endpoint at which we operate our dataplane's dfs service. (used for IO operations) This value is set on startup by the corresponding environment variable (IO_ADVERTISE_ADDRESS)
	IoEndpoint string `json:"ioEndpoint,omitempty"`
	// Endpoint at which we operate our dataplane's supervisor service (used for sync). This value is set on startup by the corresponding environment variable (SUPERVISOR_ADVERTISE_ADDRESS)
	SupervisorEndpoint string `json:"supervisorEndpoint,omitempty"`
	// Endpoint at which we operate our health checking service. This value is set on startup by the corresponding environment variable (GOSSIP_ADVERTISE_ADDRESS)
	GossipEndpoint string `json:"gossipEndpoint,omitempty"`
	// Endpoint at which we operate our clustering GRPC API. This value is set on startup by the corresponding environment variable (INTERNAL_API_ADVERTISE_ADDRESS)
	ClusteringEndpoint string `json:"clusteringEndpoint,omitempty"`
	// A set of arbitrary key value labels to apply to the entity.
	Labels map[string]string `json:"labels,omitempty"`
	// The time the entity was created. This timestamp is set by the node that created the entity, and may not be correct if the node's local clock was skewed. This value is for the user's informative purposes only, and correctness is not required. String format is RFC3339.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// The time the entity was last updated. This timestamp is set by the node that last updated the entity, and may not be correct if the node's local clock was skewed. This value is for the user's informative purposes only, and correctness is not required. String format is RFC3339.
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// An opaque representation of an entity version at the time it was obtained from the API. All operations that mutate the entity must include this version field in the request unchanged. The format of this type is undefined and may change but the defined properties will not change.
	Version string `json:"version,omitempty"`
}
