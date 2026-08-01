package events

import (
	"context"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) RemoveEventMember(ctx context.Context, req openapi.RemoveEventMemberRequestObject) (openapi.RemoveEventMemberResponseObject, error) {
	err := h.ops.RemoveMember(ctx, req.EventId, models.RemoveMemberInput{
		Email: req.Body.Email,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveEventMember204Response{}, nil
}
