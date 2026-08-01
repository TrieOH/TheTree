package wallets

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) ListOrganizationWallets(ctx context.Context, req openapi.ListOrganizationWalletsRequestObject) (openapi.ListOrganizationWalletsResponseObject, error) {
	wallets, err := h.ops.ListFromOrg(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationWallets201JSONResponse{
		Code: 201, Data: &wallets, Timestamp: time.Now(), Module: module,
	}, nil
}
