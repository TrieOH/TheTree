package mercado_pago_provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lib/database"
	"net/http"
	"payssage/errs"
	"payssage/models"
)

func (p *Provider) Parse(ctx context.Context, _ *http.Request, rawBody []byte) (*models.WebhookParseResult, error) {
	var envelope struct {
		Type   string `json:"type"`
		Action string `json:"action"`
		Data   struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err := json.Unmarshal(rawBody, &envelope)
	if err != nil {
		return nil, fmt.Errorf("%w: decode envelope: %w", errs.ErrWebhookPayloadInvalid, err)
	}

	// Only "payment" events are handled today. Other MP topics
	// (merchant_order, chargebacks, etc.) are acknowledged but ignored.
	if envelope.Type != "payment" {
		return nil, &errs.UnhandledWebhookEventError{Type: envelope.Type}
	}
	if envelope.Data.ID == "" {
		return nil, errs.ErrWebhookMissingDataID
	}

	intent, err := p.intents.GetByProviderTransactionID(ctx, "mercadopago", envelope.Data.ID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, &errs.WebhookIntentNotFoundError{PaymentID: envelope.Data.ID}
		}
		return nil, fmt.Errorf("resolve intent for payment %s: %w", envelope.Data.ID, err)
	}

	seller, err := p.sellers.GetByID(ctx, intent.SellerID)
	if err != nil {
		return nil, fmt.Errorf("resolve seller: %w", err)
	}

	var creds models.MercadoPagoCredentials
	err = json.Unmarshal(seller.Credentials, &creds)
	if err != nil {
		return nil, fmt.Errorf("unmarshal seller credentials: %w", err)
	}
	if creds.AccessToken == "" {
		return nil, &errs.WebhookSellerTokenMissingError{SellerID: seller.ID, IntentID: intent.ID}
	}

	var mpResp mpPaymentResponse
	var mpErr map[string]any

	resp, err := p.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+creds.AccessToken).
		SetResult(&mpResp).
		SetResultError(&mpErr).
		Get("https://api.mercadopago.com/v1/payments/" + envelope.Data.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: mercadopago fetch payment request: %w", errs.ErrProviderRequestFailed, err)
	}
	if resp.IsStatusFailure() {
		return nil, mapMPError(resp.StatusCode(), resp.ResultError())
	}

	intent.Status = p.NormalizeStatus(mpResp.Status)
	intent.StatusDetail = p.MapStatusDetail(mpResp.StatusDetail)

	var providerData models.MercadoPagoIntentData
	_ = json.Unmarshal(intent.ProviderData, &providerData)
	providerData.OrderStatus = mpResp.Status
	providerData.OrderStatusDetail = mpResp.StatusDetail
	applySettlement(&providerData, &mpResp)
	providerDataBytes, err := json.Marshal(&providerData)
	if err != nil {
		return nil, fmt.Errorf("marshal provider data: %w", err)
	}
	intent.ProviderData = providerDataBytes

	updatedIntent, err := p.intents.Update(ctx, *intent)
	if err != nil {
		return nil, fmt.Errorf("persist updated intent: %w", err)
	}

	return &models.WebhookParseResult{
		WalletID:     updatedIntent.WalletID,
		IntentID:     updatedIntent.ID,
		ExternalID:   envelope.Data.ID,
		EventType:    "payment." + string(updatedIntent.Status),
		StatusDetail: (*string)(updatedIntent.StatusDetail),
	}, nil
}
