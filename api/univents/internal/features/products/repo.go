package products

import (
	"context"
	"errors"
	"time"

	"univents/internal/platform/database"
	"univents/internal/platform/database/sqlc"
	"univents/internal/shared/contracts"
	"univents/internal/shared/errx"
	"univents/internal/shared/ports"

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

var _ ports.ProductsRepository = (*productsRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ProductsRepository {
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

func mapProductFromDB(src *sqlc.Product) *contracts.Product {
	return &contracts.Product{
		ID:                 src.ID,
		ScopeID:            src.ScopeID,
		EditionID:          src.EditionID,
		Name:               src.Name,
		Description:        src.Description,
		Type:               contracts.ProductType(src.Type),
		TicketID:           src.TicketID,
		PriceCents:         src.PriceCents,
		Status:             contracts.ProductStatus(src.Status),
		AvailableFrom:      src.AvailableFrom,
		AvailableUntil:     src.AvailableUntil,
		HasInventory:       src.HasInventory,
		InventoryQuantity:  src.InventoryQuantity,
		InventoryRemaining: src.InventoryRemaining,
		ThumbnailURL:       src.ThumbnailUrl,
		GalleryURLs:        src.GalleryUrls,
		CreatedBy:          src.CreatedBy,
		CreatedAt:          src.CreatedAt,
		UpdatedAt:          src.UpdatedAt,
		DeletedAt:          src.DeletedAt,
	}
}

func (repo *productsRepo) Create(ctx context.Context, toCreate contracts.Product) (*contracts.Product, error) {
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

func (repo *productsRepo) Publish(ctx context.Context, id uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.Publish")
	defer span.End()

	err := repo.queries(ctx).PublishProduct(ctx, id)
	if err != nil {
		return errx.FromDB(err, "product")
	}

	return nil
}

func (repo *productsRepo) GetByID(ctx context.Context, id uuid.UUID) (*contracts.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.GetByID")
	defer span.End()

	sqlcProduct, err := repo.queries(ctx).GetProductByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	return mapProductFromDB(&sqlcProduct), nil
}

func (repo *productsRepo) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]contracts.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.GetByIDs")
	defer span.End()

	sqlcProducts, err := repo.queries(ctx).GetProductsByIDs(ctx, ids)
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	out := make([]contracts.Product, 0, len(sqlcProducts))
	for _, product := range sqlcProducts {
		out = append(out, *mapProductFromDB(&product))
	}
	return out, nil
}

func (repo *productsRepo) List(ctx context.Context, editionID uuid.UUID) ([]contracts.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.List")
	defer span.End()

	sqlcProducts, err := repo.queries(ctx).ListEditionProducts(ctx, editionID)
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	out := make([]contracts.Product, 0, len(sqlcProducts))
	for _, product := range sqlcProducts {
		out = append(out, *mapProductFromDB(&product))
	}
	return out, nil
}

func (repo *productsRepo) AdminList(ctx context.Context, editionID uuid.UUID) ([]contracts.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.AdminList")
	defer span.End()

	sqlcProducts, err := repo.queries(ctx).ListEditionProductsAdmin(ctx, editionID)
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	out := make([]contracts.Product, 0, len(sqlcProducts))
	for _, product := range sqlcProducts {
		out = append(out, *mapProductFromDB(&product))
	}
	return out, nil
}

func (repo *productsRepo) Delete(ctx context.Context, productID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.Delete")
	defer span.End()

	if err := repo.queries(ctx).SoftDeleteProduct(ctx, productID); err != nil {
		return errx.FromDB(err, "product")
	}

	return nil
}

func (repo *productsRepo) Restore(ctx context.Context, productID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.Restore")
	defer span.End()

	if err := repo.queries(ctx).RestoreProduct(ctx, productID); err != nil {
		return errx.FromDB(err, "product")
	}

	return nil
}

func (repo *productsRepo) ItemHasCompletedPurchases(ctx context.Context, productID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.ItemHasCompletedPurchases")
	defer span.End()

	has, err := repo.queries(ctx).ItemHasCompletedPurchases(ctx, productID)
	if err != nil {
		return true, errx.FromDB(err, "product")
	}

	return has, nil
}

func (repo *productsRepo) ReserveItems(ctx context.Context, sessionID uuid.UUID, items []contracts.CartItem, expiresAt time.Time) (contracts.ReservationOutcome, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.ReserveItems")
	defer span.End()

	var outcome contracts.ReservationOutcome

	for _, item := range items {
		if !item.HasInventory {
			err := repo.queries(ctx).ReserveProductNoInventory(ctx, sqlc.ReserveProductNoInventoryParams{
				SessionID: sessionID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				ExpiresAt: expiresAt,
			})
			if err != nil {
				_, _ = repo.queries(ctx).UnreserveProducts(ctx, sessionID)
				return contracts.ReservationOutcome{}, errx.FromDB(err, "product")
			}
			outcome.Reserved = append(outcome.Reserved, item)
			continue
		}

		row, err := repo.queries(ctx).ReserveProduct(ctx, sqlc.ReserveProductParams{
			SessionID: sessionID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			ExpiresAt: expiresAt,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				outcome.Unavailable = append(outcome.Unavailable, contracts.InvalidProduct{
					ProductID: item.ProductID,
					Requested: item.Quantity,
					Reserved:  0,
					Reason:    "out_of_stock",
				})
				continue
			}
			_, _ = repo.queries(ctx).UnreserveProducts(ctx, sessionID)
			return contracts.ReservationOutcome{}, errx.FromDB(err, "product")
		}

		outcome.InventoryUpdates = append(outcome.InventoryUpdates, contracts.InventoryUpdate{
			ProductID:          item.ProductID,
			InventoryRemaining: int(row.InventoryRemaining),
		})

		if row.ReservedQuantity < item.Quantity {
			outcome.Reserved = append(outcome.Reserved, contracts.CartItem{
				ProductID:    item.ProductID,
				Quantity:     int(row.ReservedQuantity),
				HasInventory: item.HasInventory,
			})
			outcome.Unavailable = append(outcome.Unavailable, contracts.InvalidProduct{
				ProductID: item.ProductID,
				Requested: item.Quantity,
				Reserved:  int(row.ReservedQuantity),
				Reason:    "insufficient_inventory",
			})
			continue
		}

		outcome.Reserved = append(outcome.Reserved, item)
	}

	return outcome, nil
}

func (repo *productsRepo) UnreserveItems(ctx context.Context, sessionID uuid.UUID) ([]contracts.InventoryUpdate, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.UnreserveItems")
	defer span.End()

	rows, err := repo.queries(ctx).UnreserveProducts(ctx, sessionID)
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	updates := make([]contracts.InventoryUpdate, 0, len(rows))
	for _, row := range rows {
		updates = append(updates, contracts.InventoryUpdate{
			ProductID:          row.ID,
			InventoryRemaining: row.InventoryRemaining,
		})
	}

	return updates, nil
}

func (repo *productsRepo) DeleteReservation(ctx context.Context, sessionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.DeleteReservation")
	defer span.End()

	if err := repo.queries(ctx).DeleteReservation(ctx, sessionID); err != nil {
		return errx.FromDB(err, "product")
	}

	return nil
}

func (repo *productsRepo) AddGalleryImage(ctx context.Context, id uuid.UUID, url string) (*contracts.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.AddGalleryImage")
	defer span.End()

	sqlcProduct, err := repo.queries(ctx).AddProductGalleryImage(ctx, sqlc.AddProductGalleryImageParams{
		ID:  id,
		Url: url,
	})
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	return mapProductFromDB(&sqlcProduct), nil
}

func (repo *productsRepo) RemoveGalleryImage(ctx context.Context, id uuid.UUID, url string) (*contracts.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.RemoveGalleryImage")
	defer span.End()

	sqlcProduct, err := repo.queries(ctx).RemoveProductGalleryImage(ctx, sqlc.RemoveProductGalleryImageParams{
		ID:  id,
		Url: url,
	})
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	return mapProductFromDB(&sqlcProduct), nil
}

func (repo *productsRepo) SetThumbnail(ctx context.Context, id uuid.UUID, url string) (*contracts.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.SetThumbnail")
	defer span.End()

	sqlcProduct, err := repo.queries(ctx).SetThumbnail(ctx, sqlc.SetThumbnailParams{
		ID:  id,
		Url: url,
	})
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	return mapProductFromDB(&sqlcProduct), nil
}

func (repo *productsRepo) UnsetThumbnail(ctx context.Context, id uuid.UUID) (*contracts.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.UnsetThumbnail")
	defer span.End()

	sqlcProduct, err := repo.queries(ctx).UnsetThumbnail(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "product")
	}

	return mapProductFromDB(&sqlcProduct), nil
}
