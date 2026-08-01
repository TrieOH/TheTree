package signatures

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) RevokeSignature(ctx context.Context, req openapi.RevokeSignatureRequestObject) (openapi.RevokeSignatureResponseObject, error) {
	err := h.ops.RevokeSignature(ctx, req.Params.Token)
	if err != nil {
		return nil, err
	}
	return openapi.RevokeSignature204Response{}, nil
}
