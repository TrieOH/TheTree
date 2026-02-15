package errx

import (
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
)

type HTTPResponse struct {
	resp.Response
}

type HTTPTranslator struct{}

var (
	CannotTranslateUnsupportedError = fail.ID(5, "TR", 0, false, "TRCannotTranslateUnsupportedError")
	ErrCannotTranslate              = fail.Form(CannotTranslateUnsupportedError, "cannot translate unsupported error ID(%s) to http", true, nil, "MISSING ID")
	CannotTranslateNilError         = fail.ID(5, "TR", 1, false, "TRCannotTranslateNilError")
	ErrCannotTranslateNilError      = fail.Form(CannotTranslateNilError, "cannot translate nil error to %s", true, nil, "MISSING DOMAIN")
)

func HTTPResponseTranslator() *HTTPTranslator {
	return &HTTPTranslator{}
}

func (h *HTTPTranslator) Name() string { return "http" }

func (h *HTTPTranslator) Supports(err *fail.Error) error {
	if err == nil {
		return fail.New(CannotTranslateNilError).WithArgs("http").Render()
	}

	code, ok := err.Meta["code"].(int)
	if !ok {
		code = 500
		return fail.New(CannotTranslateUnsupportedError).WithArgs(err.ID).Trace("Couldn't cast code")
	} else if code < 100 || code > 599 {
		code = 500
		return fail.New(CannotTranslateUnsupportedError).WithArgs(err.ID).Trace("Outside valid code range")
	}

	return nil
}

func (h *HTTPTranslator) Translate(err *fail.Error) (any, error) {
	traces := toStringSlice(err.Meta["traces"])
	delete(err.Meta, "traces")

	if err.Cause != nil {
		traces = append(traces, err.Cause.Error())
	}

	var module string
	var code int
	var ok bool
	if err.Meta != nil {
		code, ok = err.Meta["code"].(int)
		delete(err.Meta, "code")
		if !ok {
			code = 500
		} else if code < 100 || code > 599 {
			code = 500
		}

		module, ok = err.Meta["module"].(string)
		delete(err.Meta, "module")
		if !ok {
			module = "univents"
		}
	} else {
		err.Meta = map[string]any{}
		code = 500
	}

	_ = err.Render()

	r := resp.Response{
		Module:         module,
		Message:        err.Message,
		Data:           err.Meta,
		Trace:          traces,
		Timestamp:      time.Now(),
		PaginationData: nil,
		Code:           code,
		ErrorID:        err.ID.String(),
		ContentType:    "application/json",
		TracePrefix:    "",
	}

	return &r, nil
}

func toStringSlice(v any) []string {
	if v == nil {
		return []string{}
	}

	if s, ok := v.([]string); ok {
		return s
	}

	if arr, ok := v.([]any); ok {
		result := make([]string, 0, len(arr))
		for _, item := range arr {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}

	return []string{}
}
