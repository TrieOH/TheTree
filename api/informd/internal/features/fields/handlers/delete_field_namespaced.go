package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) DeleteFieldNamespaced(ctx context.Context, req openapi.DeleteFieldNamespacedRequestObject) (openapi.DeleteFieldNamespacedResponseObject, error) {
	err := h.ops.DeleteNamespaced(ctx, req.NamespaceId, req.FormId, req.FieldId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteFieldNamespaced204Response{}, nil
}
