package mercado_pago_provider

import (
	"context"
	"encoding/json"
	"lib/telemetry"
	"lib/utils"
	"payssage/models"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
	"resty.dev/v3"
)

func (p *Provider) CancelPendingPayment(ctx context.Context, intent *models.Intent) error {
	var providerData models.MercadoPagoIntentData
	err := utils.MapTo(&providerData, intent.ProviderData)
	if err != nil {
		return fun.Errf("error mapping mercadopago provider data: %v", err).Internal()
	}
	if providerData.TransactionID == "" {
		return fun.Err("intent has no mercadopago transaction to cancel").Conflict()
	}

	seller, err := p.sellers.GetByID(ctx, intent.SellerID)
	if err != nil {
		return fun.Errf("resolve seller: %v", err).NotFound()
	}
	if seller.RevokedAt != nil {
		return fun.Errf("seller %s has revoked mercadopago access; cannot cancel — payment may still be live at mercadopago and requires manual review", seller.ID).Conflict()
	}

	var creds models.MercadoPagoCredentials
	err = utils.MapTo(&creds, seller.Credentials)
	if err != nil {
		return fun.Errf("error mapping seller credentials: %v", err).Internal()
	}
	if creds.AccessToken == "" {
		return fun.Errf("seller %s missing mercadopago access token", seller.ID).Internal()
	}

	body := map[string]any{"status": "cancelled"}

	// TODO: hoist into shared *resty.Client (see Checkout TODO)
	client := resty.New()
	defer client.Close()

	var mpResp struct {
		Status       string `json:"status"`
		StatusDetail string `json:"status_detail"`
	}
	var mpErr map[string]any

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+creds.AccessToken).
		SetBody(body).
		SetResult(&mpResp).
		SetResultError(&mpErr).
		Put("https://api.mercadopago.com/v1/payments/" + providerData.TransactionID)
	if err != nil {
		return fun.Errf("mercadopago cancel payment request: %v", err).BadGateway()
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
		return fun.Errf("marshal provider data: %v", err).Internal()
	}
	intent.ProviderData = providerDataBytes
	intent.Status = p.NormalizeStatus(mpResp.Status)
	intent.StatusDetail = p.MapStatusDetail(mpResp.StatusDetail)

	return nil
}
