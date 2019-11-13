package openapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/openapi"
)

// TODO: Surely there's a better way?
// apiErr defines a JSON encodable struct with an error field.
type apiErr struct {
	Error string `json:"error"`
}

// mapOpenAPIError will given err and the corresponding resp attempt to map the
// error value to an apiclient error type.
//
// err is returned as is when any of the following are true:
// 	→ resp is nil
// 	→ err is not a GenericOpenAPIError
func mapOpenAPIError(err error, resp *http.Response) error {
	if resp == nil {
		return err
	}

	var oerr openapi.GenericOpenAPIError
	if ok := errors.As(err, &oerr); !ok {
		return err
	}

	var details string
	switch resp.Header.Get("Content-Type") {
	case "application/json":
		// TODO: the error doesn't do anything useful here, as we default
		// to an empty string.
		details, _ = extractErrorStringJSON(oerr.Body())
	default:
	}

	switch resp.StatusCode {

	// 4XX
	case http.StatusBadRequest:
		return apiclient.NewBadRequestError(details)

	case http.StatusUnauthorized:
		return apiclient.NewAuthenticationError(details)

	case http.StatusForbidden:
		return apiclient.NewUnauthorisedError(details)

	case http.StatusNotFound:
		return apiclient.NewNotFoundError(details)

	case http.StatusConflict:
		return apiclient.NewConflictError(details)

	case http.StatusPreconditionFailed:
		return apiclient.NewStaleWriteError(details)

	case http.StatusUnprocessableEntity:
		return apiclient.NewInvalidStateTransitionError(details)

	// TODO(CP-3925): This may need changing to present a friendly error, or
	// it may be done up the call stack.
	case http.StatusUnavailableForLegalReasons:
		return apiclient.NewLicenceCapabilityError(details)

	// 5XX
	case http.StatusInternalServerError:
		return apiclient.NewServerError(details)

	case http.StatusServiceUnavailable:

		return apiclient.NewStoreError(details)
	default:
		return err
	}
}

func extractErrorStringJSON(body []byte) (string, error) {
	var e apiErr

	if err := json.Unmarshal(body, &e); err != nil {
		return "", err
	}

	return e.Error, nil
}
