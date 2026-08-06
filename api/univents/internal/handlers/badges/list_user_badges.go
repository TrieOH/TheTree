package badges

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListUserBadges(ctx context.Context, req openapi.ListUserBadgesRequestObject) (openapi.ListUserBadgesResponseObject, error) {
	groups, err := h.ops.ListByUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return openapi.ListUserBadges200JSONResponse{
		Code: 200, Data: groups, Timestamp: time.Now(), Module: module,
	}, nil
}
