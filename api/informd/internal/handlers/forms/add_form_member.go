package forms

import (
	"context"

	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) AddFormMember(ctx context.Context, req openapi.AddFormMemberRequestObject) (openapi.AddFormMemberResponseObject, error) {
	err := h.ops.AddMember(ctx, models.AddFormMemberInput{
		UserID: req.Body.UserId,
		FormID: req.FormId,
		Role:   req.Body.Role,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddFormMember201Response{}, nil
}
