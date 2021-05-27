/*
 * StorageOS API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 2.4.0
 * Contact: info@storageos.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// AttachType The attachment type of a volume. \"host\" indicates that the volume is consumed by the node it is attached to.
type AttachType string

// List of AttachType
const (
	ATTACHTYPE_UNKNOWN  AttachType = "unknown"
	ATTACHTYPE_DETACHED AttachType = "detached"
	ATTACHTYPE_NFS      AttachType = "nfs"
	ATTACHTYPE_HOST     AttachType = "host"
)
