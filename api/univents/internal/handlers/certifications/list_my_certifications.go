package certifications

import (
	idx "sdk/identityx"

	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListMyCertifications(ctx context.Context, _ openapi.ListMyCertificationsRequestObject) (openapi.ListMyCertificationsResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	certs, err := h.ops.ListCertsByUser(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.ListMyCertifications200JSONResponse{
		Code: 200, Data: &certs, Timestamp: time.Now(), Module: module,
	}, nil
}
