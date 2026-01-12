package apierr

import (
	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

var IncludeDebugCauses = false

// MapAPIErrorWithTrace maps an API error to a response with full trace information.
// Debug causes are included based on the global IncludeDebugCauses flag.
func MapAPIErrorWithTrace(e *Error) *resp.Response {
	r := respFromCode(e.Code)

	if e.Message != "" {
		r = r.WithMsg(e.Message)
	}

	if e.ID != "" {
		if e.ID == RequestValidationError {
			r.WithModule("validation")
		}
		r = r.WithErrID(string(e.ID))
	}

	// Add all causes to trace
	allCauses := e.GetAllCauses()
	for _, cause := range allCauses {
		r = r.AddTrace(cause.Error())
	}

	// Optionally add debug causes based on global flag
	if IncludeDebugCauses {
		for _, debugCause := range e.DebugCauses {
			r = r.AddTrace("[DEBUG] " + debugCause.Error())
		}
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
