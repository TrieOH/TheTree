package handlers

import (
	"context"
	"time"

	"IdentityX/errs"
	"IdentityX/internal/openapi"
)

func (h *Handlers) GetOAuthCallback(ctx context.Context, req openapi.GetOAuthCallbackRequestObject) (openapi.GetOAuthCallbackResponseObject, error) {
	if req.Params.Code == "" {
		return nil, errs.ErrMissingOAuthCode
	}
	if req.Params.State == "" {
		return nil, errs.ErrMissingOAuthState
	}
	tokens, err := h.ops.Callback(ctx, string(req.Provider), req.Params.Code, req.Params.State)
	if err != nil {
		return nil, err
	}
	return openapi.GetOAuthCallback201JSONResponse{
		Code: 201, Data: tokens, Timestamp: time.Now(), Module: module,
	}, nil
}
