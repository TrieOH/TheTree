package signatures

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) GetSignatureRequest(ctx context.Context, req openapi.GetSignatureRequestRequestObject) (openapi.GetSignatureRequestResponseObject, error) {
	request, err := h.ops.GetRequestByID(ctx, req.RequestId)
	if err != nil {
		return nil, err
	}
	return openapi.GetSignatureRequest200JSONResponse{
		Code: 200, Data: request, Timestamp: time.Now(), Module: module,
	}, nil
}
