package apierr

import (
	"errors"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail"
)

// ErrToResp converts an error to a response.
// It handles API errors and returns a formatted response.
// For unhandled errors, it returns an internal server error response.
// Debug causes are included based on the global IncludeDebugCauses flag.
func ErrToResp(err error) *resp.Response {
	if err == nil {
		return nil
	}

	var ae *Error
	if errors.As(err, &ae) {
		return MapAPIErrorWithTrace(ae)
	}

	var fe *fail.Error
	if errors.As(err, &fe) {
		return Sender2(fe)
	}

	// unknown error = 500
	return resp.InternalServerError().
		WithTracePrefix("unhandled-error").
		AddTrace(err.Error())
}

func Sender2(e *fail.Error) *resp.Response {
	trrs, err := HTTPResponseTranslator().Translate(e)
	if err != nil {
		return resp.InternalServerError().WithData(err)
	}
	if rs, ok := trrs.(*resp.Response); ok {
		if rs != nil {
			return rs
		}
		return resp.InternalServerError("response was nil").WithData(err)
	} else {
		return resp.InternalServerError("couldn't cast to response")
	}
}
