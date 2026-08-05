package wallets

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/models"
)

func (h *Handlers) SetWalletFee(ctx context.Context, req openapi.SetWalletFeeRequestObject) (openapi.SetWalletFeeResponseObject, error) {
	err := h.ops.SetFeeBPS(ctx, models.SetFeeBPSInput{
		WalletID:       req.WalletId,
		FeeBps:         req.Body.FeeBps,
		OrganizationID: req.Body.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.SetWalletFee200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
