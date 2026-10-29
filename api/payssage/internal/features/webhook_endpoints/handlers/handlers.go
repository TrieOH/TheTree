package handlers

import (
	"context"

	"payssage/internal/features/webhook_endpoints/service"
	"payssage/internal/openapi"
	"payssage/models"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) CreateWebhookEndpoint(ctx context.Context, req openapi.CreateWebhookEndpointRequestObject) (openapi.CreateWebhookEndpointResponseObject, error) {
	endpoint, err := h.ops.Create(ctx, models.CreateWebhookEndpointInput{
		WalletID: req.WalletId,
		Name:     req.Body.Name,
		URL:      req.Body.Url,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateWebhookEndpoint201JSONResponse(*endpoint), nil
}

func (h *Handlers) DeleteWebhookEndpoint(ctx context.Context, req openapi.DeleteWebhookEndpointRequestObject) (openapi.DeleteWebhookEndpointResponseObject, error) {
	err := h.ops.Delete(ctx, req.EndpointId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteWebhookEndpoint204Response{}, nil
}

func (h *Handlers) GetWebhookEndpoint(ctx context.Context, req openapi.GetWebhookEndpointRequestObject) (openapi.GetWebhookEndpointResponseObject, error) {
	endpoint, err := h.ops.GetByID(ctx, req.EndpointId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWebhookEndpoint200JSONResponse(*endpoint), nil
}

func (h *Handlers) ListWebhookEndpoints(ctx context.Context, req openapi.ListWebhookEndpointsRequestObject) (openapi.ListWebhookEndpointsResponseObject, error) {
	endpoints, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWebhookEndpoints200JSONResponse{Items: endpoints}, nil
}
