package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionCertifications(ctx context.Context, req openapi.ListEditionCertificationsRequestObject) (openapi.ListEditionCertificationsResponseObject, error) {
	certs, err := h.ops.ListCertsByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionCertifications200JSONResponse{Items: certs}, nil
}
