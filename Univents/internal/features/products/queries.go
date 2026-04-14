package products

import (
	"context"
	"univents/internal/platform/database"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	domain2 "univents/internal/shared/ports"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	products  domain2.ProductsRepository
	purchases domain2.PurchaseRepository
	editions  domain2.EditionsRepository
	inventory domain2.InventorySubscriber
	gaClient  *goauth.Client
	tracer    trace.Tracer
	az        *authzed.Client
	tx        database.TxRunner
}

func NewQueryService(
	products domain2.ProductsRepository,
	purchases domain2.PurchaseRepository,
	editions domain2.EditionsRepository,
	inventory domain2.InventorySubscriber,
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

func (uc *QueryService) StreamInventory(ctx context.Context, editionID uuid.UUID) (<-chan []contracts.InventoryUpdate, error) {
	return uc.inventory.Subscribe(ctx, editionID)
}

func (uc *QueryService) List(ctx context.Context, editionID uuid.UUID) (out []contracts.Product, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "ProductService.List")
	defer span.End()

	return uc.products.List(ctx, editionID)
}

func (uc *QueryService) AdminList(ctx context.Context, editionID uuid.UUID) (out []contracts.Product, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "ProductService.AdminList")
	defer span.End()

	edition, err := uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("view_products"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return nil, err
	}

	return uc.products.AdminList(ctx, editionID)
}
