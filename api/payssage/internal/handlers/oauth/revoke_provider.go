package oauth

import (
	"context"

	"payssage/internal/openapi"
	"payssage/models"
)

func (h *Handlers) RevokeProvider(ctx context.Context, req openapi.RevokeProviderRequestObject) (openapi.RevokeProviderResponseObject, error) {
	err := h.ops.Revoke(ctx, models.RevokeInput{
		Flow:     req.Body.Flow,
		ID:       req.Body.Id,
		Provider: req.Provider,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RevokeProvider204Response{}, nil
}
