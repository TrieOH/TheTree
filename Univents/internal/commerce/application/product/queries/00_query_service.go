package queries

import (
	"univents/internal/commerce/domain"
	domain2 "univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/authzed/authzed-go/v1"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	products  domain.ProductsRepository
	purchases domain.PurchaseRepository
	editions  domain2.EditionsRepository
	inventory domain.InventorySubscriber
	gaClient  *goauth.Client
	tracer    trace.Tracer
	az        *authzed.Client
	tx        database.TxRunner
}

func New(
	products domain.ProductsRepository,
	purchases domain.PurchaseRepository,
	editions domain2.EditionsRepository,
	inventory domain.InventorySubscriber,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	az *authzed.Client,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		products:  products,
		purchases: purchases,
		editions:  editions,
		inventory: inventory,
		gaClient:  gaClient,
		tracer:    tracer,
		az:        az,
		tx:        tx,
	}
}
