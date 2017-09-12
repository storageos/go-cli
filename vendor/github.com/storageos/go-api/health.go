package storageos

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/storageos/go-api/types"
)

var (
	// HealthAPIPrefix is a partial path to the HTTP endpoint.
	HealthAPIPrefix = "health"
)

// CPHealth returns the health of the control plane server at a given url.
func (c *Client) CPHealth(ctx context.Context, hostname string) (*types.CPHealthStatus, error) {

	req, err := http.NewRequest("GET", "http://"+hostname+":5705/v1/"+HealthAPIPrefix, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	if c.username != "" && c.secret != "" {
		req.SetBasicAuth(c.username, c.secret)
	}

	resp, err := c.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var status *types.CPHealthStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}

	return status, nil
}

// DPHealth returns the health of the data plane server at a given url.
func (c *Client) DPHealth(ctx context.Context, hostname string) (*types.DPHealthStatus, error) {

	req, err := http.NewRequest("GET", "http://"+hostname+":8001/v1/"+HealthAPIPrefix, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	if c.username != "" && c.secret != "" {
		req.SetBasicAuth(c.username, c.secret)
	}

	resp, err := c.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var status *types.DPHealthStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}

	return status, nil
}
