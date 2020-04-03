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

// UserNotFoundError indicates that the API could not find the StorageOS user
// specified.
type UserNotFoundError struct {
	msg string

	uid  id.User
	name string
}

// Error returns an error message indicating that the user with a given
// ID or name was not found, as configured.
func (e UserNotFoundError) Error() string {
	return e.msg
}

// NewUserNotFoundError returns a UserNotFoundError using details as the
// the error message. This can be used when provided an opaque but detailed
// error strings.
func NewUserNotFoundError(details string, uID id.User) UserNotFoundError {
	return UserNotFoundError{
		msg: details,
		uid: uID,
	}
}

// NewUserNameNotFoundError returns a UserNotFoundError for the user
// with name, constructing a user friendly message and storing the name inside
// the error.
func NewUserNameNotFoundError(name string) UserNotFoundError {
	return UserNotFoundError{
		msg:  fmt.Sprintf("user with name %v not found", name),
		name: name,
	}
}

// GetUserByName requests the details of a StorageOS user account with username
// and returns it to the caller.
func (c *Client) GetUserByName(ctx context.Context, username string) (*user.Resource, error) {
	list, err := c.Transport.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	for _, u := range list {
		if u.Username == username {
			return u, nil
		}
	}

	return nil, NewUserNameNotFoundError(username)
}

// GetListUsersByUID returns all the users with the ID listed in the uids parameter.
func (c *Client) GetListUsersByUID(ctx context.Context, uids []id.User) ([]*user.Resource, error) {
	list, err := c.Transport.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	toMap := make(map[id.User]*user.Resource)
	for _, u := range list {
		toMap[u.ID] = u
	}

	filtered := make([]*user.Resource, 0)
	for _, id := range uids {
		user, ok := toMap[id]
		if !ok {
			return nil, NewUserNotFoundError("user not found", id)
		}
		filtered = append(filtered, user)
	}

	return filtered, nil
}

// GetListUsersByUsername returns all the users with the username listed in the
// usernames parameter.
func (c *Client) GetListUsersByUsername(ctx context.Context, usernames []string) ([]*user.Resource, error) {
	list, err := c.Transport.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	toMap := make(map[string]*user.Resource)
	for _, u := range list {
		toMap[u.Username] = u
	}

	filtered := make([]*user.Resource, 0)
	for _, username := range usernames {
		user, ok := toMap[username]
		if !ok {
			return nil, NewUserNameNotFoundError(username)
		}
		filtered = append(filtered, user)
	}

	return filtered, nil
}
