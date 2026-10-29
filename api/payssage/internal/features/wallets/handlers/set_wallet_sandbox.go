package handlers

import (
	"context"

	"payssage/internal/openapi"
	"payssage/models"
)

func (h *Handlers) SetWalletSandbox(ctx context.Context, req openapi.SetWalletSandboxRequestObject) (openapi.SetWalletSandboxResponseObject, error) {
	err := h.ops.SetSandbox(ctx, models.SetSandboxInput{
		WalletID:       req.WalletId,
		Sandbox:        req.Body.Sandbox,
		OrganizationID: req.Body.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.SetWalletSandbox204Response{}, nil
}
