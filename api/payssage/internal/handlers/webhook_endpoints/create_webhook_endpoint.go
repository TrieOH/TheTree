package webhook_endpoints

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/models"
)

func (h *Handlers) CreateWebhookEndpoint(ctx context.Context, req openapi.CreateWebhookEndpointRequestObject) (openapi.CreateWebhookEndpointResponseObject, error) {
	endpoint, err := h.ops.Create(ctx, models.CreateWebhookEndpointInput{
		WalletID: req.WalletId,
		Name:     req.Body.Name,
		URL:      req.Body.Url,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateWebhookEndpoint201JSONResponse{
		Code: 201, Data: endpoint, Timestamp: time.Now(), Module: module,
	}, nil
}
