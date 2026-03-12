package domain

import "context"

type ChargeRequest struct {
	Intent          Intent
	CardToken       string
	PaymentMethodID string
	Installments    int
	PayerEmail      string
	ApplicationFee  float64
	SellerToken     string
	SponsorID       int
}

type ChargeResult struct {
	ProviderPaymentID string
	Status            IntentStatus
}

type PaymentProvider interface {
	Charge(ctx context.Context, req ChargeRequest) (*ChargeResult, error)
}
