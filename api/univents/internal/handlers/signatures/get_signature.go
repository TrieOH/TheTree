package signatures

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) GetSignature(ctx context.Context, req openapi.GetSignatureRequestObject) (openapi.GetSignatureResponseObject, error) {
	signature, err := h.ops.GetByID(ctx, req.SignatureId)
	if err != nil {
		return nil, err
	}
	return openapi.GetSignature200JSONResponse{
		Code: 200, Data: signature, Timestamp: time.Now(), Module: module,
	}, nil
}
