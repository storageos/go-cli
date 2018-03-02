package types

import "context"

// PoolUpdateOptions are available parameters for updating existing controllers.
type PoolUpdateOptions struct {

	// Pool unique ID.
	// Read Only: true
	ID string `json:"id"`

	// Name of the storage pool.
	// Read Only: true
	Name string `json:"name"`

	// Description of the storage pool.
	Description string `json:"description"`

	// Default determines whether the pool is the default when a volume is
	// provisioned without a pool specified. There can only be one default
	// pool.
	Default bool `string:"default"`

	// DefaultDriver specifies the storage driver to use by default if there
	// are multiple drivers in the storage pool.
	DefaultDriver string `json:"defaultDriver"`

	// ControllerNames is the list of controllers which are participating in
	// the storage pool.
	ControllerNames []string `json:"controllerNames"`

	// DriverNames is the list of backend storage drivers that are available
	// in the storage pool.
	DriverNames []string `json:"driverNames"`

	// Active status of the pool
	Active bool `json:"active"`

	// Labels that describe the pool.
	Labels map[string]string `json:"labels"`

	// Context can be set with a timeout or can be used to cancel a request.
	Context context.Context `json:"-"`
}
