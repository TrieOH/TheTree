package commands

import (
	"context"
	"encoding/json"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Create(ctx context.Context, in domain.CreateProductSpec) (out *domain.Product, err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	ga := uc.gaClient

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

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("products").
		Action("create").
		Scope(in.EditionScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("product").SetMessage("insufficient permissions")
	}

	span.SetAttributes(attribute.String("prodcut.id", validProduct.ID.String()))

	meta := json.RawMessage(`{"color": "#3bde09", "icon": "Gift", "folder": "products"}`)
	var scope *goauth.Scope
	var idStr = validProduct.ID.String()
	scope, err = ga.Scopes.CreateWithParent(ctx, validProduct.Name, &idStr, &in.EditionScopeID, meta)
	if err != nil {
		return nil, err
	}
	validProduct.AddScope(scope.ID)

	var created *domain.Product
	created, err = uc.products.Create(ctx, *validProduct)
	if err != nil {
		return nil, err
	}

	return created, nil
}
