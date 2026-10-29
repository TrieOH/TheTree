package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) GetSignatureRequest(ctx context.Context, req openapi.GetSignatureRequestRequestObject) (openapi.GetSignatureRequestResponseObject, error) {
	request, err := h.ops.GetRequestByID(ctx, req.RequestId)
	if err != nil {
		return nil, err
	}
	return openapi.GetSignatureRequest200JSONResponse(*request), nil
}
