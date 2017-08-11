package storageos

import (
	"context"
	"encoding/json"
	//"errors"
	//"fmt"
	"github.com/storageos/go-api/types"
	"net/http"
	//"net/url"
)

var (
	// HealthAPIPrefix is a partial path to the HTTP endpoint.
	HealthAPIPrefix = "health"
)

// CPHealth returns the health of the control plane server at a given url.
func (c *Client) CPHealth(hostname string) (*types.CPHealthStatus, error) {

	req, err := http.NewRequest("GET", "http://"+hostname+":5705/v1/"+HealthAPIPrefix, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	if c.username != "" && c.secret != "" {
		req.SetBasicAuth(c.username, c.secret)
	}

	resp, err := c.HTTPClient.Do(req.WithContext(context.Background()))
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
func (c *Client) DPHealth(hostname string) (*types.DPHealthStatus, error) {

	req, err := http.NewRequest("GET", "http://"+hostname+":8001/v1/"+HealthAPIPrefix, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	if c.username != "" && c.secret != "" {
		req.SetBasicAuth(c.username, c.secret)
	}

	resp, err := c.HTTPClient.Do(req.WithContext(context.Background()))
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
