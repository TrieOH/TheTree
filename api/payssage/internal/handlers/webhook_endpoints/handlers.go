// Package webhook_endpoints implements the StrictServerInterface methods
// for the webhook_endpoints feature.
package webhook_endpoints

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/internal/services"
	"payssage/models"
)

const module = "Payssage"

type Handlers struct {
	ops *services.WebhookEndpoints
}

func New(ops *services.WebhookEndpoints) *Handlers { return &Handlers{ops: ops} }

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

func (h *Handlers) ListWebhookEndpoints(ctx context.Context, req openapi.ListWebhookEndpointsRequestObject) (openapi.ListWebhookEndpointsResponseObject, error) {
	endpoints, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWebhookEndpoints200JSONResponse{
		Code: 200, Data: &endpoints, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetWebhookEndpoint(ctx context.Context, req openapi.GetWebhookEndpointRequestObject) (openapi.GetWebhookEndpointResponseObject, error) {
	endpoint, err := h.ops.GetByID(ctx, req.EndpointId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWebhookEndpoint200JSONResponse{
		Code: 200, Data: endpoint, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) DeleteWebhookEndpoint(ctx context.Context, req openapi.DeleteWebhookEndpointRequestObject) (openapi.DeleteWebhookEndpointResponseObject, error) {
	err := h.ops.Delete(ctx, req.EndpointId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteWebhookEndpoint204Response{}, nil
}
