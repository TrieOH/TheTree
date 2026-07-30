package mercado_pago

import (
	"encoding/json"
	"fmt"
	"lib/jsonschema"
)

// CheckoutData is the MercadoPago-specific payload carried in
// CreateIntentRequest.ProviderData. Validate it against ProviderDataSchema
// (via jsonschema.Validate) before unmarshaling — this struct assumes
// the instance already conforms to the schema.
type CheckoutData struct {
	Installments        int             `json:"installments"`
	Token               string          `json:"token"`
	PaymentMethodID     string          `json:"payment_method_id"`
	MarketplaceFeeBPS   int             `json:"marketplace_fee_bps"`
	StatementDescriptor string          `json:"statement_descriptor,omitempty"`
	Payer               Payer           `json:"payer"`
	AdditionalInfo      *AdditionalInfo `json:"additional_info,omitempty"`
}

type Payer struct {
	Email                string `json:"email"`
	IdentificationType   string `json:"identification_type"`
	IdentificationNumber string `json:"identification_number"`
}

type AdditionalInfo struct {
	Items []Item `json:"items"`
}

type Item struct {
	Title          string `json:"title"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents *int64 `json:"unit_price_cents,omitempty"` // nil -> fall back to intent's amount_cents
}

// ParseCheckoutData validates raw against CheckoutDataSchema, then unmarshals
// it into a CheckoutData struct. Always use this instead of unmarshaling
// provider_data directly.
func ParseCheckoutData(raw json.RawMessage) (*CheckoutData, error) {
	err := jsonschema.Validate(CheckoutDataSchema, raw)
	if err != nil {
		return nil, fmt.Errorf("mercadopago checkout_data: %w", err)
	}

	var cd CheckoutData
	err := json.Unmarshal(raw, &cd)
	if err != nil {
		return nil, fmt.Errorf("mercadopago checkout_data: unmarshal: %w", err)
	}
	return &cd, nil
}
