package oauth_providers

import (
	"IdentityX/internal/openapi"
	"context"
	"time"
)

func (h *Handlers) DisableOAuthProvider(ctx context.Context, req openapi.DisableOAuthProviderRequestObject) (openapi.DisableOAuthProviderResponseObject, error) {
	provider, err := h.ops.SetEnabled(ctx, req.OauthProviderId, false)
	if err != nil {
		return nil, err
	}
	return openapi.DisableOAuthProvider200JSONResponse{
		Code: 200, Data: provider, Timestamp: time.Now(), Module: module,
	}, nil
}
