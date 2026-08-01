package webhook_endpoints

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) GetWebhookEndpoint(ctx context.Context, req openapi.GetWebhookEndpointRequestObject) (openapi.GetWebhookEndpointResponseObject, error) {
	endpoint, err := h.ops.GetByID(ctx, req.EndpointId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWebhookEndpoint200JSONResponse{
		Code: 200, Data: endpoint, Timestamp: time.Now(), Module: module,
	}, nil
}
