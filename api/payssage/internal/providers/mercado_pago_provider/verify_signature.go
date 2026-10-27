package mercado_pago_provider

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"

	"payssage/errs"
)

func (p *Provider) VerifySignature(_ context.Context, r *http.Request, _ []byte) error {
	signatureHeader := r.Header.Get("x-signature")
	if signatureHeader == "" {
		return &errs.MissingParamError{In: "header", Name: "x-signature"}
	}
	requestID := r.Header.Get("x-request-id")
	if requestID == "" {
		return &errs.MissingParamError{In: "header", Name: "x-request-id"}
	}
	dataID := r.URL.Query().Get("data.id")
	if dataID == "" {
		return &errs.MissingParamError{In: "query", Name: "data.id"}
	}

	// MercadoPago's x-signature header is a comma-separated list of
	// key=value pairs. As of their current webhook spec, it always
	// contains exactly two: "ts" (the Unix timestamp the signature was
	// generated at) and "v1" (the HMAC-SHA256 signature itself, hex
	// encoded). Both key names ("ts", "v1") are MercadoPago's own wire
	// format and must match exactly — they are not ours to rename.
	// Example header value: "ts=1704067200,v1=7f3a9c1e2b8d4f6a..."
	var signatureTimestamp, providedSignature string
	for part := range strings.SplitSeq(signatureHeader, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		fieldName := strings.TrimSpace(kv[0])
		fieldValue := strings.TrimSpace(kv[1])
		switch fieldName {
		case "ts":
			signatureTimestamp = fieldValue
		case "v1":
			providedSignature = fieldValue
		}
	}
	if signatureTimestamp == "" || providedSignature == "" {
		return errs.ErrWebhookSignatureMissing
	}

	// The manifest is the exact string MercadoPago signs on their end —
	// this format (field:value; pairs, this specific field order) is
	// dictated by MP and must be reproduced exactly or the HMAC won't match.
	manifest := fmt.Sprintf("id:%s;request-id:%s;ts:%s;", dataID, requestID, signatureTimestamp)

	mac := hmac.New(sha256.New, []byte(os.Getenv("MP_WEBHOOK_SECRET")))
	mac.Write([]byte(manifest))
	computedSignature := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(computedSignature), []byte(providedSignature)) {
		return errs.ErrWebhookSignatureMismatch
	}
	return nil
}
