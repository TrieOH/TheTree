package authn

import (
	"context"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetJWKS(ctx context.Context, req openapi.GetJWKSRequestObject) (openapi.GetJWKSResponseObject, error) {
	jwks, err := h.ops.JWKS(ctx, req.Params.ProjectId)
	if err != nil {
		return nil, err
	}
	keys, _ := jwks["keys"].([]map[string]any)
	return openapi.GetJWKS200JSONResponse{Keys: &keys}, nil
}
