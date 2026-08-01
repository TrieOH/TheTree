package signatures

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) DenySignatureRequest(ctx context.Context, req openapi.DenySignatureRequestRequestObject) (openapi.DenySignatureRequestResponseObject, error) {
	err := h.ops.DenyRequest(ctx, req.Params.Token, req.Body.Reason)
	if err != nil {
		return nil, err
	}
	return openapi.DenySignatureRequest204Response{}, nil
}
