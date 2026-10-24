package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetVariantByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.ProductVariant, error) {
	result, err := database.Queries(ctx, repo.q).GetProductVariantByVendorCode(ctx, sqlc.GetProductVariantByVendorCodeParams{
		EditionID:  editionID,
		VendorCode: vendorCode,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapVariant(result)), nil
}
