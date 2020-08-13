package openapi

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/openapi"
)

// GetDiagnostics makes a request to the StorageOS API for a cluster diagnostic
// bundle to be generated and returned to the client.
//
// Because the OpenAPI code generator produces broken code for this method,
// we source the target path, authorization token and http client from it but
// handle the response ourselves.
func (o *OpenAPI) GetDiagnostics(ctx context.Context) (io.ReadCloser, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	// Get the appropriate config settings from the openapi client
	token := o.client.GetConfig().DefaultHeader["Authorization"]
	targetEndpoint := o.client.GetConfig().Scheme + "://" + o.client.GetConfig().Host + "/" + o.client.GetConfig().BasePath + "/diagnostics"
	client := o.client.GetConfig().HTTPClient

	// Construct the request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header["Authorization"] = []string{token}
	req.Header["Accept"] = []string{"application/gzip", "application/json"}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		// Carry on.
	case http.StatusBadGateway:
		// Check if the response content-type indicates a partial bundle. That
		// is, it has a gzip content type.
		for _, value := range resp.Header["Content-Type"] {
			if value == "application/gzip" {
				return nil, apiclient.NewIncompleteDiagnosticsError(resp.Body)
			}
		}

		// If not, use the normal error handling code.
		fallthrough
	default:
		defer resp.Body.Close()
		// Try to read the response body and unmarshal it into an openapi.Error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var errModel openapi.Error
		err = json.Unmarshal(body, &errModel)
		if err != nil {
			return nil, err
		}

		// Construct an openAPIError from it and hand off to the
		// OpenAPI error mapping code.
		return nil, mapOpenAPIError(
			newOpenAPIError(errModel),
			resp,
		)
	}

	return resp.Body, nil
}
