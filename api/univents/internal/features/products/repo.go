package products

import (
	"context"
	"lib/database"
	"univents/errs"
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ProductRepo = (*Repo)(nil)
var _ ports.ProductPurchaseRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrProductNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapProduct(src sqlc.Product) models.Product {
	return models.Product{
		ID:                   src.ID,
		EditionID:            src.EditionID,
		VendorCode:           src.VendorCode,
		RequiresRegistration: src.RequiresRegistration,
		CreatedAt:            src.CreatedAt,
		UpdatedAt:            src.UpdatedAt,
		DeletedAt:            src.DeletedAt,
	}
}

func mapVariant(src sqlc.ProductVariant) models.ProductVariant {
	return models.ProductVariant{
		ID:          src.ID,
		EditionID:   src.EditionID,
		ProductID:   src.ProductID,
		VendorCode:  src.VendorCode,
		Name:        src.Name,
		Description: src.Description,
		Price:       src.Price,
		Stock:       src.Stock,
		GalleryURLs: src.GalleryUrls,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		DeletedAt:   src.DeletedAt,
	}
}

func mapProductPurchase(src sqlc.ProductPurchase) models.ProductPurchase {
	return models.ProductPurchase{
		ID:               src.ID,
		EditionID:        src.EditionID,
		VariantID:        src.VariantID,
		PurchaserID:      src.PurchaserID,
		RecipientID:      src.RecipientID,
		RegistrationID:   src.RegistrationID,
		Quantity:         src.Quantity,
		Status:           models.ProductPurchaseStatus(src.Status),
		StatusReason:     src.StatusReason,
		PayssageIntentID: src.PayssageIntentID,
		CreatedAt:        src.CreatedAt,
		UpdatedAt:        src.UpdatedAt,
		DeletedAt:        src.DeletedAt,
	}
}

func (repo *Repo) CreateProduct(ctx context.Context, toCreate *models.Product) (*models.Product, error) {
	result, err := repo.db(ctx).CreateProduct(ctx, sqlc.CreateProductParams{
		EditionID:            toCreate.EditionID,
		VendorCode:           toCreate.VendorCode,
		RequiresRegistration: toCreate.RequiresRegistration,
	})
	return database.One(result, err, repo.dbe, mapProduct)
}

func (repo *Repo) CreateVariant(ctx context.Context, toCreate *models.ProductVariant) (*models.ProductVariant, error) {
	result, err := repo.db(ctx).CreateProductVariant(ctx, sqlc.CreateProductVariantParams{
		EditionID:   toCreate.EditionID,
		ProductID:   toCreate.ProductID,
		VendorCode:  toCreate.VendorCode,
		Name:        toCreate.Name,
		Description: toCreate.Description,
		Price:       toCreate.Price,
		Stock:       toCreate.Stock,
	})
	return database.One(result, err, repo.dbe, mapVariant)
}

func (repo *Repo) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).DeleteProduct(ctx, id)
	return repo.dbe(err)
}

func (repo *Repo) DeleteVariant(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).DeleteProductVariant(ctx, id)
	return repo.dbe(err)
}

func (repo *Repo) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	result, err := repo.db(ctx).GetProductByID(ctx, id)
	return database.One(result, err, repo.dbe, mapProduct)
}

func (repo *Repo) GetProductByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.Product, error) {
	result, err := repo.db(ctx).GetProductByVendorCode(ctx, sqlc.GetProductByVendorCodeParams{
		EditionID:  editionID,
		VendorCode: vendorCode,
	})
	return database.One(result, err, repo.dbe, mapProduct)
}

func (repo *Repo) GetVariantByID(ctx context.Context, id uuid.UUID) (*models.ProductVariant, error) {
	result, err := repo.db(ctx).GetProductVariantByID(ctx, id)
	return database.One(result, err, repo.dbe, mapVariant)
}

// GetVariantByIDForUpdate is the row-lock variant used inside the checkout
// tx (split 7): serializes concurrent checkouts on the same variant before
// availability is checked.
func (repo *Repo) GetVariantByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.ProductVariant, error) {
	result, err := repo.db(ctx).GetProductVariantByIDForUpdate(ctx, id)
	return database.One(result, err, repo.dbe, mapVariant)
}

func (repo *Repo) GetVariantByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.ProductVariant, error) {
	result, err := repo.db(ctx).GetProductVariantByVendorCode(ctx, sqlc.GetProductVariantByVendorCodeParams{
		EditionID:  editionID,
		VendorCode: vendorCode,
	})
	return database.One(result, err, repo.dbe, mapVariant)
}

func (repo *Repo) ListProductsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Product, error) {
	results, err := repo.db(ctx).ListProductsByEdition(ctx, editionID)
	return database.Many(results, err, repo.dbe, mapProduct)
}

func (repo *Repo) ListVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]models.ProductVariant, error) {
	results, err := repo.db(ctx).ListProductVariantsByProduct(ctx, productID)
	return database.Many(results, err, repo.dbe, mapVariant)
}

func (repo *Repo) PatchProduct(ctx context.Context, id uuid.UUID, product *models.Product) (*models.Product, error) {
	result, err := repo.db(ctx).PatchProduct(ctx, sqlc.PatchProductParams{
		VendorCode:           product.VendorCode,
		RequiresRegistration: product.RequiresRegistration,
		ID:                   id,
	})
	return database.One(result, err, repo.dbe, mapProduct)
}

func (repo *Repo) PatchVariant(ctx context.Context, id uuid.UUID, variant *models.ProductVariant) (*models.ProductVariant, error) {
	result, err := repo.db(ctx).PatchProductVariant(ctx, sqlc.PatchProductVariantParams{
		VendorCode:  variant.VendorCode,
		Name:        variant.Name,
		Description: variant.Description,
		Price:       variant.Price,
		Stock:       variant.Stock,
		GalleryUrls: variant.GalleryURLs,
		ID:          id,
	})
	return database.One(result, err, repo.dbe, mapVariant)
}

// CreateProductPurchase inserts a pending product purchase at checkout
// (split 7) — one of the purchase's materialized rows (D4).
func (repo *Repo) CreateProductPurchase(ctx context.Context, toCreate *models.ProductPurchase) (*models.ProductPurchase, error) {
	result, err := repo.db(ctx).CreateProductPurchase(ctx, sqlc.CreateProductPurchaseParams{
		EditionID:        toCreate.EditionID,
		VariantID:        toCreate.VariantID,
		PurchaserID:      toCreate.PurchaserID,
		RegistrationID:   toCreate.RegistrationID,
		Quantity:         toCreate.Quantity,
		Status:           string(toCreate.Status),
		StatusReason:     toCreate.StatusReason,
		PayssageIntentID: toCreate.PayssageIntentID,
	})
	return database.One(result, err, repo.dbe, mapProductPurchase)
}

// UpdateProductPurchaseStatus flips a product purchase (pending→confirmed/
// cancelled/expired) — the webhook receiver (split 4) and the expiry worker
// (split 7) write side.
func (repo *Repo) UpdateProductPurchaseStatus(ctx context.Context, id uuid.UUID, status models.ProductPurchaseStatus, reason *string) (*models.ProductPurchase, error) {
	result, err := repo.db(ctx).UpdateProductPurchaseStatus(ctx, sqlc.UpdateProductPurchaseStatusParams{
		ID:           id,
		Status:       string(status),
		StatusReason: reason,
	})
	return database.One(result, err, repo.dbe, mapProductPurchase)
}

// GetProductPurchaseByID reads one product purchase (expiry-worker tests and
// future read surfaces).
func (repo *Repo) GetProductPurchaseByID(ctx context.Context, id uuid.UUID) (*models.ProductPurchase, error) {
	result, err := repo.db(ctx).GetProductPurchaseByID(ctx, id)
	return database.One(result, err, repo.dbe, mapProductPurchase)
}
