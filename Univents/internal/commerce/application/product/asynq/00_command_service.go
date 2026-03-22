package async

import (
	"univents/internal/commerce/domain"
	"univents/internal/plataform/database"
	"univents/internal/shared/sockets"

	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type AsynqHandlers struct {
	products  domain.ProductsRepository
	purchases domain.PurchaseRepository
	payments  *paymentsSDK.Client
	inventory domain.InventoryPublisher
	sessions  domain.PurchaseSessionStore
	ws        *sockets.Registry
	gaClient  *goauth.Client
	tracer    trace.Tracer
	tx        database.TxRunner
}

func New(
	products domain.ProductsRepository,
	purchases domain.PurchaseRepository,
	payments *paymentsSDK.Client,
	inventory domain.InventoryPublisher,
	sessions domain.PurchaseSessionStore,
	ws *sockets.Registry,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *AsynqHandlers {
	return &AsynqHandlers{
		products:  products,
		purchases: purchases,
		payments:  payments,
		inventory: inventory,
		sessions:  sessions,
		ws:        ws,
		gaClient:  gaClient,
		tracer:    tracer,
		tx:        tx,
	}
}
