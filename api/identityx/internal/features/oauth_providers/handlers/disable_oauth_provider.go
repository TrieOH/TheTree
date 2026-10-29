package handlers

import (
	"IdentityX/internal/openapi"
	"context"
)

func (h *Handlers) DisableOAuthProvider(ctx context.Context, req openapi.DisableOAuthProviderRequestObject) (openapi.DisableOAuthProviderResponseObject, error) {
	provider, err := h.ops.SetEnabled(ctx, req.OauthProviderId, false)
	if err != nil {
		return nil, err
	}
	return openapi.DisableOAuthProvider200JSONResponse(*provider), nil
}
