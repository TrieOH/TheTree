package handlers

import (
	idx "sdk/identityx"

	"context"

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
	return openapi.ListMyCertifications200JSONResponse{Items: certs}, nil
}
