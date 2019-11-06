// Package apiclient provides a type which implements an abstraction layer
// for consuming the storageos API programmatically.
package apiclient

import (
	"context"
	"fmt"
	"time"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/node"
	"code.storageos.net/storageos/c2-cli/pkg/volume"
)

// ConfigProvider describes the access to configuration values required by
// the apiclient package.
type ConfigProvider interface {
	APIEndpoint() (string, error)
	CommandTimeout() (time.Duration, error)
	Username() (string, error)
	Password() (string, error)
}

// Transport describes the set of methods required by an API client to use a
// given transport implementation provider.
type Transport interface {
	Authenticate(ctx context.Context, username, password string) error

	GetNode(ctx context.Context, uid id.Node) (*node.Resource, error)
	GetListNodes(ctx context.Context) ([]*node.Resource, error)
	GetVolume(context.Context, id.Namespace, id.Volume) (*volume.Resource, error)

	DescribeNode(ctx context.Context, uid id.Node) (*node.Resource, error)
	DescribeListNodes(ctx context.Context) ([]*node.Resource, error)
}

// Client provides a collection of methods for consumers to interact with the
// StorageOS API.
type Client struct {
	transport Transport
	// TODO: Config options?
	username string
	password string
	timeout  time.Duration
}

// GetNode requests the node resource corresponding to uid from the StorageOS
// API.
func (c *Client) GetNode(uid id.Node) (*node.Resource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Pre-authenticate request
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	n, err := c.transport.GetNode(ctx, uid)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// GetListNodes requests the list of all node resources which are members of
// the cluster, returning a list of node resources which have an ID present
// in argument list.
//
// If no IDs are specified, all are returned.
func (c *Client) GetListNodes(uids ...id.Node) ([]*node.Resource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Pre-authenticate request
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	nodes, err := c.transport.GetListNodes(ctx)
	if err != nil {
		return nil, err
	}

	if len(uids) == 0 {
		return nodes, nil
	}

	retrieved := map[id.Node]*node.Resource{}

	for _, n := range nodes {
		retrieved[n.ID] = n
	}

	filtered := make([]*node.Resource, len(uids))

	i := 0
	for _, id := range uids {
		n, ok := retrieved[id]
		if ok {
			filtered[i] = n
			i++
		} else {
			return nil, fmt.Errorf("node not found: %v", id)
		}
	}

	return nodes, nil
}

func (c *Client) GetVolume(namespace id.Namespace, uid id.Volume) (*volume.Resource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Pre-authenticate request
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	v, err := c.transport.GetVolume(ctx, namespace, uid)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// DescribeNode requests detailed information for the node resource which
// corresponds to uid from the StorageOS API.
func (c *Client) DescribeNode(uid id.Node) (*node.Resource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Pre-authenticate request
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	n, err := c.transport.DescribeNode(ctx, uid)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// DescribeListNodes requests a list containing detailed information on each
// node resource in the cluster.
//
// The returned list is filtered using uids so that it contains only those
// resources which have a matching ID. If none are specified, all are returned.
func (c *Client) DescribeListNodes(uids ...id.Node) ([]*node.Resource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Pre-authenticate request
	err := c.transport.Authenticate(ctx, c.username, c.password)
	if err != nil {
		return nil, err
	}

	nodes, err := c.transport.DescribeListNodes(ctx)
	if err != nil {
		return nil, err
	}

	if len(uids) == 0 {
		return nodes, nil
	}

	retrieved := map[id.Node]*node.Resource{}

	for _, n := range nodes {
		retrieved[n.ID] = n
	}

	filtered := make([]*node.Resource, len(uids))

	i := 0
	for _, id := range uids {
		n, ok := retrieved[id]
		if ok {
			filtered[i] = n
			i++
		} else {
			return nil, fmt.Errorf("node not found: %v", id)
		}
	}

	return nodes, nil
}

// New initialises a new Client using timeout for each operation, with
// transport providing the underlying implementation for encoding requests and
// decoding responses.
func New(transport Transport, config ConfigProvider) (*Client, error) {
	username, err := config.Username()
	if err != nil {
		return nil, err
	}

	password, err := config.Password()
	if err != nil {
		return nil, err
	}

	requestTimeout, err := config.CommandTimeout()
	if err != nil {
		return nil, err
	}

	return &Client{
		transport: transport,
		username:  username,
		password:  password,
		timeout:   requestTimeout,
	}, nil
}
