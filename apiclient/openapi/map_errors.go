package openapi

import (
	"net/http"

	"code.storageos.net/storageos/c2-cli/apiclient"
)

func mapResponseToError(resp *http.Response) error {
	switch resp.StatusCode {
	// 2XX
	case http.StatusOK, http.StatusAccepted:
		return nil
	// 4XX
	case http.StatusBadRequest:
		return apiclient.ErrBadRequest
	case http.StatusUnauthorized:
		return apiclient.ErrAuthenticationRequired
	case http.StatusForbidden:
		return apiclient.ErrUnauthorised
	case http.StatusNotFound:
		return apiclient.ErrNotFound
	// TODO: StatusConflict maps to ErrAlreadyExists and ErrInUse
	case http.StatusConflict:
		return apiclient.ErrInUse
	case http.StatusPreconditionFailed:
		return apiclient.ErrStaleWrite
	case http.StatusUnprocessableEntity:
		return apiclient.ErrInvalidStateTransition
	case http.StatusUnavailableForLegalReasons:
		return apiclient.ErrLicenceCapacityExceeded
	// 5XX
	case http.StatusInternalServerError: // 500
		return apiclient.ErrServerError
	case http.StatusServiceUnavailable: // 503
		return apiclient.ErrStoreError
	default:
		return apiclient.ErrUnknown
	}
}
