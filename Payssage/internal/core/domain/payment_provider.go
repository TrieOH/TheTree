package domain

import (
	"context"

	"github.com/google/uuid"
)

type ChargeRequest struct {
	Amount          float64
	CardToken       string
	PaymentMethodID string
	Installments    int
	PayerEmail      string
	ApplicationFee  float64
	SellerToken     string
	IntentID        uuid.UUID
	OrderID         string
}

type ChargeResult struct {
	OrderID string
	Status  IntentStatus
}

type PaymentProvider interface {
	Charge(ctx context.Context, req ChargeRequest) (*ChargeResult, error)
}
