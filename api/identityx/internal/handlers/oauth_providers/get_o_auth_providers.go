package oauth_providers

import (
	"IdentityX/internal/openapi"
	"context"
	"time"
)

func (h *Handlers) GetOAuthProviders(ctx context.Context, req openapi.GetOAuthProvidersRequestObject) (openapi.GetOAuthProvidersResponseObject, error) {
	providers, err := h.ops.ListEnabledProviders(ctx, req.Params.ProjectId)
	if err != nil {
		return nil, err
	}
	items := make([]openapi.OAuthProviderDiscoveryItem, 0, len(providers))
	for _, p := range providers {
		items = append(items, openapi.OAuthProviderDiscoveryItem{Provider: p})
	}
	return openapi.GetOAuthProviders200JSONResponse{
		Code: 200, Data: &items, Timestamp: time.Now(), Module: module,
	}, nil
}
