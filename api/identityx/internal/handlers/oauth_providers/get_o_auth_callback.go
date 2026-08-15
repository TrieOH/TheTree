package oauth_providers

import (
	"context"
	"time"

	"github.com/MintzyG/fun"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetOAuthCallback(ctx context.Context, req openapi.GetOAuthCallbackRequestObject) (openapi.GetOAuthCallbackResponseObject, error) {
	if req.Params.Code == "" {
		return nil, fun.ErrBadRequest("missing code")
	}
	if req.Params.State == "" {
		return nil, fun.ErrBadRequest("missing state")
	}
	tokens, err := h.ops.Callback(ctx, string(req.Provider), req.Params.Code, req.Params.State)
	if err != nil {
		return nil, err
	}
	return openapi.GetOAuthCallback201JSONResponse{
		Code: 201, Data: tokens, Timestamp: time.Now(), Module: module,
	}, nil
}
