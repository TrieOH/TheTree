package signatures

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionSignatures(ctx context.Context, req openapi.ListEditionSignaturesRequestObject) (openapi.ListEditionSignaturesResponseObject, error) {
	signatures, err := h.ops.ListByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionSignatures200JSONResponse{
		Code: 200, Data: &signatures, Timestamp: time.Now(), Module: module,
	}, nil
}
