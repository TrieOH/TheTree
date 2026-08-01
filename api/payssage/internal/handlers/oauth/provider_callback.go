package oauth

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) ProviderCallback(ctx context.Context, req openapi.ProviderCallbackRequestObject) (openapi.ProviderCallbackResponseObject, error) {
	finalRedirectURI, err := h.ops.Callback(ctx, req.Provider, req.Params.Code, req.Params.State)
	if err != nil {
		return nil, err
	}
	return openapi.ProviderCallback200JSONResponse{
		Code: 200, Data: &finalRedirectURI, Timestamp: time.Now(), Module: module,
	}, nil
}
