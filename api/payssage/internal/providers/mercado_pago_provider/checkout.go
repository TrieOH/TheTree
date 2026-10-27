package mercado_pago_provider

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"payssage/errs"
	"payssage/internal/providers/mercado_pago"
	"payssage/models"
	"strconv"
	"time"

	"lib/telemetry"

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
		return fmt.Errorf("resolve seller: %w", err)
	}
	if seller.RevokedAt != nil {
		return &errs.SellerRevokedError{SellerID: seller.ID, Op: "checkout"}
	}
	if seller.Provider != "mercadopago" {
		return &errs.SellerProviderMismatchError{SellerID: seller.ID, Provider: seller.Provider}
	}

	wallet, err := p.wallets.GetByID(ctx, intent.WalletID)
	if err != nil {
		return fmt.Errorf("resolve wallet: %w", err)
	}

	// The wallet's fee (bps) is the platform's marketplace fee — the
	// authoritative source for application_fee. A caller-supplied
	// marketplace_fee_bps is honored only when the wallet has no fee
	// configured (0), so a caller cannot undercut platform policy.
	feeBps := wallet.FeeBps
	if feeBps == 0 {
		feeBps = checkoutData.MarketplaceFeeBPS
	}

	var creds models.MercadoPagoCredentials
	err = json.Unmarshal(seller.Credentials, &creds)
	if err != nil {
		return fmt.Errorf("unmarshal seller credentials: %w", err)
	}
	if creds.AccessToken == "" {
		return fmt.Errorf("seller %s missing mercadopago access token", seller.ID)
	}

	// TODO: check token expiry / refresh via creds.RefreshToken before use
	// once MP token TTL tracking is added to Credentials.

	isPix := checkoutData.PaymentMethodID == "pix"

	body := map[string]any{
		"transaction_amount":   json.Number(formatAmount(intent.AmountCents)),
		"application_fee":      json.Number(formatAmount(calcApplicationFee(intent.AmountCents, feeBps))),
		"payment_method_id":    checkoutData.PaymentMethodID,
		"external_reference":   intent.ID.String(),
		"statement_descriptor": cmp.Or(checkoutData.StatementDescriptor, "payssage"),
		"description":          paymentDescription(checkoutData.AdditionalInfo),
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
		if checkoutData.IssuerID != "" {
			body["issuer_id"] = checkoutData.IssuerID
		}
	}

	telemetry.Log().Info("MP Create Payment Request", zap.Any("body", body))

	// TODO remember to use external group feature we made
	// TODO: fix the idemp key not being actual idemp
	idempotencyKey, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generate idempotency key: %w", err)
	}

	var mpResp mpPaymentResponse
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
		return fmt.Errorf("%w: mercadopago create payment request: %w", errs.ErrProviderRequestFailed, err)
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
	applySettlement(providerData, &mpResp)
	if isPix {
		providerData.PaymentMethodType = "bank_transfer"
		providerData.PixQRCode = mpResp.PointOfInteraction.TransactionData.QRCode
		providerData.PixQRCodeB64 = mpResp.PointOfInteraction.TransactionData.QRCodeBase64
	}

	providerDataBytes, err := json.Marshal(providerData)
	if err != nil {
		return fmt.Errorf("marshal provider data: %w", err)
	}
	intent.ProviderData = providerDataBytes
	intent.Status = p.NormalizeStatus(mpResp.Status)
	intent.StatusDetail = p.MapStatusDetail(mpResp.StatusDetail)

	return nil
}
