package mercado_pago_provider

import (
	"encoding/json"
	"math"

	"payssage/models"
)

// mpPaymentResponse is the slice of the MercadoPago /v1/payments response
// payssage consumes — create, get, and refund all return this shape. Kept in
// one place so checkout (POST), webhook re-sync (GET), and refund (POST
// refunds) read identical fields.
type mpPaymentResponse struct {
	ID                 int64           `json:"id"`
	Status             string          `json:"status"`
	StatusDetail       string          `json:"status_detail"`
	FeeDetails         json.RawMessage `json:"fee_details"`
	MoneyReleaseDate   string          `json:"money_release_date"`
	MoneyReleaseStatus string          `json:"money_release_status"`
	TransactionDetails struct {
		NetReceivedAmount float64 `json:"net_received_amount"`
	} `json:"transaction_details"`
	PointOfInteraction struct {
		TransactionData struct {
			QRCode       string `json:"qr_code"`
			QRCodeBase64 string `json:"qr_code_base64"`
		} `json:"transaction_data"`
	} `json:"point_of_interaction"`
}

// applySettlement captures MP's fee + settlement fields onto the intent's
// provider data (refund plan Slice 1 — fee observability): the raw
// fee_details array (entries {type, amount, fee_payer}; type
// "application_fee" = the marketplace cut), the seller's net received
// amount, and the money-release date/status (answers "where's my money" —
// the credit sits a liberar until the release date). Absent fields stay nil.
func applySettlement(providerData *models.MercadoPagoIntentData, resp *mpPaymentResponse) {
	if len(resp.FeeDetails) > 0 && string(resp.FeeDetails) != "null" {
		providerData.FeeDetails = resp.FeeDetails
	}
	if resp.TransactionDetails.NetReceivedAmount != 0 {
		cents := int64(math.Round(resp.TransactionDetails.NetReceivedAmount * 100))
		providerData.NetReceivedAmountCents = &cents
	}
	if resp.MoneyReleaseDate != "" {
		d := resp.MoneyReleaseDate
		providerData.MoneyReleaseDate = &d
	}
	if resp.MoneyReleaseStatus != "" {
		s := resp.MoneyReleaseStatus
		providerData.MoneyReleaseStatus = &s
	}
}
