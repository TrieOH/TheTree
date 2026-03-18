package queries

import (
	"univents/internal/commerce/domain"
	domain2 "univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	products  domain.ProductsRepository
	purchases domain.PurchaseRepository
	editions  domain2.EditionsRepository
	inventory domain.InventorySubscriber
	gaClient  *goauth.Client
	tracer    trace.Tracer
	tx        database.TxRunner
}

func New(
	products domain.ProductsRepository,
	purchases domain.PurchaseRepository,
	editions domain2.EditionsRepository,
	inventory domain.InventorySubscriber,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		products:  products,
		purchases: purchases,
		editions:  editions,
		inventory: inventory,
		gaClient:  gaClient,
		tracer:    tracer,
		tx:        tx,
	}
}
