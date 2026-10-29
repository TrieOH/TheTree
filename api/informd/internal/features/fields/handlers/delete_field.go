package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) DeleteField(ctx context.Context, req openapi.DeleteFieldRequestObject) (openapi.DeleteFieldResponseObject, error) {
	err := h.ops.Delete(ctx, req.FormId, req.FieldId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteField204Response{}, nil
}
