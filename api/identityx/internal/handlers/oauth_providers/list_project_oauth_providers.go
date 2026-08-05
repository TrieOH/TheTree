package oauth_providers

import (
	"IdentityX/internal/openapi"
	"context"
	"time"
)

func (h *Handlers) ListProjectOAuthProviders(ctx context.Context, req openapi.ListProjectOAuthProvidersRequestObject) (openapi.ListProjectOAuthProvidersResponseObject, error) {
	providers, err := h.ops.ListByProject(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListProjectOAuthProviders200JSONResponse{
		Code: 200, Data: &providers, Timestamp: time.Now(), Module: module,
	}, nil
}
