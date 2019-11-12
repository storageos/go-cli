package openapi

import (
	"context"
	"errors"
	"strings"
	"sync"

	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/openapi"
)

type ConfigProvider interface {
	APIEndpoints() ([]string, error)
}

// OpenAPI provides functionality for consuming the REST API exposed by
// StorageOS, implemented with a client generated from the OpenAPI
// specification.
//
// The codec stored on the type is responsible for translating the returned
// OpenAPI models into the internal types which are returned to consumers of
// the type.
type OpenAPI struct {
	mu *sync.RWMutex

	config ConfigProvider
	client *openapi.APIClient
	codec  codec
}

func (o *OpenAPI) Authenticate(ctx context.Context, username, password string) (*user.Resource, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	model, resp, err := o.client.DefaultApi.AuthenticateUser(
		ctx,
		openapi.AuthUserData{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		if resp != nil {
			return nil, mapResponseToError(resp)
		} else {
			return nil, err
		}
	}

	token := strings.TrimPrefix(resp.Header.Get("Authorization"), "Bearer ")
	o.client.GetConfig().AddDefaultHeader("Authorization", token)

	return o.codec.decodeUser(model)
}

func NewOpenAPI(config ConfigProvider, userAgent string) (*OpenAPI, error) {
	hosts, err := config.APIEndpoints()
	if err != nil || len(hosts) == 0 {
		return nil, errors.New("unable to determine target host")
	}

	// Create the OpenAPI client configuration
	// and initialise.
	apiCfg := &openapi.Configuration{
		BasePath:      "v2",
		DefaultHeader: map[string]string{},
		// TODO(CP-3924): For now the CLI supports only sending requests to the
		// first host provided. There should be a way to utilise multiple
		// hosts.
		Host: hosts[0],
		// TODO(CP-3913): Support TLS.
		Scheme:    "http",
		UserAgent: userAgent,
	}

	client := openapi.NewAPIClient(apiCfg)

	return &OpenAPI{
		mu: &sync.RWMutex{},

		client: client,
		codec:  codec{},
	}, nil
}
