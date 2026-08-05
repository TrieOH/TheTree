package wallets

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/models"
)

func (h *Handlers) CreateWallet(ctx context.Context, req openapi.CreateWalletRequestObject) (openapi.CreateWalletResponseObject, error) {
	wallet, err := h.ops.Create(ctx, models.CreateWalletInput{
		Name:           req.Body.Name,
		OrganizationID: req.Body.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateWallet201JSONResponse{
		Code: 201, Data: wallet, Timestamp: time.Now(), Module: module,
	}, nil
}
