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

// CapacityStats struct for CapacityStats
type CapacityStats struct {
	// Total bytes in the filesystem
	Total uint64 `json:"total,omitempty"`
	// Free bytes in the filesystem available to root user
	Free uint64 `json:"free,omitempty"`
	// Byte value available to an unprivileged user
	Available uint64 `json:"available,omitempty"`
}
