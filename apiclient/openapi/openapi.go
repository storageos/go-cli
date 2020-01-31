package openapi

import (
	"context"
	"errors"
	"strings"
	"sync"

	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/openapi"
)

// ConfigProvider abstracts the functionality required by the OpenAPI transport
// implementation for client configuration.
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

// Authenticate attempts to authenticate against the target API using username
// and password. If successful, o's underlying OpenAPI client will use the
// returned token in the Authorization header for future operations.
//
// The returned *user.Resource corresponds to the authenticated user.
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
		return nil, mapOpenAPIError(err, resp)
	}

	token := strings.TrimPrefix(resp.Header.Get("Authorization"), "Bearer ")
	o.client.GetConfig().AddDefaultHeader("Authorization", token)

	return o.codec.decodeUser(model)
}

// NewOpenAPI initialises a new OpenAPI transport using config to source the
// target host endpoints and userAgent as the HTTP user agent string.
func NewOpenAPI(config ConfigProvider, userAgent string) (*OpenAPI, error) {
	hosts, err := config.APIEndpoints()
	if err != nil || len(hosts) == 0 {
		return nil, errors.New("unable to determine target host")
	}

	// TODO(CP-3924): This is not good - fix how we get API endpoints from the config.
	// This should be done as part of the work in supporting multiple endpoints.
	parts := strings.Split(hosts[0], "://")
	switch len(parts) {
	case 1:
		parts = []string{"http", parts[0]}
	case 2:
	default:
		return nil, errors.New("unable to parse target host")
	}

	// Create the OpenAPI client configuration
	// and initialise.
	apiCfg := &openapi.Configuration{
		BasePath:      "v2",
		DefaultHeader: map[string]string{},
		// TODO(CP-3924): For now the CLI supports only sending requests to the
		// first host provided. There should be a way to utilise multiple
		// hosts.
		Host: parts[1],
		// TODO(CP-3913): Support TLS.
		Scheme:    parts[0],
		UserAgent: userAgent,
	}

	client := openapi.NewAPIClient(apiCfg)

	return &OpenAPI{
		mu: &sync.RWMutex{},

		config: config,
		client: client,
		codec:  codec{},
	}, nil
}
