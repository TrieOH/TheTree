package handlers

import (
	"context"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) AddEventMember(ctx context.Context, req openapi.AddEventMemberRequestObject) (openapi.AddEventMemberResponseObject, error) {
	member, err := h.ops.AddMember(ctx, req.EventId, models.AddEventMemberInput{
		Email: req.Body.Email,
		Role:  req.Body.Role,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddEventMember201JSONResponse(*member), nil
}
