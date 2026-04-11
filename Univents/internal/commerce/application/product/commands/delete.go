package commands

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Delete(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.Delete")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("delete.success", err == nil))
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
		authz.Permission("delete"),
		authz.Resource("product", product.ID.String()),
	); err != nil {
		return err
	}

	var hasPurchase bool
	hasPurchase, err = uc.products.ItemHasCompletedPurchases(ctx, product.ID)
	if err != nil {
		return err
	}

	if hasPurchase {
		return errx.Forbidden("product").SetMessage("cannot delete a product that was already purchased")
	}

	if err = uc.products.Delete(ctx, product.ID); err != nil {
		return err
	}

	return nil
}
