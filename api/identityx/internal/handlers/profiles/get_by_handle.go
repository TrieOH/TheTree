package profiles

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetProfileByHandle(ctx context.Context, req openapi.GetProfileByHandleRequestObject) (openapi.GetProfileByHandleResponseObject, error) {
	profile, err := h.ops.GetProfileByHandle(ctx, req.Handle)
	if err != nil {
		return nil, err
	}
	return openapi.GetProfileByHandle200JSONResponse{
		Code: 200, Data: profile, Timestamp: time.Now(), Module: module,
	}, nil
}
