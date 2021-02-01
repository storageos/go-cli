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
import (
	"time"
)
// Volume struct for Volume
type Volume struct {
	// A unique identifier for a volume. The format of this type is undefined and may change but the defined properties will not change. 
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	AttachedOn string `json:"attachedOn,omitempty"`
	Nfs NfsConfig `json:"nfs,omitempty"`
	NamespaceID string `json:"namespaceID,omitempty"`
	// A set of arbitrary key value labels to apply to the entity. 
	Labels map[string]string `json:"labels,omitempty"`
	FsType FsType `json:"fsType,omitempty"`
	AttachmentType AttachType `json:"attachmentType,omitempty"`
	Master MasterDeploymentInfo `json:"master,omitempty"`
	Replicas *[]ReplicaDeploymentInfo `json:"replicas,omitempty"`
	// A volume's size in bytes 
	SizeBytes uint64 `json:"sizeBytes,omitempty"`
	// The time the entity was created. This timestamp is set by the node that created the entity, and may not be correct if the node's local clock was skewed. This value is for the user's informative purposes only, and correctness is not required. String format is RFC3339. 
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// The time the entity was last updated. This timestamp is set by the node that last updated the entity, and may not be correct if the node's local clock was skewed. This value is for the user's informative purposes only, and correctness is not required. String format is RFC3339. 
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// An opaque representation of an entity version at the time it was obtained from the API. All operations that mutate the entity must include this version field in the request unchanged. The format of this type is undefined and may change but the defined properties will not change. 
	Version string `json:"version,omitempty"`
}
