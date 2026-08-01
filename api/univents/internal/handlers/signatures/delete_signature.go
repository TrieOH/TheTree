package signatures

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) DeleteSignature(ctx context.Context, req openapi.DeleteSignatureRequestObject) (openapi.DeleteSignatureResponseObject, error) {
	err := h.ops.Delete(ctx, req.SignatureId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteSignature204Response{}, nil
}
