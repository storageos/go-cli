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
// ReplicaHealth The operational health of a volume replica deployment 
type ReplicaHealth string

// List of ReplicaHealth
const (
	REPLICAHEALTH_RECOVERING ReplicaHealth = "recovering"
	REPLICAHEALTH_PROVISIONING ReplicaHealth = "provisioning"
	REPLICAHEALTH_PROVISIONED ReplicaHealth = "provisioned"
	REPLICAHEALTH_SYNCING ReplicaHealth = "syncing"
	REPLICAHEALTH_READY ReplicaHealth = "ready"
	REPLICAHEALTH_DELETED ReplicaHealth = "deleted"
	REPLICAHEALTH_FAILED ReplicaHealth = "failed"
	REPLICAHEALTH_UNKNOWN ReplicaHealth = "unknown"
)