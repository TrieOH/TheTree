package service

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) PatchVariant(ctx context.Context, payload models.PatchProductVariantInput) (*models.ProductVariant, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.PatchVariant")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := o.Products.GetVariantByID(ctx, payload.VariantID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
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
