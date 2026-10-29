package handlers

import (
	"context"

	"payssage/internal/features/wallets/service"
	"payssage/internal/openapi"
	"payssage/models"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) BindCollector(ctx context.Context, req openapi.BindCollectorRequestObject) (openapi.BindCollectorResponseObject, error) {
	err := h.ops.BindCollector(ctx, req.WalletId, req.Body.CollectorId)
	if err != nil {
		return nil, err
	}
	return openapi.BindCollector204Response{}, nil
}

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

func (h *Handlers) GetWallet(ctx context.Context, req openapi.GetWalletRequestObject) (openapi.GetWalletResponseObject, error) {
	wallet, err := h.ops.GetByID(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWallet201JSONResponse(*wallet), nil
}

func (h *Handlers) ListOrganizationWallets(ctx context.Context, req openapi.ListOrganizationWalletsRequestObject) (openapi.ListOrganizationWalletsResponseObject, error) {
	wallets, err := h.ops.ListFromOrg(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationWallets201JSONResponse{Items: wallets}, nil
}

func (h *Handlers) ListWallets(ctx context.Context, _ openapi.ListWalletsRequestObject) (openapi.ListWalletsResponseObject, error) {
	wallets, err := h.ops.List(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListWallets201JSONResponse{Items: wallets}, nil
}

func (h *Handlers) SetWalletFee(ctx context.Context, req openapi.SetWalletFeeRequestObject) (openapi.SetWalletFeeResponseObject, error) {
	err := h.ops.SetFeeBPS(ctx, models.SetFeeBPSInput{
		WalletID:       req.WalletId,
		FeeBps:         req.Body.FeeBps,
		OrganizationID: req.Body.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.SetWalletFee204Response{}, nil
}

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

func (h *Handlers) UnbindCollector(ctx context.Context, req openapi.UnbindCollectorRequestObject) (openapi.UnbindCollectorResponseObject, error) {
	err := h.ops.UnbindCollector(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.UnbindCollector204Response{}, nil
}
