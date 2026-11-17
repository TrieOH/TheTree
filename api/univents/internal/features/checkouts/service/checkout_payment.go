package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lib/telemetry"
	"log/slog"
	payssage "sdk/payssage"
	"univents/errs"
	"univents/models"
	"uuid"
)

// attachIntent stores the intent on the purchase in a second, short
// transaction (post-commit): seller, intent id (the D2 correlation key),
// and the pix QR from the intent's provider_data.
func (o *Operations) attachIntent(ctx context.Context, res *reservation, intent *payssage.Intent) error {
	qr, qrB64 := pixQR(intent.ProviderData)
	return o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.Purchases.AttachIntent(ctx, res.purchase.ID, intent.SellerID, intent.ID, qr, qrB64)
		if err != nil {
			return err
		}
		res.purchase = updated
		return nil
	})
}

// cancelIntentBestEffort cancels a stray intent (Tx 2 failure): pix intents
// are cancellable; a card already charged comes back non-cancellable and is
// left for the refund follow-up (D1). Log-only.
func (o *Operations) cancelIntentBestEffort(ctx context.Context, intentID uuid.UUID) {
	_, err := o.Payssage.CancelIntent(ctx, intentID)
	if err != nil {
		telemetry.Log().Warn("checkout: best-effort intent cancel failed",
			slog.String("intent_id", intentID.String()),
			telemetry.Err(err))
	}
}

// buildIntentRequest maps the purchase + checkout input to the Payssage
// create-intent payload. For pix the server derives `payment_method_id` =
// "pix" (no tokenization exists); for cards the front's tokenization
// supplies it. `additional_info.items` carries the real product lines (name,
// quantity, unit price) so the provider's payment screen and risk engine
// see what the buyer is paying for — an intent without item data surfaces
// as a nameless product and can be flagged high-risk. Metadata carries
// purchase_id + edition_id for correlation.
func (o *Operations) buildIntentRequest(purchase *models.Purchase, in CheckoutInput, sellerID uuid.UUID, priced []pricedLine) payssage.CreateIntentRequest {
	providerData := map[string]any{
		"payer": map[string]any{
			"email":                 in.Payer.Email,
			"identification_type":   in.Payer.IdentificationType,
			"identification_number": in.Payer.IdentificationNumber,
		},
	}
	switch in.PaymentMethod {
	case paymentPix:
		providerData["payment_method_id"] = paymentPix
	case paymentCreditCard, paymentDebitCard:
		installments := 1
		if in.Installments != nil {
			installments = *in.Installments
		}
		providerData["payment_method_id"] = *in.PaymentMethodID
		providerData["token"] = *in.CardToken
		providerData["installments"] = installments
		if in.IssuerID != nil && *in.IssuerID != "" {
			providerData["issuer_id"] = *in.IssuerID
		}
	}

	if len(priced) > 0 {
		items := make([]map[string]any, 0, len(priced))
		for _, pl := range priced {
			items = append(items, map[string]any{
				"title":            pl.Name,
				"quantity":         pl.Line.Quantity,
				"unit_price_cents": pl.UnitPrice,
			})
		}
		providerData["additional_info"] = map[string]any{"items": items}
	}

	return payssage.CreateIntentRequest{
		SellerID:             sellerID,
		Currency:             currencyBRL,
		AmountCents:          purchase.TotalCents,
		CheckoutProviderData: providerData,
		Metadata: map[string]any{
			"purchase_id": purchase.ID.String(),
			"edition_id":  purchase.EditionID.String(),
		},
		// First-class correlation (refund plan A7): the caller ids payssage
		// stores as external_id (purchase) + external_group (edition). The
		// metadata copy is kept for backward compat.
		ExternalID:    new(purchase.ID.String()),
		ExternalGroup: new(purchase.EditionID.String()),
	}
}

// pixQR extracts the pix QR from the intent's provider_data
// (`pix_qr_code` / `pix_qr_code_base64` — payssage's MercadoPagoIntentData).
// Cards carry none — both stay nil.
func pixQR(providerData json.RawMessage) (qr, qrBase64 *string) {
	var d struct {
		PixQRCode       string `json:"pix_qr_code"`
		PixQRCodeBase64 string `json:"pix_qr_code_base64"`
	}
	if len(providerData) == 0 || json.Unmarshal(providerData, &d) != nil {
		return nil, nil
	}
	if d.PixQRCode != "" {
		qr = &d.PixQRCode
	}
	if d.PixQRCodeBase64 != "" {
		qrBase64 = &d.PixQRCodeBase64
	}
	return qr, qrBase64
}

// mapPayssageError converts a provider/API failure into an error the buyer
// can act on: 4xx (rejected card, invalid token...) → a PaymentRejectedError
// with the provider's message; anything else stays an internal error. Without
// this the strict server would turn every payssage failure into a 500.
func mapPayssageError(err error) error {
	var apiErr *payssage.APIError
	if !errors.As(err, &apiErr) {
		return err
	}
	if apiErr.StatusCode >= 400 && apiErr.StatusCode < 500 {
		return &errs.PaymentRejectedError{Message: apiErr.Message}
	}
	return fmt.Errorf("payment provider error: %w", err)
}
