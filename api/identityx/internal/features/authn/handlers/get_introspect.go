package handlers

import (
	"context"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) GetIntrospect(ctx context.Context, _ openapi.GetIntrospectRequestObject) (openapi.GetIntrospectResponseObject, error) {
	identity, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.GetIntrospect200JSONResponse(*identity), nil
}
