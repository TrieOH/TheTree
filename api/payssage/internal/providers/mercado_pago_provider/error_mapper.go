package mercado_pago_provider

import (
	"encoding/json"
	"log/slog"

	"payssage/errs"

	"lib/telemetry"
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
// production — check CauseCode on logged errs.ProviderError values.
var knownMPErrorCodes = map[int]string{
	13253: "seller is not eligible for Pix: no active Pix key registered on their MercadoPago account",
	// 2034: "invalid card token: token expired or already used",
	// 4037: "invalid payment_method_id for this collector's country/account",
}

// mapMPError converts a MercadoPago error response into an
// errs.ProviderError. statusCode is the HTTP status MP returned; resultErr
// is whatever was populated via resty's SetResultError (map[string]any or
// similar). Payssage's problem is picked from statusCode by the problem switch
// (internal/app). It logs MP's error code and message under "MP <op>
// failed" with attrs, never the response body: MP echoes the payer back.
func mapMPError(op string, statusCode int, resultErr any, attrs ...any) error {
	var body mpErrorBody
	raw, err := json.Marshal(resultErr)
	if err == nil {
		_ = json.Unmarshal(raw, &body)
	}

	msg := body.Message
	if msg == "" {
		msg = "mercadopago request failed"
	}

	perr := &errs.ProviderError{
		Provider: "mercadopago",
		Status:   statusCode,
		Message:  msg,
		Code:     body.Error,
	}
	if len(body.Cause) > 0 {
		cause := body.Cause[0]
		perr.CauseCode = cause.Code
		perr.CauseDescription = cause.Description
		perr.Reason = knownMPErrorCodes[cause.Code]
	}
	telemetry.Log().Warn("MP "+op+" failed", append(attrs,
		slog.Int("status", statusCode),
		slog.String("mp_error", perr.Code),
		slog.String("mp_message", perr.Message),
		slog.Int("mp_cause_code", perr.CauseCode),
	)...)
	return perr
}
