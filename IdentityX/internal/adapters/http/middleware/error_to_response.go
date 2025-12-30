package middleware

import (
	"GoAuth/internal/apierr"
	"errors"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

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
