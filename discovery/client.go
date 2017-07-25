package discovery

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/storageos/go-cli/types"
)

const (
	userAgent         = "go-storageosclient"
	unixProtocol      = "unix"
	namedPipeProtocol = "npipe"
	DefaultVersionStr = "1"
	DefaultVersion    = 1
	defaultNamespace  = "default"
)

var (
	// ErrInvalidEndpoint is returned when the endpoint is not a valid HTTP URL.
	ErrInvalidEndpoint = errors.New("invalid endpoint")

	// ErrInvalidVersion is returned when a versioned client was requested but no version specified.
	ErrInvalidVersion = errors.New("invalid version")

	// DefaultHost is the default API host
	DefaultHost = "https://discovery.storageos.cloud"
)

// APIVersion is an internal representation of a version of the Remote API.
type APIVersion int

// Client is the basic type of this package. It provides methods for
// interaction with the API.
type Client struct {
	SkipServerVersionCheck bool
	HTTPClient             *http.Client
	TLSConfig              *tls.Config
	endpoint               string
	username               string
	secret                 string
}

// NewClient returns a Client instance ready for communication with the given
// server endpoint. It will use the latest remote API version available in the
// server.
func NewClient(endpoint, username, secret string) (*Client, error) {
	if endpoint == "" {
		endpoint = DefaultHost
	}

	client := &Client{
		endpoint:   endpoint,
		HTTPClient: defaultClient(),
		secret:     secret,
		username:   username,
	}
	client.SkipServerVersionCheck = true
	return client, nil
}

// SetAuth sets the API username and secret to be used for all API requests.
// It should not be called concurrently with any other Client methods.
func (c *Client) SetAuth(username string, secret string) {
	if username != "" {
		c.username = username
	}
	if secret != "" {
		c.secret = secret
	}
}

// defaultClient returns a new http.Client with similar default values to
// http.Client, but with a non-shared Transport, idle connections disabled, and
// keepalives disabled.
func defaultClient() *http.Client {
	return &http.Client{
		Transport: defaultTransport(),
	}
}

// defaultPooledTransport returns a new http.Transport with similar default
// values to http.DefaultTransport. Do not use this for transient transports as
// it can leak file descriptors over time. Only use this for transports that
// will be re-used for the same host(s).
func defaultPooledTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 1,
	}
	return transport
}

// defaultTransport returns a new http.Transport with the same default values
// as http.DefaultTransport, but with idle connections and keepalives disabled.
func defaultTransport() *http.Transport {
	transport := defaultPooledTransport()
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

// ClusterCreate - creates new cluster token
func (c *Client) ClusterCreate(clusterName string, size int) (token string, err error) {
	path := c.endpoint + "/clusters"
	vals := url.Values{}
	vals.Set("size", fmt.Sprintf("%d", size))
	vals.Set("name", clusterName)
	path = fmt.Sprintf("%s?%s", path, vals.Encode())
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		return
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("discovery service responded with %d: %s", resp.StatusCode, string(body))
	}

	var cluster types.Cluster
	err = json.Unmarshal(body, &cluster)
	if err != nil {
		return
	}

	return cluster.ID, nil
}

// ClusterStatus - current cluster
func (c *Client) ClusterStatus(id string) (*types.Cluster, error) {
	req, err := http.NewRequest("GET", c.endpoint+"/clusters/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cluster types.Cluster
	err = json.Unmarshal(body, &cluster)
	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

// ClusterDelete - delete specified cluster
func (c *Client) ClusterDelete(id string) error {
	req, err := http.NewRequest("DELETE", c.endpoint+"/clusters/"+id, nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from discovery service: %d", resp.StatusCode)
	}

	return nil
}
