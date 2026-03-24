package commands

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Restore(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.Restore")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("restore.success", err == nil))
	}()

	ga := uc.gaClient

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

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("products").
		Action("restore").
		Scope(product.ScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return errx.Forbidden("product").SetMessage("insufficient permissions")
	}

	if err = uc.products.Restore(ctx, product.ID); err != nil {
		return err
	}

	return nil
}
