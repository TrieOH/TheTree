package mercado_pago_provider

import (
	"context"
	"encoding/json"
	"lib/utils"
	"net/http"
	"payssage/models"

	"github.com/MintzyG/fun"
)

func (p *Provider) Parse(ctx context.Context, _ *http.Request, rawBody []byte) (*models.WebhookParseResult, error) {
	var envelope struct {
		Type   string `json:"type"`
		Action string `json:"action"`
		Data   struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rawBody, &envelope); err != nil {
		return nil, fun.Errf("mercadopago webhook: decode envelope: %v", err).BadRequest()
	}

	// Only "payment" events are handled today. Other MP topics
	// (merchant_order, chargebacks, etc.) are acknowledged but ignored.
	if envelope.Type != "payment" {
		return nil, fun.Errf("mercadopago webhook: unhandled event type %q", envelope.Type).NotFound()
	}
	if envelope.Data.ID == "" {
		return nil, fun.Err("mercadopago webhook: missing data.id").BadRequest()
	}

	intent, err := p.intents.GetByProviderTransactionID(ctx, "mercadopago", envelope.Data.ID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil, fun.Errf("mercadopago webhook: no intent found for payment %s", envelope.Data.ID).NotFound()
		}
		return nil, fun.Errf("resolve intent for payment %s: %v", envelope.Data.ID, err).Internal()
	}

	seller, err := p.sellers.GetByID(ctx, intent.SellerID)
	if err != nil {
		return nil, fun.Errf("resolve seller: %v", err).NotFound()
	}

	var creds models.MercadoPagoCredentials
	err = utils.MapTo(&creds, seller.Credentials)
	if err != nil {
		return nil, fun.Errf("unmarshal seller credentials: %v", err).Internal()
	}
	if creds.AccessToken == "" {
		return nil, fun.Errf("seller %s has no mercadopago access token; cannot resolve webhook for intent %s", seller.ID, intent.ID).Conflict()
	}

	var mpResp struct {
		ID           int64  `json:"id"`
		Status       string `json:"status"`
		StatusDetail string `json:"status_detail"`
	}
	var mpErr map[string]any

	resp, err := p.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+creds.AccessToken).
		SetResult(&mpResp).
		SetResultError(&mpErr).
		Get("https://api.mercadopago.com/v1/payments/" + envelope.Data.ID)
	if err != nil {
		return nil, fun.Errf("mercadopago fetch payment request: %v", err).BadGateway()
	}
	if resp.IsStatusFailure() {
		return nil, mapMPError(resp.StatusCode(), resp.ResultError())
	}

	intent.Status = p.NormalizeStatus(mpResp.Status)
	intent.StatusDetail = p.MapStatusDetail(mpResp.StatusDetail)

	var providerData models.MercadoPagoIntentData
	_ = utils.MapTo(&providerData, intent.ProviderData)
	providerData.OrderStatus = mpResp.Status
	providerData.OrderStatusDetail = mpResp.StatusDetail
	providerDataBytes, err := json.Marshal(&providerData)
	if err != nil {
		return nil, fun.Errf("marshal provider data: %v", err).Internal()
	}
	intent.ProviderData = providerDataBytes

	updatedIntent, err := p.intents.Update(ctx, *intent)
	if err != nil {
		return nil, fun.Errf("persist updated intent: %v", err).Internal()
	}

	return &models.WebhookParseResult{
		WalletID:   updatedIntent.WalletID,
		IntentID:   updatedIntent.ID,
		ExternalID: envelope.Data.ID,
		EventType:  "payment." + string(updatedIntent.Status),
	}, nil
}
