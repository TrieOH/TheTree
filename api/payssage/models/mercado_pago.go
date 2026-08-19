package models

import "encoding/json"

type MercadoPagoCredentials struct {
	PublicKey      string `json:"public_key"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	ProviderUserID int    `json:"provider_user_id"`
}

type MercadoPagoIntentData struct {
	OrderID                 string `json:"order_id"`
	OrderStatus             string `json:"order_status"`
	OrderStatusDetail       string `json:"order_status_detail"`
	TransactionID           string `json:"transaction_id"`
	TransactionStatus       string `json:"transaction_status"`
	TransactionStatusDetail string `json:"transaction_status_detail"`
	PaymentMethodID         string `json:"payment_method_id"`
	PaymentMethodType       string `json:"payment_method_type"`
	PixQRCode               string `json:"pix_qr_code,omitempty"` // FIXME maybe dont send this or the one below
	PixQRCodeB64            string `json:"pix_qr_code_base64,omitempty"`
	// Fee + settlement observability (refund plan Slice 1): captured from the
	// MP payment response. FeeDetails is the raw fee_details array (entries
	// carry {type, amount, fee_payer} — type application_fee = the
	// marketplace cut, mercadopago_fee = MP's own processing fee).
	// NetReceivedAmountCents is transaction_details.net_received_amount (the
	// seller's net). MoneyReleaseDate/Status answer "where's my money" — the
	// credit appears as a liberar until the release date.
	FeeDetails             json.RawMessage `json:"fee_details,omitempty"`
	NetReceivedAmountCents *int64          `json:"net_received_amount_cents,omitempty"`
	MoneyReleaseDate       *string         `json:"money_release_date,omitempty"`
	MoneyReleaseStatus     *string         `json:"money_release_status,omitempty"`
	// Refund observability (refund plan A2): set by the refund op. The intent
	// status stays succeeded until the payment.refunded webhook confirms;
	// these fields record the MP refund object (id, status, status_detail,
	// amount) for R1 fee-on-refund verification.
	RefundID           *string `json:"refund_id,omitempty"`
	RefundStatus       *string `json:"refund_status,omitempty"`
	RefundStatusDetail *string `json:"refund_status_detail,omitempty"`
	RefundAmountCents  *int64  `json:"refund_amount_cents,omitempty"`
}

// Payer is required by MP, optional for Stripe.
type Payer struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// Document is required in some MP countries (CPF in Brazil).
	DocumentType   *string `json:"document_type"` // "CPF", "CNPJ", etc.
	DocumentNumber *string `json:"document_number"`
}

type MercadoPagoCheckoutData struct {
	Payer                Payer  `json:"payer"`
	Installments         int    `json:"installments"`
	IdentificationNumber string `json:"identification_number"`
	IdentificationType   string `json:"identification_type"`
	MPSellerToken        string `json:"mp_seller_token"`
	MPMarketplaceFeeBPS  int    `json:"mp_marketplace_fee_bps"`
	MPPaymentMethodID    string `json:"mp_payment_method_id"`
	MPPaymentMethodType  string `json:"mp_payment_method_type"`
	MPCardToken          string `json:"mp_card_token"`
}
