package queries

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"

	"github.com/google/uuid"
)

func (uc *QueryService) List(ctx context.Context, editionID uuid.UUID) (out []domain.Product, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "ProductService.List")
	defer span.End()

	return uc.products.List(ctx, editionID)
}

func (uc *QueryService) AdminList(ctx context.Context, editionID uuid.UUID) (out []domain.Product, err error) { // FIXME Pagination
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
