package handlers

import (
	"context"

	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) RemoveNamespaceMember(ctx context.Context, req openapi.RemoveNamespaceMemberRequestObject) (openapi.RemoveNamespaceMemberResponseObject, error) {
	err := h.ops.RemoveMember(ctx, models.RemoveNamespaceMemberInput{
		UserID:      req.Body.UserId,
		NamespaceID: req.NamespaceId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveNamespaceMember204Response{}, nil
}
