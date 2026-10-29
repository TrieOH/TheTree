package handlers

import (
	"context"

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
	return openapi.CreateWallet201JSONResponse(*wallet), nil
}
