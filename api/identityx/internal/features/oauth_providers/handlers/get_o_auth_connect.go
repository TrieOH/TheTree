package handlers

import (
	"context"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetOAuthConnect(ctx context.Context, req openapi.GetOAuthConnectRequestObject) (openapi.GetOAuthConnectResponseObject, error) {
	acceptedTos := false
	if req.Params.AcceptedTos != nil {
		acceptedTos = *req.Params.AcceptedTos
	}
	url, err := h.ops.Connect(ctx, string(req.Provider), req.Params.ProjectId, acceptedTos)
	if err != nil {
		return nil, err
	}
	return openapi.GetOAuthConnect200JSONResponse{Url: url}, nil
}
