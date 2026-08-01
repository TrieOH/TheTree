// Package oauth implements the StrictServerInterface methods for the oauth
// feature. The provider callback is public (browser redirect target) and
// returns the final redirect URL as the envelope data string.
package oauth

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/internal/services"
	"payssage/models"
)

const module = "Payssage"

type Handlers struct {
	ops *services.OAuth
}

func New(ops *services.OAuth) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ConnectProvider(ctx context.Context, req openapi.ConnectProviderRequestObject) (openapi.ConnectProviderResponseObject, error) {
	url, err := h.ops.Connect(ctx, models.ConnectInput{
		Provider:            req.Provider,
		Flow:                req.Body.Flow,
		ProviderRedirectURL: req.Body.ProviderRedirectUrl,
		FinalRedirectURL:    req.Body.FinalRedirectUrl,
		WalletID:            req.Body.WalletId,
		OrganizationID:      req.Body.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.ConnectProvider200JSONResponse{
		Code: 200, Data: &url, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ProviderCallback(ctx context.Context, req openapi.ProviderCallbackRequestObject) (openapi.ProviderCallbackResponseObject, error) {
	finalRedirectURI, err := h.ops.Callback(ctx, req.Provider, req.Params.Code, req.Params.State)
	if err != nil {
		return nil, err
	}
	return openapi.ProviderCallback200JSONResponse{
		Code: 200, Data: &finalRedirectURI, Timestamp: time.Now(), Module: module,
	}, nil
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
