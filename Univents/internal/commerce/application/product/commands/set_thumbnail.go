package commands

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) SetThumbnail(ctx context.Context, id uuid.UUID, url string) (product *domain.Product, err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.SetThumbnail")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("set_thumbnail.success", err == nil)) }()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	product, err = uc.products.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("product", product.ID.String()),
	); err != nil {
		return nil, err
	}

	product, err = uc.products.SetThumbnail(ctx, product.ID, url)
	if err != nil {
		return nil, err
	}

	return product, nil
}
