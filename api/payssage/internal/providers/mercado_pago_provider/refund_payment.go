package mercado_pago_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"lib/telemetry"
	"math"
	"payssage/errs"
	"payssage/models"
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Refund fully reverses an approved MercadoPago payment (full refund only —
// no amount, v1). The intent's Status stays `succeeded` on success: the
// payment.refunded webhook is what flips it (webhook-only confirmation,
// mirroring approval, D3). The refund id/status lands in provider_data for
// observability (and for the fee-on-refund verification, R1).
func (p *Provider) Refund(ctx context.Context, intent *models.Intent) error {
	var providerData models.MercadoPagoIntentData
	err := json.Unmarshal(intent.ProviderData, &providerData)
	if err != nil {
		return fmt.Errorf("error mapping mercadopago provider data: %w", err)
	}
	if providerData.TransactionID == "" {
		return errs.ErrNoTransactionToRefund
	}

	seller, err := p.sellers.GetByID(ctx, intent.SellerID)
	if err != nil {
		return fmt.Errorf("resolve seller: %w", err)
	}
	if seller.RevokedAt != nil {
		return &errs.SellerRevokedError{SellerID: seller.ID, Op: "refund"}
	}

	var creds models.MercadoPagoCredentials
	err = json.Unmarshal(seller.Credentials, &creds)
	if err != nil {
		return fmt.Errorf("error mapping seller credentials: %w", err)
	}
	if creds.AccessToken == "" {
		return fmt.Errorf("seller %s missing mercadopago access token", seller.ID)
	}

	// MP requires X-Idempotency-Key on refunds (4292 otherwise); a fresh key
	// per attempt is enough for v1 — presence is what MP validates, and the
	// service-layer succeeded-only guard plus webhook confirmation make a
	// duplicate refund a no-op rather than a double charge.
	idempotencyKey, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generate idempotency key: %w", err)
	}

	var mpResp struct {
		ID           int64   `json:"id"`
		Status       string  `json:"status"`
		StatusDetail string  `json:"status_detail"`
		Amount       float64 `json:"amount"`
	}
	var mpErr map[string]any

	resp, err := p.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+creds.AccessToken).
		SetHeader("X-Idempotency-Key", idempotencyKey.String()).
		SetBody(map[string]any{}).
		SetResult(&mpResp).
		SetResultError(&mpErr).
		Post("https://api.mercadopago.com/v1/payments/" + providerData.TransactionID + "/refunds")
	if err != nil {
		return fmt.Errorf("%w: mercadopago refund payment request: %w", errs.ErrProviderRequestFailed, err)
	}

	telemetry.Log().Info("MP Refund Payment Response",
		zap.String("payment_id", providerData.TransactionID),
		zap.String("body", resp.String()),
	)

	if resp.IsStatusFailure() {
		return mapMPError(resp.StatusCode(), resp.ResultError())
	}

	// The refund POST returns the refund object; the payment itself is still
	// `approved` until MP processes it (the payment.refunded webhook follows).
	refundID := strconv.FormatInt(mpResp.ID, 10)
	providerData.RefundID = &refundID
	if mpResp.Status != "" {
		status := mpResp.Status
		providerData.RefundStatus = &status
	}
	if mpResp.StatusDetail != "" {
		detail := mpResp.StatusDetail
		providerData.RefundStatusDetail = &detail
	}
	if mpResp.Amount != 0 {
		cents := int64(math.Round(mpResp.Amount * 100))
		providerData.RefundAmountCents = &cents
	}

	providerDataBytes, err := json.Marshal(&providerData)
	if err != nil {
		return fmt.Errorf("marshal provider data: %w", err)
	}
	intent.ProviderData = providerDataBytes
	// Intent status intentionally untouched — the webhook confirms (D3/D-2).

	return nil
}
