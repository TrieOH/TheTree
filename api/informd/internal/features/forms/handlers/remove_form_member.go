package handlers

import (
	"context"

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
	return openapi.RemoveFormMember204Response{}, nil
}
