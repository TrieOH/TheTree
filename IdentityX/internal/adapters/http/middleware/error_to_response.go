package middleware

import (
	"GoAuth/internal/apierr"
	"errors"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// ErrToResp converts an error to a response.
// It handles API errors and returns a formatted response.
// For unhandled errors, it returns an internal server error response.
func ErrToResp(err error) *resp.Response {
	if err == nil {
		return nil
	}

	var ae *apierr.Error
	if errors.As(err, &ae) {
		return apierr.MapAPIError(ae)
	}

	// unknown error = 500
	return resp.InternalServerError().
		WithTracePrefix("unhandled-error").
		AddTrace(err.Error())
}
