package commands

import (
	"context"
	"univents/internal/commerce/domain"
	domain2 "univents/internal/core/domain"
	"univents/internal/shared/authz"

	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Create(ctx context.Context, in domain.CreateProductSpec) (out *domain.Product, err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var validProduct *domain.Product
	validProduct, err = domain.NewProduct(sub.ID, in)
	if err != nil {
		return nil, err
	}

	var edition *domain2.Edition
	edition, err = uc.editions.GetByID(ctx, in.EditionID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_products"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return nil, err
	}

	var created *domain.Product
	created, err = uc.products.Create(ctx, *validProduct)
	if err != nil {
		return nil, err
	}

	return created, nil
}
