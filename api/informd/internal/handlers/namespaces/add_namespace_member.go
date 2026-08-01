package namespaces

import (
	"context"

	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) AddNamespaceMember(ctx context.Context, req openapi.AddNamespaceMemberRequestObject) (openapi.AddNamespaceMemberResponseObject, error) {
	err := h.ops.AddMember(ctx, models.AddNamespaceMemberInput{
		UserID:      req.Body.UserId,
		Role:        req.Body.Role,
		NamespaceID: req.NamespaceId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddNamespaceMember201Response{}, nil
}
