package apierr

import (
	"errors"
	"univents/internal/adapters/observability/logs"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
)

func ErrToResp(err error) *resp.Response {
	if err == nil {
		return nil
	}

	var fe *fail.Error
	if errors.As(err, &fe) {
		var rs *resp.Response
		rs, err = fail.ToAs[*resp.Response](fe, "http")
		if err != nil {
			return resp.InternalServerError().WithData(err)
		}
		return rs
	}

	logs.L().Error("FAILED ErrToResp")
	// unknown error = 500
	return resp.InternalServerError().
		WithTracePrefix("unhandled-error").
		AddTrace(err.Error())
}
