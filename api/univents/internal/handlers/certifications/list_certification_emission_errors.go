package certifications

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListCertificationEmissionErrors(ctx context.Context, req openapi.ListCertificationEmissionErrorsRequestObject) (openapi.ListCertificationEmissionErrorsResponseObject, error) {
	errors, err := h.ops.ListEmissionErrors(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListCertificationEmissionErrors200JSONResponse{
		Code: 200, Data: &errors, Timestamp: time.Now(), Module: module,
	}, nil
}
