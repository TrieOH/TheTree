package signatures

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) CancelSignatureRequest(ctx context.Context, req openapi.CancelSignatureRequestRequestObject) (openapi.CancelSignatureRequestResponseObject, error) {
	err := h.ops.CancelRequest(ctx, req.RequestId, req.Body.Reason)
	if err != nil {
		return nil, err
	}
	return openapi.CancelSignatureRequest204Response{}, nil
}
