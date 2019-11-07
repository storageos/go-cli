package openapi

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"code.storageos.net/storageos/c2-cli/pkg/cluster"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/namespace"
	"code.storageos.net/storageos/c2-cli/pkg/node"
	"code.storageos.net/storageos/c2-cli/pkg/volume"

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

func (o *OpenAPI) GetCluster(ctx context.Context) (*cluster.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, _, err := o.client.DefaultApi.GetCluster(ctx)
	if err != nil {
		return nil, err
	}

	return o.codec.decodeCluster(model)
}

func (o *OpenAPI) GetNode(ctx context.Context, uid id.Node) (*node.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, _, err := o.client.DefaultApi.GetNode(ctx, uid.String())
	if err != nil {
		// TODO: Maybe do the error mapping at the transport level?
		// → if so change below as well.
		// → Error mapping could use the resp object to be a bit
		// intelligent?
		return nil, err
	}

	n, err := o.codec.decodeNode(model)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (o *OpenAPI) GetVolume(ctx context.Context, namespace id.Namespace, uid id.Volume) (*volume.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	model, _, err := o.client.DefaultApi.GetVolume(ctx, namespace.String(), uid.String())
	if err != nil {
		return nil, err
	}

	v, err := o.codec.decodeVolume(model)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (o *OpenAPI) ListNodes(ctx context.Context) ([]*node.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	models, _, err := o.client.DefaultApi.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	nodes := make([]*node.Resource, len(models))
	for i, m := range models {

		// → If we error here then there's an incompatibility somewhere so
		// aborting is probably a good shout.
		n, err := o.codec.decodeNode(m)
		if err != nil {
			return nil, err
		}

		nodes[i] = n
	}

	return nodes, nil
}

func (o *OpenAPI) ListVolumes(ctx context.Context, namespace id.Namespace) ([]*volume.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	models, _, err := o.client.DefaultApi.ListVolumes(ctx, namespace.String())
	if err != nil {
		return nil, err
	}

	volumes := make([]*volume.Resource, len(models))
	for i, m := range models {
		v, err := o.codec.decodeVolume(m)
		if err != nil {
			return nil, err
		}

		volumes[i] = v
	}

	return volumes, nil
}

func (o *OpenAPI) ListNamespaces(ctx context.Context) ([]*namespace.Resource, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	models, _, err := o.client.DefaultApi.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	namespaces := make([]*namespace.Resource, len(models))
	for i, m := range models {
		ns, err := o.codec.decodeNamespace(m)
		if err != nil {
			return nil, err
		}

		namespaces[i] = ns
	}

	return namespaces, nil
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
