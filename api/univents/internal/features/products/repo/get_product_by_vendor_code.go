package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetProductByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.Product, error) {
	result, err := database.Queries(ctx, repo.q).GetProductByVendorCode(ctx, sqlc.GetProductByVendorCodeParams{
		EditionID:  editionID,
		VendorCode: vendorCode,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProduct(result)), nil
}
