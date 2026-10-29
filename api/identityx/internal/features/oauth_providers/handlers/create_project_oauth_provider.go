package handlers

import (
	"IdentityX/internal/openapi"
	"IdentityX/models"
	"context"
)

func (h *Handlers) CreateProjectOAuthProvider(ctx context.Context, req openapi.CreateProjectOAuthProviderRequestObject) (openapi.CreateProjectOAuthProviderResponseObject, error) {
	created, err := h.ops.Create(ctx, models.CreateOAuthProviderInput{
		ProjectID:    req.ProjectId,
		Provider:     req.Body.Provider,
		ClientID:     req.Body.ClientId,
		ClientSecret: req.Body.ClientSecret,
		CallbackURL:  req.Body.CallbackUrl,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateProjectOAuthProvider201JSONResponse(*created), nil
}
