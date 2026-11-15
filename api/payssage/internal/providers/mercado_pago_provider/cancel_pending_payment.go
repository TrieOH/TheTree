package mercado_pago_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"lib/telemetry"
	"log/slog"
	"payssage/errs"
	"payssage/models"
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

	if resp.IsStatusFailure() {
		return mapMPError("cancel payment", resp.StatusCode(), resp.ResultError(),
			slog.String("intent_id", intent.ID.String()),
			slog.String("payment_id", providerData.TransactionID),
		)
	}

	telemetry.Log().Info("MP payment cancelled",
		slog.String("intent_id", intent.ID.String()),
		slog.String("payment_id", providerData.TransactionID),
		slog.String("status", mpResp.Status),
		slog.String("status_detail", mpResp.StatusDetail),
	)

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
