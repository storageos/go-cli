package openapi

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/diagnostics"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/openapi"
)

var (
	// errExtractingFilename is an error indicating a filename was not extracted
	// from the response header.
	errExtractingFilename = errors.New("failed to extract filename from response header")
)

// GetDiagnostics makes a request to the StorageOS API for a cluster diagnostic
// bundle to be generated and returned to the client.
//
// Because the OpenAPI code generator produces broken code for this method,
// we source the target path, authorization token and http client from it but
// handle the response ourselves.
func (o *OpenAPI) GetDiagnostics(ctx context.Context) (*diagnostics.BundleReadCloser, error) {
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
	req.Header["Accept"] = []string{"application/octet-stream", "application/gzip", "application/json"}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bundleCloser, err := o.getFileFromResp(resp)
	if err != nil {
		return nil, err
	}

	return bundleCloser, nil
}

// GetSingleNodeDiagnostics makes a request to the StorageOS API for a single
// node cluster diagnostic bundle to be generated and returned to the client.
//
// Because the OpenAPI code generator produces broken code for this method, we
// source the target path, authorization token and http client from it but
// handle the response ourselves.
func (o *OpenAPI) GetSingleNodeDiagnostics(ctx context.Context, nodeID id.Node) (*diagnostics.BundleReadCloser, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	// Get the appropriate config settings from the openapi client
	token := o.client.GetConfig().DefaultHeader["Authorization"]
	targetEndpoint := o.client.GetConfig().Scheme + "://" + o.client.GetConfig().Host + "/" + o.client.GetConfig().BasePath + "/diagnostics" + "/" + nodeID.String()
	client := o.client.GetConfig().HTTPClient

	// Construct the request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header["Authorization"] = []string{token}
	req.Header["Accept"] = []string{"application/octet-stream", "application/gzip", "application/json"}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bundleCloser, err := o.getFileFromResp(resp)
	if err != nil {
		return nil, err
	}

	return bundleCloser, nil
}

func (o *OpenAPI) getFileFromResp(resp *http.Response) (*diagnostics.BundleReadCloser, error) {
	var name string
	if extracted, err := getFilenameFromHeader(resp.Header); err == nil {
		name = extracted
	}

	switch resp.StatusCode {
	case http.StatusOK:
		// Carry on.
	case http.StatusBadGateway:
		// Check if the response content-type indicates a partial bundle. That
		// is, it has a gzip or octet-stream content type.
		for _, value := range resp.Header["Content-Type"] {
			switch value {
			case "application/gzip", "application/octet-stream":
				return nil, apiclient.NewIncompleteDiagnosticsError(
					diagnostics.NewBundleReadCloser(resp.Body, name),
				)
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

	return diagnostics.NewBundleReadCloser(resp.Body, name), nil
}

// getFileNameFromHeader attempts to extract an attachment filename from the
// provided HTTP response header.
func getFilenameFromHeader(header http.Header) (string, error) {

	// Try grab a name from the content disposition header.
	//
	// Expected form if present is `attachment; filename="some-name.ext"`.
	for _, value := range header["Content-Disposition"] {
		parts := strings.Split(value, ";")
		// If it doesn't split at least in two, can't have filename key and be
		// correct
		if len(parts) != 2 {
			continue
		}

		if parts[0] != "attachment" {
			continue
		}

		parts = strings.Split(parts[1], "=")
		// If the second part doesn't split in two on an equals sign, can't be
		// correct
		if len(parts) != 2 {
			continue
		}

		if strings.Trim(parts[0], " ") != "filename" {
			continue
		}

		// Cut quotes and whitespace from head and tail, then break out
		name := strings.Trim(parts[1], " \"")
		return name, nil
	}

	return "", errExtractingFilename
}
