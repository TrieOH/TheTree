package mercado_pago_provider

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"payssage/internal/providers/mercado_pago"
	"payssage/models"
	"strconv"

	"lib/telemetry"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"resty.dev/v3"
)

func (p *Provider) Checkout(ctx context.Context, intent *models.Intent, providerCheckoutData json.RawMessage) (json.RawMessage, error) {
	checkoutData, err := mercado_pago.ParseCheckoutData(providerCheckoutData)
	if err != nil {
		return nil, err
	}

	seller, err := p.sellers.GetByID(ctx, intent.SellerID)
	if err != nil {
		return nil, fmt.Errorf("resolve seller: %w", err)
	}
	if seller.RevokedAt != nil {
		return nil, fmt.Errorf("seller %s has revoked provider access", seller.ID)
	}
	if seller.Provider != "mercadopago" {
		return nil, fmt.Errorf("seller %s is not a mercadopago seller (got %q)", seller.ID, seller.Provider)
	}

	var creds models.MercadoPagoCredentials
	err = json.Unmarshal(seller.Credentials, &creds)
	if err != nil {
		return nil, fmt.Errorf("unmarshal seller credentials: %w", err)
	}
	if creds.AccessToken == "" {
		return nil, fmt.Errorf("seller %s missing mercadopago access token", seller.ID)
	}

	// TODO: check token expiry / refresh via creds.RefreshToken before use
	// once MP token TTL tracking is added to Credentials.

	body := map[string]any{
		"transaction_amount":   json.Number(formatAmount(intent.AmountCents)),
		"application_fee":      json.Number(formatAmount(calcApplicationFee(intent.AmountCents, checkoutData.MarketplaceFeeBPS))),
		"installments":         checkoutData.Installments,
		"token":                checkoutData.Token,
		"payment_method_id":    checkoutData.PaymentMethodID,
		"external_reference":   intent.ID.String(),
		"statement_descriptor": cmp.Or(checkoutData.StatementDescriptor, "payssage"),
		"payer": map[string]any{
			"email": checkoutData.Payer.Email,
			"identification": map[string]any{
				"type":   checkoutData.Payer.IdentificationType,
				"number": checkoutData.Payer.IdentificationNumber,
			},
		},
		"additional_info": buildAdditionalInfo(checkoutData.AdditionalInfo, intent.AmountCents),
	}

	telemetry.Log().Info("MP Create Payment Request", zap.Any("body", body))

	// TODO remember to use external group feature we made

	// TODO: fix the idemp key not being actual idemp
	idempotencyKey, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	// TODO: hoist this into a shared *resty.Client on Provider (constructed once,
	// with base URL / timeout / retry policy) instead of creating one per call.
	client := resty.New()
	defer client.Close()

	var mpResp struct {
		ID           int    `json:"id"`
		Status       string `json:"status"`
		StatusDetail string `json:"status_detail"`
	}
	var mpErr map[string]any

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+creds.AccessToken).
		SetHeader("X-Idempotency-Key", idempotencyKey.String()).
		SetBody(body).
		SetResult(&mpResp).
		SetResultError(&mpErr).
		Post("https://api.mercadopago.com/v1/payments")
	if err != nil {
		return nil, fmt.Errorf("mercadopago create payment request: %w", err)
	}

	telemetry.Log().Info("MP Create Payment Request Raw Body", zap.String("body", resp.String()))

	if resp.IsStatusFailure() {
		return nil, fmt.Errorf("mercadopago create payment error %d: %v", resp.StatusCode(), resp.ResultError())
	}

	providerData := &models.MercadoPagoIntentData{
		OrderID:           strconv.Itoa(mpResp.ID),
		TransactionID:     strconv.Itoa(mpResp.ID),
		OrderStatus:       mpResp.Status,
		OrderStatusDetail: mpResp.StatusDetail,
		PaymentMethodID:   checkoutData.PaymentMethodID,
	}
	providerDataBytes, err := json.Marshal(providerData)
	if err != nil {
		return nil, err
	}
	return providerDataBytes, nil
}

// Charge performs a direct server-side charge.
// Use for server-to-server flows where you already have a payment method token.
func (p *Provider) Charge(ctx context.Context, intent *models.Intent) (*models.Intent, error) {
	return nil, errors.New("not implemented")
}

// Refund issues a full or partial refund against a prior transaction.
func (p *Provider) Refund(ctx context.Context, intent *models.Intent) (*models.Intent, error) {
	return nil, errors.New("not implemented")
}
