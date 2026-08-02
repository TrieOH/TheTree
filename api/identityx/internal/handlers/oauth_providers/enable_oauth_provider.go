package oauth_providers

import (
	"IdentityX/internal/openapi"
	"context"
	"time"
)

func (h *Handlers) EnableOAuthProvider(ctx context.Context, req openapi.EnableOAuthProviderRequestObject) (openapi.EnableOAuthProviderResponseObject, error) {
	provider, err := h.ops.SetEnabled(ctx, req.OauthProviderId, true)
	if err != nil {
		return nil, err
	}
	return openapi.EnableOAuthProvider200JSONResponse{
		Code: 200, Data: provider, Timestamp: time.Now(), Module: module,
	}, nil
}
