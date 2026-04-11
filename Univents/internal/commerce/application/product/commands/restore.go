package commands

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Restore(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.Restore")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("restore.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var product *domain.Product
	product, err = uc.products.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("restore"),
		authz.Resource("product", product.ID.String()),
	); err != nil {
		return err
	}

	if err = uc.products.Restore(ctx, product.ID); err != nil {
		return err
	}

	return nil
}
