package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) ProviderCallback(ctx context.Context, req openapi.ProviderCallbackRequestObject) (openapi.ProviderCallbackResponseObject, error) {
	finalRedirectURI, err := h.ops.Callback(ctx, req.Provider, req.Params.Code, req.Params.State)
	if err != nil {
		return nil, err
	}
	return openapi.ProviderCallback200JSONResponse(finalRedirectURI), nil
}
