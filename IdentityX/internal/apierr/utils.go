package apierr

import (
	"GoAuth/internal/adapters/observability/logs"
	"errors"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail"
)

func ErrToResp(err error) *resp.Response {
	if err == nil {
		return nil
	}

	var ae *Error
	if errors.As(err, &ae) {
		logs.L().Error("CALLED MapAPIErrorWithTrace")
		return MapAPIErrorWithTrace(ae)
	}

	var fe *fail.Error
	if errors.As(err, &fe) {
		return Sender2(fe)
	}

	logs.L().Error("FAILED ErrToResp")
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
