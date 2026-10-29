package handlers

import (
	"IdentityX/internal/openapi"
	"context"
)

func (h *Handlers) EnableOAuthProvider(ctx context.Context, req openapi.EnableOAuthProviderRequestObject) (openapi.EnableOAuthProviderResponseObject, error) {
	provider, err := h.ops.SetEnabled(ctx, req.OauthProviderId, true)
	if err != nil {
		return nil, err
	}
	return openapi.EnableOAuthProvider200JSONResponse(*provider), nil
}
