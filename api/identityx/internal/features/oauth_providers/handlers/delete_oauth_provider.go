package handlers

import (
	"IdentityX/internal/openapi"
	"context"
)

func (h *Handlers) DeleteOAuthProvider(ctx context.Context, req openapi.DeleteOAuthProviderRequestObject) (openapi.DeleteOAuthProviderResponseObject, error) {
	err := h.ops.Delete(ctx, req.OauthProviderId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteOAuthProvider204Response{}, nil
}
