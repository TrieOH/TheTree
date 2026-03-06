package infrastructure

import (
	"context"
	"time"
	"univents/internal/commerce/domain"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type productsRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ domain.ProductsRepository = (*productsRepo)(nil)

func NewProductsRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) domain.ProductsRepository {
	return &productsRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *productsRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapProductFromDB(src *sqlc.Product) *domain.Product {
	return &domain.Product{
		ID:                 src.ID,
		ScopeID:            src.ScopeID,
		EditionID:          src.EditionID,
		Name:               src.Name,
		Description:        src.Description,
		Type:               domain.ProductType(src.Type),
		TicketID:           src.TicketID,
		PriceCents:         src.PriceCents,
		Status:             domain.ProductStatus(src.Status),
		AvailableFrom:      src.AvailableFrom,
		AvailableUntil:     src.AvailableUntil,
		HasInventory:       src.HasInventory,
		InventoryQuantity:  src.InventoryQuantity,
		InventoryRemaining: src.InventoryRemaining,
		CreatedBy:          src.CreatedBy,
		CreatedAt:          src.CreatedAt,
		UpdatedAt:          src.UpdatedAt,
		DeletedAt:          src.DeletedAt,
	}
}

func (repo *productsRepo) Create(ctx context.Context, toCreate domain.Product) (*domain.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.Create")
	defer span.End()

	sqlcProduct, err := repo.queries(ctx).CreateProduct(ctx, sqlc.CreateProductParams{
		ID:                 toCreate.ID,
		ScopeID:            toCreate.ScopeID,
		EditionID:          toCreate.EditionID,
		Name:               toCreate.Name,
		Description:        toCreate.Description,
		Type:               sqlc.ProductType(toCreate.Type),
		PriceCents:         toCreate.PriceCents,
		AvailableFrom:      toCreate.AvailableFrom,
		AvailableUntil:     toCreate.AvailableUntil,
		HasInventory:       toCreate.HasInventory,
		InventoryQuantity:  toCreate.InventoryQuantity,
		InventoryRemaining: toCreate.InventoryRemaining,
		CreatedBy:          toCreate.CreatedBy,
		TicketID:           toCreate.TicketID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	return mapProductFromDB(&sqlcProduct), nil
}

func (repo *productsRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.GetByID")
	defer span.End()

	sqlcProduct, err := repo.queries(ctx).GetProductByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	return mapProductFromDB(&sqlcProduct), nil
}

func (repo *productsRepo) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.GetByIDs")
	defer span.End()

	sqlcProducts, err := repo.queries(ctx).GetProductsByIDs(ctx, ids)
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	out := make([]domain.Product, 0, len(sqlcProducts))
	for _, product := range sqlcProducts {
		out = append(out, *mapProductFromDB(&product))
	}
	return out, nil
}

func (repo *productsRepo) List(ctx context.Context, editionID uuid.UUID) ([]domain.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.List")
	defer span.End()

	sqlcProducts, err := repo.queries(ctx).ListEditionProducts(ctx, editionID)
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	out := make([]domain.Product, 0, len(sqlcProducts))
	for _, product := range sqlcProducts {
		out = append(out, *mapProductFromDB(&product))
	}
	return out, nil
}

func (repo *productsRepo) AdminList(ctx context.Context, editionID uuid.UUID) ([]domain.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.AdminList")
	defer span.End()

	sqlcProducts, err := repo.queries(ctx).ListEditionProductsAdmin(ctx, editionID)
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	out := make([]domain.Product, 0, len(sqlcProducts))
	for _, product := range sqlcProducts {
		out = append(out, *mapProductFromDB(&product))
	}
	return out, nil
}

func (repo *productsRepo) ReserveItems(ctx context.Context, sessionID uuid.UUID, items []domain.CartItem, expiresAt time.Time) error {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.ReserveItems")
	defer span.End()

	for _, item := range items {
		if item.HasInventory {
			_, err := repo.queries(ctx).ReserveProduct(ctx, sqlc.ReserveProductParams{
				SessionID: sessionID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				ExpiresAt: expiresAt,
			})
			if err != nil {
				return errx.FromDB(err, "product")
			}
		} else {
			err := repo.queries(ctx).ReserveProductNoInventory(ctx, sqlc.ReserveProductNoInventoryParams{
				SessionID: sessionID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				ExpiresAt: expiresAt,
			})
			if err != nil {
				return errx.FromDB(err, "product")
			}
		}
	}

	return nil
}

func (repo *productsRepo) UnreserveItems(ctx context.Context, sessionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.UnreserveItems")
	defer span.End()

	if err := repo.queries(ctx).UnreserveProducts(ctx, sessionID); err != nil {
		return errx.FromDB(err, "product")
	}

	return nil
}

func (repo *productsRepo) DeleteReservation(ctx context.Context, sessionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.DeleteReservation")
	defer span.End()

	if err := repo.queries(ctx).DeleteReservation(ctx, sessionID); err != nil {
		return errx.FromDB(err, "product")
	}

	return nil
}
