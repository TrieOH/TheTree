package async

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type AsynqHandlers struct {
	deliveries domain.WebhookDeliveryRepo
	gaClient   *goauth.Client
	tracer     trace.Tracer
	tx         database.TxRunner
}

func New(
	deliveries domain.WebhookDeliveryRepo,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *AsynqHandlers {
	return &AsynqHandlers{
		deliveries: deliveries,
		gaClient:   gaClient,
		tracer:     tracer,
		tx:         tx,
	}
}
