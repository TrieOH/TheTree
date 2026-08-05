package certifications

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionCertifications(ctx context.Context, req openapi.ListEditionCertificationsRequestObject) (openapi.ListEditionCertificationsResponseObject, error) {
	certs, err := h.ops.ListCertsByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionCertifications200JSONResponse{
		Code: 200, Data: &certs, Timestamp: time.Now(), Module: module,
	}, nil
}
