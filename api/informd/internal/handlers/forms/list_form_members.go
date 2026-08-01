package forms

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) ListFormMembers(ctx context.Context, req openapi.ListFormMembersRequestObject) (openapi.ListFormMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ListFormMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}
