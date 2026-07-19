package models

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
	PixQRCode               string `json:"pix_qr_code,omitempty"` //FIXME maybe dont send this or the one below
	PixQRCodeB64            string `json:"pix_qr_code_base64,omitempty"`
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
