package apierr

import (
	"log"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail"
)

type HTTPResponse struct {
	resp.Response
}

type HTTPTranslator struct{}

var (
	CannotTranslateToHTTP = fail.ID(5, "TR", 0, false, "TRCannotTranslateToHTTP")
	ErrCannotTranslate    = fail.Form(CannotTranslateToHTTP, "cannot translate error to http", true, nil)
)

func HTTPResponseTranslator() *HTTPTranslator {
	return &HTTPTranslator{}
}

func (h *HTTPTranslator) Name() string { return "http" }
func (h *HTTPTranslator) Supports(err *fail.Error) bool {
	if !err.IsTrusted() || err == nil {
		return false
	}
	switch err.ID {
	case RequestMissingQueryParamValue, RequestMissingQueryParam, RequestMissingSchemaCustomFields,
		RequestInvalidJSONFormat, RequestValidationError, RequestNotApplicationJSON, RequestEmptyCookie,
		RequestUnknownQueryParam:
		return true
	default:
		return false
	}
}

func (h *HTTPTranslator) Translate(err *fail.Error) (any, error) {
	if !h.Supports(err) {
		return nil, fail.New(CannotTranslateToHTTP).With(err)
	}

	traces, ok := err.Meta["traces"].([]string)
	if !ok {
		traces = []string{}
	}

	if err.Cause != nil {
		traces = append(traces, err.Cause.Error())
	}

	var code int
	if err.Meta != nil {
		err.Meta["traces"] = ""
		code, ok = err.Meta["code"].(int)
		if !ok {
			code = 500
		} else if code < 100 || code > 599 {
			code = 500
		}
	} else {
		err.Meta = map[string]any{}
		code = 500
	}

	err.Render()

	r := resp.Response{
		Module:         "translator",
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
	log.Println(r)

	return &r, nil
}
