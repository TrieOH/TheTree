package mercado_pago_provider

import (
	"cmp"
	"context"
	"encoding/json"
	"payssage/internal/providers/mercado_pago"
	"payssage/models"
	"strconv"
	"time"

	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (p *Provider) Checkout(ctx context.Context, intent *models.Intent, providerCheckoutData json.RawMessage) error {
	checkoutData, err := mercado_pago.ParseCheckoutData(providerCheckoutData)
	if err != nil {
		return err
	}

	seller, err := p.sellers.GetByID(ctx, intent.SellerID)
	if err != nil {
		return fun.Errf("resolve seller: %v", err).NotFound()
	}
	if seller.RevokedAt != nil {
		return fun.Errf("seller %s has revoked provider access", seller.ID).Forbidden()
	}
	if seller.Provider != "mercadopago" {
		return fun.Errf("seller %s is not a mercadopago seller (got %q)", seller.ID, seller.Provider).Conflict()
	}

	var creds models.MercadoPagoCredentials
	err = json.Unmarshal(seller.Credentials, &creds)
	if err != nil {
		return fun.Errf("unmarshal seller credentials: %v", err).Internal()
	}
	if creds.AccessToken == "" {
		return fun.Errf("seller %s missing mercadopago access token", seller.ID).Internal()
	}

	// TODO: check token expiry / refresh via creds.RefreshToken before use
	// once MP token TTL tracking is added to Credentials.

	isPix := checkoutData.PaymentMethodID == "pix"

	body := map[string]any{
		"transaction_amount":   json.Number(formatAmount(intent.AmountCents)),
		"application_fee":      json.Number(formatAmount(calcApplicationFee(intent.AmountCents, checkoutData.MarketplaceFeeBPS))),
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

	if isPix {
		loc := time.FixedZone("BRT", -3*60*60)
		body["date_of_expiration"] = time.Now().In(loc).Add(30 * time.Minute).Format("2006-01-02T15:04:05.000-07:00")
	} else {
		body["installments"] = checkoutData.Installments
		body["token"] = checkoutData.Token
	}

	telemetry.Log().Info("MP Create Payment Request", zap.Any("body", body))

	// TODO remember to use external group feature we made
	// TODO: fix the idemp key not being actual idemp
	idempotencyKey, err := uuid.NewV7()
	if err != nil {
		return fun.Errf("generate idempotency key: %v", err).Internal()
	}

	var mpResp struct {
		ID                 int64  `json:"id"`
		Status             string `json:"status"`
		StatusDetail       string `json:"status_detail"`
		PointOfInteraction struct {
			TransactionData struct {
				QRCode       string `json:"qr_code"`
				QRCodeBase64 string `json:"qr_code_base64"`
			} `json:"transaction_data"`
		} `json:"point_of_interaction"`
	}
	var mpErr map[string]any

	resp, err := p.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+creds.AccessToken).
		SetHeader("X-Idempotency-Key", idempotencyKey.String()).
		SetBody(body).
		SetResult(&mpResp).
		SetResultError(&mpErr).
		Post("https://api.mercadopago.com/v1/payments")
	if err != nil {
		return fun.Errf("mercadopago create payment request: %v", err).BadGateway()
	}

	telemetry.Log().Info("MP Create Payment Request Raw Body", zap.String("body", resp.String()))

	if resp.IsStatusFailure() {
		return mapMPError(resp.StatusCode(), resp.ResultError())
	}

	providerData := &models.MercadoPagoIntentData{
		OrderID:           strconv.FormatInt(mpResp.ID, 10),
		TransactionID:     strconv.FormatInt(mpResp.ID, 10),
		OrderStatus:       mpResp.Status,
		OrderStatusDetail: mpResp.StatusDetail,
		PaymentMethodID:   checkoutData.PaymentMethodID,
	}
	if isPix {
		providerData.PaymentMethodType = "bank_transfer"
		providerData.PixQRCode = mpResp.PointOfInteraction.TransactionData.QRCode
		providerData.PixQRCodeB64 = mpResp.PointOfInteraction.TransactionData.QRCodeBase64
	}

	providerDataBytes, err := json.Marshal(providerData)
	if err != nil {
		return fun.Errf("marshal provider data: %v", err).Internal()
	}
	intent.ProviderData = providerDataBytes
	intent.Status = p.NormalizeStatus(mpResp.Status)
	intent.StatusDetail = p.MapStatusDetail(mpResp.StatusDetail)

	return nil
}
