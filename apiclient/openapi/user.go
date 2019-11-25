package openapi

import (
	"context"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/openapi"
)

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
		return nil, mapOpenAPIError(err, resp)
	}

	return o.codec.decodeUser(model)
}
