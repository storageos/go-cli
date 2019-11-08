package openapi

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"code.storageos.net/storageos/openapi"
)

type OpenAPI struct {
	mu *sync.RWMutex

	client *openapi.APIClient
	codec  codec
}

func (o *OpenAPI) Authenticate(ctx context.Context, username, password string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	_, resp, err := o.client.DefaultApi.AuthenticateUser(
		ctx,
		openapi.AuthUserData{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	token := strings.TrimPrefix(resp.Header.Get("Authorization"), "Bearer ")
	o.client.GetConfig().AddDefaultHeader("Authorization", token)

	return nil
}

func NewOpenAPI(apiEndpoint, userAgent string) *OpenAPI {
	// Init the OpenAPI client
	cfg := &openapi.Configuration{
		BasePath:      "v2",
		DefaultHeader: map[string]string{},
		Host:          apiEndpoint,
		Scheme:        "http",
		UserAgent:     userAgent,
	}

	client := openapi.NewAPIClient(cfg)

	return &OpenAPI{
		mu: &sync.RWMutex{},

		client: client,
		codec:  codec{},
	}
}
