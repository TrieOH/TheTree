package oauth_providers

import (
	"IdentityX/internal/openapi"
	"context"
	"time"
)

func (h *Handlers) DeleteOAuthProvider(ctx context.Context, req openapi.DeleteOAuthProviderRequestObject) (openapi.DeleteOAuthProviderResponseObject, error) {
	err := h.ops.Delete(ctx, req.OauthProviderId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteOAuthProvider200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
