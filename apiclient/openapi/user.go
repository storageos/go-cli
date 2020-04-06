package openapi

import (
	"context"

	"github.com/antihax/optional"

	"code.storageos.net/storageos/openapi"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/user"
)

// CreateUser requests the creation of a new StorageOS user account using the
// provided parameters.
func (o *OpenAPI) CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error) {

	gs := make([]string, len(groups))
	for i, g := range groups {
		gs[i] = g.String()
	}

	createData := openapi.CreateUserData{
		Username: username,
		Password: password,
		IsAdmin:  withAdmin,
		Groups:   &gs,
	}

	model, resp, err := o.client.DefaultApi.CreateUser(ctx, createData)
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case badRequestError:
			return nil, apiclient.NewInvalidUserCreationError(v.msg)
		case conflictError:
			return nil, apiclient.NewUserExistsError(username)
		default:
			return nil, v
		}
	}

	return o.codec.decodeUser(model)
}

// GetUser requests all the details of a StorageOS user account, starting from
// its id.
func (o *OpenAPI) GetUser(ctx context.Context, uID id.User) (*user.Resource, error) {
	model, resp, err := o.client.DefaultApi.GetUser(ctx, uID.String())
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case notFoundError:
			return nil, apiclient.NewUserNotFoundError(v.msg, uID)
		default:
			return nil, v
		}
	}

	return o.codec.decodeUser(model)
}

// ListUsers requests the list of all users accounts.
func (o *OpenAPI) ListUsers(ctx context.Context) ([]*user.Resource, error) {
	list, resp, err := o.client.DefaultApi.ListUsers(ctx)
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		default:
			return nil, v
		}
	}

	users := make([]*user.Resource, 0, len(list))

	for _, u := range list {
		user, err := o.codec.decodeUser(u)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// DeleteUser makes a delete request for a user given its ID.
//
// The behaviour of the operation is dictated by params:
//
//
// 	Version constraints:
// 	- If params is nil or params.CASVersion is empty then the delete request is
// 	unconditional
// 	- If params.CASVersion is set, the request is conditional upon it matching
// 	the volume entity's version as seen by the server.
func (o *OpenAPI) DeleteUser(ctx context.Context, uid id.User, params *apiclient.DeleteUserRequestParams) error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	var casVersion string
	var ignoreVersion optional.Bool = optional.NewBool(true)

	if params != nil {
		if params.CASVersion.String() != "" {
			ignoreVersion = optional.NewBool(false)
			casVersion = params.CASVersion.String()
		}
	}

	resp, err := o.client.DefaultApi.DeleteUser(
		ctx,
		uid.String(),
		casVersion,
		&openapi.DeleteUserOpts{
			IgnoreVersion: ignoreVersion,
		},
	)
	if err != nil {
		switch v := mapOpenAPIError(err, resp).(type) {
		case notFoundError:
			return apiclient.NewUserNotFoundError(v.msg, uid)
		default:
			return v
		}
	}

	return nil
}
