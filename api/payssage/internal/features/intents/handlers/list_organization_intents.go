package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) ListOrganizationIntents(ctx context.Context, req openapi.ListOrganizationIntentsRequestObject) (openapi.ListOrganizationIntentsResponseObject, error) {
	intents, err := h.ops.ListByOrg(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationIntents200JSONResponse{Items: intents}, nil
}
