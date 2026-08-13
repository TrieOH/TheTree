package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type ProductRepo interface {
	CreateProduct(ctx context.Context, toCreate *models.Product) (*models.Product, error)
	CreateVariant(ctx context.Context, toCreate *models.ProductVariant) (*models.ProductVariant, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProductByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.Product, error)
	ListProductsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Product, error)
	GetVariantByID(ctx context.Context, id uuid.UUID) (*models.ProductVariant, error)
	GetVariantByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.ProductVariant, error)
	GetVariantByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.ProductVariant, error)
	ListVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]models.ProductVariant, error)
	PatchProduct(ctx context.Context, id uuid.UUID, product *models.Product) (*models.Product, error)
	PatchVariant(ctx context.Context, id uuid.UUID, variant *models.ProductVariant) (*models.ProductVariant, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	DeleteVariant(ctx context.Context, id uuid.UUID) error
}
