package oauth

import (
	"context"

	"payssage/internal/features/oauth/service"
	"payssage/internal/openapi"
	"payssage/models"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ConnectProvider(ctx context.Context, req openapi.ConnectProviderRequestObject) (openapi.ConnectProviderResponseObject, error) {
	url, err := h.ops.Connect(ctx, models.ConnectInput{
		Provider:         req.Provider,
		Flow:             req.Body.Flow,
		FinalRedirectURL: req.Body.FinalRedirectUrl,
		WalletID:         req.Body.WalletId,
		OrganizationID:   req.Body.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.ConnectProvider200JSONResponse(url), nil
}

func (h *Handlers) ProviderCallback(ctx context.Context, req openapi.ProviderCallbackRequestObject) (openapi.ProviderCallbackResponseObject, error) {
	finalRedirectURI, err := h.ops.Callback(ctx, req.Provider, req.Params.Code, req.Params.State)
	if err != nil {
		return nil, err
	}
	return openapi.ProviderCallback200JSONResponse(finalRedirectURI), nil
}

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
