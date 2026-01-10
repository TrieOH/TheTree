package apierr

import (
	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

func MapAPIError(e *Error) *resp.Response {
	r := respFromCode(e.Code)

	if e.Message != "" {
		r = r.WithMsg(e.Message)
	}

	if e.ID != "" {
		r = r.WithErrID(string(e.ID))
	}

	if e.Cause != nil {
		r = r.AddTrace(e.Cause.Error())
	}

	return r
}

func respFromCode(code Code) *resp.Response {
	switch code {
	case InvalidInput, BadRequest:
		return resp.BadRequest()

	case Unauthorized:
		return resp.Unauthorized()

	case Forbidden:
		return resp.Forbidden()

	case NotFound:
		return resp.NotFound()

	case Conflict:
		return resp.Conflict()

	default:
		return resp.InternalServerError()
	}
}
