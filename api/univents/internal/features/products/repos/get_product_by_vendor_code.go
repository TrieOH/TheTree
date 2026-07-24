package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) GetProductByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.GetProductByVendorCode")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetProductByVendorCode(ctx, sqlc.GetProductByVendorCodeParams{
		EditionID:  editionID,
		VendorCode: vendorCode,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProduct(result)), nil
}
