package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) ListFormMembersNamespaced(ctx context.Context, req openapi.ListFormMembersNamespacedRequestObject) (openapi.ListFormMembersNamespacedResponseObject, error) {
	members, err := h.ops.ListFormMembers(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ListFormMembersNamespaced200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}
