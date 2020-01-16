package apiclient

import (
	"context"
	"fmt"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/user"
)

// UserExistsError is returned when a user creation request is sent to the
// StorageOS API for an already taken username.
type UserExistsError struct {
	username string
}

// Error returns an error message indicating that a username is already in use.
func (e UserExistsError) Error() string {
	return fmt.Sprintf("another user with username %v already exists", e.username)
}

// NewUserExistsError returns an error indicating that a user already exists
// for username.
func NewUserExistsError(username string) UserExistsError {
	return UserExistsError{
		username: username,
	}
}

// InvalidUserCreationError is returned when an user creation request sent to
// the StorageOS API is invalid.
type InvalidUserCreationError struct {
	details string
}

// Error returns an error message indicating that a user creation request
// made to the StorageOS API is invalid, including details if available.
func (e InvalidUserCreationError) Error() string {
	msg := "user creation request is invalid"
	if e.details != "" {
		msg = fmt.Sprintf("%v: %v", msg, e.details)
	}
	return msg
}

// NewInvalidUserCreationError returns an InvalidUserCreationError, using
// details to provide information about what must be corrected.
func NewInvalidUserCreationError(details string) InvalidUserCreationError {
	return InvalidUserCreationError{
		details: details,
	}
}

// CreateUser requests the creation of a new StorageOS user account from the
// provided fields. If successful the created resource for the user account
// is returned to the caller.
func (c *Client) CreateUser(
	ctx context.Context,
	username, password string,
	withAdmin bool,
	groups ...id.PolicyGroup,
) (*user.Resource, error) {
	_, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transport.CreateUser(
		ctx,
		username,
		password,
		withAdmin,
		groups...,
	)
}
