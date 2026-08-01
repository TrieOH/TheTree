package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) ListNamespaceMembers(ctx context.Context, req openapi.ListNamespaceMembersRequestObject) (openapi.ListNamespaceMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.ListNamespaceMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}
