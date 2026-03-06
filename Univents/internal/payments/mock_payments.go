package payments

import (
	"context"
	"log"
	"univents/internal/commerce/interfaces/http/dtos"

	"github.com/google/uuid"
)

type MockPayments struct{}

func (m *MockPayments) CreatePaymentIntent(ctx context.Context, req dtos.BuyRequest) (string, string, string, error) {
	intentID := "pi_mock_" + uuid.NewString()
	secret := "secret_mock_" + uuid.NewString()
	log.Printf("[mock stripe] created payment intent %s", intentID)
	return intentID, secret, "mercado pago", nil
}

func (m *MockPayments) CancelPaymentIntent(ctx context.Context, paymentIntentID string) error {
	log.Printf("[mock stripe] cancelled payment intent %s", paymentIntentID)
	return nil
}
