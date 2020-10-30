/*
 * StorageOS API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 2.3.0-alpha
 * Contact: info@storageos.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// CreateUserData struct for CreateUserData
type CreateUserData struct {
	Username string `json:"username"`
	// If not present, the existing password is not changed
	Password string `json:"password"`
	// If true, this user is an administrator of the cluster. Administrators bypass the usual authentication checks and are granted access to all resources. Some actions (such as adding a new user) can only be performed by an administrator.
	IsAdmin bool `json:"isAdmin,omitempty"`
	// Defines a set of policy group IDs this user is a member of. Policy groups can be used to logically group users and apply authorisation  policies to all members.
	Groups *[]string `json:"groups,omitempty"`
}
