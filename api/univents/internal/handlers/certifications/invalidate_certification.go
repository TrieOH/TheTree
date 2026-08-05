package certifications

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) InvalidateCertification(ctx context.Context, req openapi.InvalidateCertificationRequestObject) (openapi.InvalidateCertificationResponseObject, error) {
	err := h.ops.InvalidateCert(ctx, req.CertId, &req.Body.Reason)
	if err != nil {
		return nil, err
	}
	return openapi.InvalidateCertification204Response{}, nil
}
