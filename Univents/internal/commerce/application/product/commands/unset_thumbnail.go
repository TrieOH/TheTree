package commands

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) UnsetThumbnail(ctx context.Context, id uuid.UUID) (product *domain.Product, err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.UnsetThumbnail")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("unset_thumbnail.success", err == nil)) }()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	product, err = uc.products.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.gaClient.Authz.Check().User(sub.ID).
		Object("products").
		Action("edit").
		Scope(product.ScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("product").SetMessage("insufficient permissions")
	}

	product, err = uc.products.UnsetThumbnail(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	return product, nil
}
