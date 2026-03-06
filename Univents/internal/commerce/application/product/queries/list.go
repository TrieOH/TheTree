package queries

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
)

func (uc *QueryService) List(ctx context.Context, editionID uuid.UUID) (out []domain.Product, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "ProductService.List")
	defer span.End()

	return uc.products.List(ctx, editionID)
}

func (uc *QueryService) AdminList(ctx context.Context, editionID uuid.UUID) (out []domain.Product, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "ProductService.List")
	defer span.End()

	edition, err := uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	allowed, err := ga.Authz.Check().
		User(sub.ID).
		Object("products").
		Action("read").
		Scope(edition.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("products").SetMessage("insufficient permissions")
	}

	return uc.products.AdminList(ctx, editionID)
}
