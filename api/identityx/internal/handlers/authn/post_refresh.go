package authn

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) PostRefresh(ctx context.Context, req openapi.PostRefreshRequestObject) (openapi.PostRefreshResponseObject, error) {
	tokens, err := h.ops.Refresh(ctx, req.Params.RefreshToken)
	if err != nil {
		return nil, err
	}
	return openapi.PostRefresh200JSONResponse{
		Code: 200, Data: tokens, Timestamp: time.Now(), Module: module,
	}, nil
}
