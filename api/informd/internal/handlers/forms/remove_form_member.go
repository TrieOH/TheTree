package forms

import (
	"context"
	"time"

	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) RemoveFormMember(ctx context.Context, req openapi.RemoveFormMemberRequestObject) (openapi.RemoveFormMemberResponseObject, error) {
	err := h.ops.RemoveMember(ctx, models.RemoveFormMemberInput{
		UserID: req.Body.UserId,
		FormID: req.FormId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveFormMember200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
