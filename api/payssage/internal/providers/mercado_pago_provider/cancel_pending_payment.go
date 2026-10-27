package mercado_pago_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"lib/telemetry"
	"payssage/errs"
	"payssage/models"

	"go.uber.org/zap"
)

func (p *Provider) CancelPendingPayment(ctx context.Context, intent *models.Intent) error {
	var providerData models.MercadoPagoIntentData
	err := json.Unmarshal(intent.ProviderData, &providerData)
	if err != nil {
		return fmt.Errorf("error mapping mercadopago provider data: %w", err)
	}
	if providerData.TransactionID == "" {
		return errs.ErrNoTransactionToCancel
	}

	seller, err := p.sellers.GetByID(ctx, intent.SellerID)
	if err != nil {
		return fmt.Errorf("resolve seller: %w", err)
	}
	if seller.RevokedAt != nil {
		return &errs.SellerRevokedError{SellerID: seller.ID, Op: "cancel"}
	}

	var creds models.MercadoPagoCredentials
	err = json.Unmarshal(seller.Credentials, &creds)
	if err != nil {
		return fmt.Errorf("error mapping seller credentials: %w", err)
	}
	if creds.AccessToken == "" {
		return fmt.Errorf("seller %s missing mercadopago access token", seller.ID)
	}

	body := map[string]any{"status": "cancelled"}

	var mpResp struct {
		Status       string `json:"status"`
		StatusDetail string `json:"status_detail"`
	}
	var mpErr map[string]any

	resp, err := p.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+creds.AccessToken).
		SetBody(body).
		SetResult(&mpResp).
		SetResultError(&mpErr).
		Put("https://api.mercadopago.com/v1/payments/" + providerData.TransactionID)
	if err != nil {
		return fmt.Errorf("%w: mercadopago cancel payment request: %w", errs.ErrProviderRequestFailed, err)
	}

	telemetry.Log().Info("MP Cancel Payment Response",
		zap.String("payment_id", providerData.TransactionID),
		zap.String("body", resp.String()),
	)

	if resp.IsStatusFailure() {
		return mapMPError(resp.StatusCode(), resp.ResultError())
	}

	providerData.OrderStatus = mpResp.Status
	providerData.OrderStatusDetail = mpResp.StatusDetail

	providerDataBytes, err := json.Marshal(&providerData)
	if err != nil {
		return fmt.Errorf("marshal provider data: %w", err)
	}
	intent.ProviderData = providerDataBytes
	intent.Status = p.NormalizeStatus(mpResp.Status)
	intent.StatusDetail = p.MapStatusDetail(mpResp.StatusDetail)

	return nil
}
