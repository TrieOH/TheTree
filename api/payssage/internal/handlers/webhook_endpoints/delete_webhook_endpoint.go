package webhook_endpoints

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) DeleteWebhookEndpoint(ctx context.Context, req openapi.DeleteWebhookEndpointRequestObject) (openapi.DeleteWebhookEndpointResponseObject, error) {
	err := h.ops.Delete(ctx, req.EndpointId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteWebhookEndpoint204Response{}, nil
}
