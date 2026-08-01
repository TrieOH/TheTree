package certifications

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) VerifyCertification(ctx context.Context, req openapi.VerifyCertificationRequestObject) (openapi.VerifyCertificationResponseObject, error) {
	cert, err := h.ops.GetCertByHash(ctx, req.Hash)
	if err != nil {
		return nil, err
	}
	resp := models.VerifyCertResponse{
		Valid:      cert.Valid,
		TemplateID: cert.TemplateID,
		Cert:       cert,
	}
	return openapi.VerifyCertification200JSONResponse{
		Code: 200, Data: &resp, Timestamp: time.Now(), Module: module,
	}, nil
}
