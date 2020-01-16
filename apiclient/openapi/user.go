package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/openapi"
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
