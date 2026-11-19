// Package products is the products feature: its repo, handlers and
// operations. Build wires it.
package products

import (
	"context"
	"lib/database"
	"univents/internal/authz"
	"univents/models"
	"univents/ports"
	"uuid"
)

type Deps struct {
	Events   ports.EventRepo
	Editions ports.EditionRepo
	Products ports.ProductRepo
	Authz    *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) CreateInitial(ctx context.Context, payload models.CreateInitialProductInput) (*models.Product, error) {
	edition, err := o.Editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	var product *models.Product
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		p, err := o.Products.CreateProduct(ctx, &models.Product{
			EditionID:            payload.EditionID,
			VendorCode:           payload.VendorCode,
			RequiresRegistration: payload.RequiresRegistration,
		})
		if err != nil {
			return err
		}

		_, err = o.Products.CreateVariant(ctx, &models.ProductVariant{
			EditionID:   payload.EditionID,
			ProductID:   p.ID,
			VendorCode:  payload.VariantVendorCode,
			Name:        payload.Name,
			Description: payload.Description,
			Price:       payload.Price,
			Stock:       payload.Stock,
		})
		if err != nil {
			return err
		}

		product = p
		return nil
	})
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (o *Operations) CreateVariant(ctx context.Context, payload models.CreateProductVariantInput) (*models.ProductVariant, error) {
	product, err := o.Products.GetProductByID(ctx, payload.ProductID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, product.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	variant := &models.ProductVariant{
		EditionID:   product.EditionID,
		ProductID:   product.ID,
		VendorCode:  payload.VendorCode,
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Stock:       payload.Stock,
	}

	return o.Products.CreateVariant(ctx, variant)
}

// DeleteProduct deletes the product and returns it as it was.
func (o *Operations) DeleteProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	existing, err := o.Products.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	err = o.Products.DeleteProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	return existing, nil
}

func (o *Operations) DeleteVariant(ctx context.Context, id uuid.UUID) error {
	existing, err := o.Products.GetVariantByID(ctx, id)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Products.DeleteVariant(ctx, id)
}

func (o *Operations) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	return o.Products.GetProductByID(ctx, id)
}

func (o *Operations) GetProductByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.Product, error) {
	return o.Products.GetProductByVendorCode(ctx, editionID, vendorCode)
}

func (o *Operations) GetVariantByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.ProductVariant, error) {
	return o.Products.GetVariantByVendorCode(ctx, editionID, vendorCode)
}

func (o *Operations) ListProductsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Product, error) {
	return o.Products.ListProductsByEdition(ctx, editionID)
}

func (o *Operations) ListVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]models.ProductVariant, error) {
	return o.Products.ListVariantsByProduct(ctx, productID)
}

func (o *Operations) PatchProduct(ctx context.Context, payload models.PatchProductInput) (*models.Product, error) {
	existing, err := o.Products.GetProductByID(ctx, payload.ProductID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		VendorCode:           payload.VendorCode,
		RequiresRegistration: payload.RequiresRegistration,
	}

	return o.Products.PatchProduct(ctx, payload.ProductID, product)
}

func (o *Operations) PatchVariant(ctx context.Context, payload models.PatchProductVariantInput) (*models.ProductVariant, error) {
	existing, err := o.Products.GetVariantByID(ctx, payload.VariantID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	variant := &models.ProductVariant{
		VendorCode:  payload.VendorCode,
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Stock:       payload.Stock,
		GalleryURLs: payload.GalleryURLs,
	}

	return o.Products.PatchVariant(ctx, existing.ID, variant)
}
