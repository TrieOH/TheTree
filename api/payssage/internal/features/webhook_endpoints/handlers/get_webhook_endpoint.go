package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) GetWebhookEndpoint(ctx context.Context, req openapi.GetWebhookEndpointRequestObject) (openapi.GetWebhookEndpointResponseObject, error) {
	endpoint, err := h.ops.GetByID(ctx, req.EndpointId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWebhookEndpoint200JSONResponse(*endpoint), nil
}
