package goauth

import (
	"fmt"
	"net/http"

	"github.com/MintzyG/fail/v3"
)

var (
	SDKUnknownErrorID              = fail.ID(0, "SDK", 0, false, "SDKUnknownError")
	SDKNetworkErrorID              = fail.ID(0, "SDK", 1, false, "SDKNetworkError")
	SDKRequestMarshalingErrorID    = fail.ID(0, "SDK", 2, false, "SDKRequestMarshalingError")
	SDKResponseUnmarshalingErrorID = fail.ID(0, "SDK", 3, false, "SDKResponseUnmarshalingError")
	SDKInvalidAPIKeyID             = fail.ID(0, "SDK", 4, false, "SDKInvalidAPIKey")
	SDKUnauthorizedID              = fail.ID(0, "SDK", 5, false, "SDKUnauthorized")
	SDKForbiddenID                 = fail.ID(0, "SDK", 6, false, "SDKForbidden")
	SDKNotFoundID                  = fail.ID(0, "SDK", 7, false, "SDKNotFound")
	SDKRateLimitedID               = fail.ID(0, "SDK", 8, false, "SDKRateLimited")
	SDKInternalServerErrorID       = fail.ID(0, "SDK", 9, false, "SDKInternalServerError")
	SDKBadRequestID                = fail.ID(0, "SDK", 10, false, "SDKBadRequest")
	SDKInvalidObjectFormatID       = fail.ID(0, "SDK", 11, false, "SDKInvalidObjectFormat")
	SDKInvalidActionFormatID       = fail.ID(0, "SDK", 12, false, "SDKInvalidActionFormat")
	SDKConflictID                  = fail.ID(0, "SDK", 13, false, "SDKConflict")
	SDKMissingUserID               = fail.ID(0, "SDK", 14, false, "SDKMissingUserID")
	SDKUnexpectedSigningMethod     = fail.ID(0, "SDK", 15, false, "SDKUnexpectedSigningMethod")
	SDKMissingTokenKID             = fail.ID(0, "SDK", 16, false, "SDKMissingTokenKID")
	SDKKeyNotInJWKS                = fail.ID(0, "SDK", 17, false, "SDKKeyNotInJWKS")
	SDKUnsupportedCurve            = fail.ID(0, "SDK", 18, false, "SDKUnsupportedCurve")
	SDKKeyDecodeFailed             = fail.ID(0, "SDK", 19, false, "SDKKeyDecodeFailed")
	SDKInvalidKeySize              = fail.ID(0, "SDK", 20, false, "SDKInvalidKeySize")

	ErrSDKUnknownError              = fail.Form(SDKUnknownErrorID, "unknown sdk error: %s", false, map[string]any{"code": 500}, "UNSET")
	ErrSDKNetworkError              = fail.Form(SDKNetworkErrorID, "network error", false, map[string]any{"code": 500})
	ErrSDKRequestMarshalingError    = fail.Form(SDKRequestMarshalingErrorID, "request marshaling error: %s", false, map[string]any{"code": 500}, "UNSET")
	ErrSDKResponseUnmarshalingError = fail.Form(SDKResponseUnmarshalingErrorID, "response unmarshalling error: %s", false, map[string]any{"code": 500}, "UNSET")
	ErrSDKInvalidAPIKey             = fail.Form(SDKInvalidAPIKeyID, "invalid api key", false, map[string]any{"code": 401})
	ErrSDKUnauthorized              = fail.Form(SDKUnauthorizedID, "unauthorized", false, map[string]any{"code": 401})
	ErrSDKForbidden                 = fail.Form(SDKForbiddenID, "forbidden", false, map[string]any{"code": 403})
	ErrSDKNotFound                  = fail.Form(SDKNotFoundID, "not found", false, map[string]any{"code": 404})
	ErrSDKRateLimited               = fail.Form(SDKRateLimitedID, "rate limited", false, map[string]any{"code": 429})
	ErrSDKInternalServerError       = fail.Form(SDKInternalServerErrorID, "internal server error", false, map[string]any{"code": 500})
	ErrSDKBadRequest                = fail.Form(SDKBadRequestID, "bad request", false, map[string]any{"code": 400})
	ErrSDKInvalidObjectFormat       = fail.Form(SDKInvalidObjectFormatID, "invalid object format: %s", false, map[string]any{"code": 400})
	ErrSDKInvalidActionFormat       = fail.Form(SDKInvalidActionFormatID, "invalid action format: %s", false, map[string]any{"code": 400})
	ErrSDKConflict                  = fail.Form(SDKConflictID, "conflict", false, map[string]any{"code": 409})
	ErrSDKMissingUserID             = fail.Form(SDKMissingUserID, "missing user id", false, map[string]any{"code": 400})
	ErrUnexpectedSigningMethod      = fail.Form(SDKUnexpectedSigningMethod, "unexpected signing method: %s", false, map[string]any{"code": 400}, "UNSET")
	ErrMissingTokenKID              = fail.Form(SDKMissingTokenKID, "missing kid in token header", false, map[string]any{"code": 400})
	ErrKeyNotInJWKS                 = fail.Form(SDKKeyNotInJWKS, "key not found in JWKS: %s", false, map[string]any{"code": 404}, "UNSET")
	ErrUnsupportedCurve             = fail.Form(SDKUnsupportedCurve, "unsupported key type or curve: %s/%s", false, map[string]any{"code": 400}, "UNSET", "UNSET")
	ErrKeyDecodeFailed              = fail.Form(SDKKeyDecodeFailed, "failed to decode public key: %s", false, map[string]any{"code": 400}, "UNSET")
	ErrInvalidKeySize               = fail.Form(SDKInvalidKeySize, "invalid public key size", false, map[string]any{"code": 400})
)

type httpStatusError struct {
	status  int
	apiCode string
	apiID   string
	msg     string
	traces  []string
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("http status error: %d, api_id: %s, msg: %s", e.status, e.apiID, e.msg)
}

type HTTPMapper struct{}

func (m *HTTPMapper) Name() string  { return "http-status" }
func (m *HTTPMapper) Priority() int { return 10 }

func (m *HTTPMapper) Map(err error) (*fail.Error, bool) {
	if hErr, ok := err.(*httpStatusError); ok {
		fe := MapHTTPStatusToErr(hErr.status)
		if fe != nil {
			_ = fe.AddMeta("api_status", hErr.status).
				AddMeta("api_code", hErr.apiCode).
				AddMeta("api_id", hErr.apiID).
				AddMeta("api_traces", hErr.traces).
				Msg(hErr.msg)
			return fe, true
		}
	}

	return nil, false
}

func MapHTTPStatusToErr(status int) *fail.Error {
	switch status {
	case http.StatusBadRequest:
		return fail.New(SDKBadRequestID)
	case http.StatusUnauthorized:
		return fail.New(SDKUnauthorizedID)
	case http.StatusConflict:
		return fail.New(SDKConflictID)
	case http.StatusForbidden:
		return fail.New(SDKForbiddenID)
	case http.StatusNotFound:
		return fail.New(SDKNotFoundID)
	case http.StatusTooManyRequests:
		return fail.New(SDKRateLimitedID)
	case http.StatusInternalServerError:
		return fail.New(SDKInternalServerErrorID)
	default:
		return nil
	}
}

func init() {
	fail.RegisterMapper(&HTTPMapper{})
}
