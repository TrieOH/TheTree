package namespaces

import (
	"context"

	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) AddFormMemberNamespaced(ctx context.Context, req openapi.AddFormMemberNamespacedRequestObject) (openapi.AddFormMemberNamespacedResponseObject, error) {
	err := h.ops.AddFormMember(ctx, models.AddNamespaceFormMemberInput{
		UserID:      req.Body.UserId,
		NamespaceID: req.NamespaceId,
		FormID:      req.FormId,
		Role:        req.Body.Role,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddFormMemberNamespaced201Response{}, nil
}
