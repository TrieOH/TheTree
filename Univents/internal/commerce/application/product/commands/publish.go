package commands

import (
	"context"
	"errors"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Publish(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.Publish")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("publish.success", err == nil))
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
		Action("publish").
		Scope(product.ScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return errx.Forbidden("product").SetMessage("insufficient permissions")
	}

	if product.Status != domain.ProductStatusDraft {
		return errors.New("can't publish products on statuses different than draft")
	}

	// TODO: ADD ASYNQ TASKS FOR PRODUCT AVAILABILITY?

	if err = uc.products.Publish(ctx, product.ID); err != nil {
		return err
	}

	return nil
}
