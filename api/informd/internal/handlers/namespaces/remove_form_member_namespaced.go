package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) RemoveFormMemberNamespaced(ctx context.Context, req openapi.RemoveFormMemberNamespacedRequestObject) (openapi.RemoveFormMemberNamespacedResponseObject, error) {
	err := h.ops.RemoveFormMember(ctx, models.RemoveNamespaceFormMemberInput{
		UserID:      req.Body.UserId,
		NamespaceID: req.NamespaceId,
		FormID:      req.FormId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveFormMemberNamespaced200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
