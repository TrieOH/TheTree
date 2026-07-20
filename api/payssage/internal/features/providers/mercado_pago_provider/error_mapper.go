package mercado_pago_provider

import (
	"encoding/json"

	"github.com/MintzyG/fun"
)

// mpErrorCause mirrors a single entry in MercadoPago's error response "cause" array.
type mpErrorCause struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Data        string `json:"data,omitempty"`
}

// mpErrorBody mirrors MercadoPago's standard error response shape,
// e.g. {"message":"...", "error":"bad_request", "status":400, "cause":[...]}
type mpErrorBody struct {
	Message string         `json:"message"`
	Error   string         `json:"error"`
	Status  int            `json:"status"`
	Cause   []mpErrorCause `json:"cause"`
}

// knownMPErrorCodes maps MercadoPago's internal cause codes to a stable,
// actionable message for callers/support. Extend as new codes are hit in
// production — check the "mp_cause_code" meta field in logged AppErrors.
var knownMPErrorCodes = map[int]string{
	13253: "seller is not eligible for Pix: no active Pix key registered on their MercadoPago account",
	// 2034: "invalid card token: token expired or already used",
	// 4037: "invalid payment_method_id for this collector's country/account",
}

// mapMPError converts a MercadoPago error response into a *fun.AppError.
// statusCode is the HTTP status MP returned; resultErr is whatever was
// populated via resty's SetResultError (map[string]any or similar).
func mapMPError(statusCode int, resultErr any) *fun.AppError {
	var body mpErrorBody
	if raw, err := json.Marshal(resultErr); err == nil {
		_ = json.Unmarshal(raw, &body)
	}

	msg := body.Message
	if msg == "" {
		msg = "mercadopago request failed"
	}

	builder := fun.Errf("mercadopago: %s", msg).
		WithMeta("mp_status", body.Status).
		WithMeta("mp_error", body.Error)

	if len(body.Cause) > 0 {
		cause := body.Cause[0]
		builder = builder.WithMeta("mp_cause_code", cause.Code).
			WithMeta("mp_cause_description", cause.Description)

		if reason, ok := knownMPErrorCodes[cause.Code]; ok {
			builder = builder.WithMeta("mp_reason", reason)
		}
	}

	switch {
	case statusCode == 401 || statusCode == 403:
		return builder.Unauthorized()
	case statusCode == 404:
		return builder.NotFound()
	case statusCode == 409:
		return builder.Conflict()
	case statusCode == 429:
		return builder.TooManyRequests()
	case statusCode >= 500:
		return builder.BadGateway() // MP itself failed — treat as upstream/gateway issue
	default:
		return builder.BadRequest() // covers MP's 400s, including code 13253
	}
}
