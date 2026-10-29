package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) GetSignature(ctx context.Context, req openapi.GetSignatureRequestObject) (openapi.GetSignatureResponseObject, error) {
	signature, err := h.ops.GetByID(ctx, req.SignatureId)
	if err != nil {
		return nil, err
	}
	return openapi.GetSignature200JSONResponse(*signature), nil
}
