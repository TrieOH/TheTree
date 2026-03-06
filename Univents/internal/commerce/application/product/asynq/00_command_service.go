package async

import (
	"univents/internal/commerce/domain"
	"univents/internal/payments"
	"univents/internal/plataform/database"
	"univents/internal/shared/sockets"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type AsynqHandlers struct {
	products  domain.ProductsRepository
	purchases domain.PurchaseRepository
	payments  *payments.MockPayments
	ws        *sockets.Registry
	gaClient  *goauth.Client
	tracer    trace.Tracer
	tx        database.TxRunner
}

func New(
	products domain.ProductsRepository,
	purchases domain.PurchaseRepository,
	ws *sockets.Registry,
	payments *payments.MockPayments,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *AsynqHandlers {
	return &AsynqHandlers{
		products:  products,
		purchases: purchases,
		payments:  payments,
		ws:        ws,
		gaClient:  gaClient,
		tracer:    tracer,
		tx:        tx,
	}
}
