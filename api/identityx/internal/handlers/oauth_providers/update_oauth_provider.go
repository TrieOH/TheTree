package oauth_providers

import (
	"IdentityX/internal/openapi"
	"IdentityX/models"
	"context"
	"time"
)

func (h *Handlers) UpdateOAuthProvider(ctx context.Context, req openapi.UpdateOAuthProviderRequestObject) (openapi.UpdateOAuthProviderResponseObject, error) {
	updated, err := h.ops.Update(ctx, models.UpdateOAuthProviderInput{
		ID:           req.OauthProviderId,
		ClientID:     req.Body.ClientId,
		ClientSecret: req.Body.ClientSecret,
	})
	if err != nil {
		return nil, err
	}
	return openapi.UpdateOAuthProvider200JSONResponse{
		Code: 200, Data: updated, Timestamp: time.Now(), Module: module,
	}, nil
}
