package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) GetCertification(ctx context.Context, req openapi.GetCertificationRequestObject) (openapi.GetCertificationResponseObject, error) {
	cert, err := h.ops.GetCertByID(ctx, req.CertId)
	if err != nil {
		return nil, err
	}
	return openapi.GetCertification200JSONResponse(*cert), nil
}
