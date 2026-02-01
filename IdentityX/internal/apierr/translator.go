package apierr

import (
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
	case RequestMissingQueryParamValue,
		RequestMissingQueryParam,
		RequestMissingSchemaCustomFields,
		RequestInvalidJSONFormat,
		RequestValidationError,
		RequestNotApplicationJSON,
		RequestEmptyCookie,
		RequestUnknownQueryParam:
		return true
	case SQLNotFound:
		return true
	case AuthEmailAlreadyUsed,
		AuthInvalidCredentials,
		AuthInvalidRefreshCookie,
		AuthInvalidAccessCookie,
		AuthMissingRefreshCookie,
		AuthMissingAccessCookie,
		AuthInvalidPrincipal,
		AuthInvalidPassword,
		AuthNotClient,
		AuthNotProjectUser,
		AuthAlreadyVerified,
		AuthPrincipalNotInContext:
		return true
	case SessionRevoked,
		SessionNotFound,
		SessionSelfRevokeForbidden,
		SessionUnauthorized:
		return true
	case TokenInvalid,
		TokenExpired,
		TokenMalformed,
		TokenSignatureInvalid,
		TokenInvalidAlg,
		TokenCouldNotSign,
		TokenInvalidAccessClaims,
		TokenNotYetValid,
		TokenUsedBeforeIssued,
		TokenInvalidIssuer,
		TokenInvalidSubject,
		TokenInvalidAudience,
		TokenRefreshInvalidID,
		TokenAccessInvalidID,
		TokenUntrusted,
		TokenSessionMismatch,
		TokenMismatchDuringAuth:
		return true
	case SCHEMANoPublishedVersion:
		return true
	default:
		return false
	}
}

func (h *HTTPTranslator) Translate(err *fail.Error) (any, error) {
	if !h.Supports(err) {
		return nil, fail.New(CannotTranslateToHTTP).With(err)
	}

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
			module = "GoAuth"
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
