// Package wallets implements the StrictServerInterface methods for the
// wallets feature. NOTE: several read endpoints return 201 (a quirk of
// the current implementation) — documented as-is in the spec.
package wallets

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/internal/services"
	"payssage/models"
)

const module = "Payssage"

type Handlers struct {
	ops *services.Wallets
}

func New(ops *services.Wallets) *Handlers { return &Handlers{ops: ops} }

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

func (h *Handlers) ListWallets(ctx context.Context, _ openapi.ListWalletsRequestObject) (openapi.ListWalletsResponseObject, error) {
	wallets, err := h.ops.List(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListWallets201JSONResponse{
		Code: 201, Data: &wallets, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetWallet(ctx context.Context, req openapi.GetWalletRequestObject) (openapi.GetWalletResponseObject, error) {
	wallet, err := h.ops.GetByID(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWallet201JSONResponse{
		Code: 201, Data: wallet, Timestamp: time.Now(), Module: module,
	}, nil
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
	return openapi.SetWalletFee200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
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
	return openapi.SetWalletSandbox200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListOrganizationWallets(ctx context.Context, req openapi.ListOrganizationWalletsRequestObject) (openapi.ListOrganizationWalletsResponseObject, error) {
	wallets, err := h.ops.ListFromOrg(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationWallets201JSONResponse{
		Code: 201, Data: &wallets, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) BindCollector(ctx context.Context, req openapi.BindCollectorRequestObject) (openapi.BindCollectorResponseObject, error) {
	err := h.ops.BindCollector(ctx, req.WalletId, req.Body.CollectorId)
	if err != nil {
		return nil, err
	}
	return openapi.BindCollector204Response{}, nil
}

func (h *Handlers) UnbindCollector(ctx context.Context, req openapi.UnbindCollectorRequestObject) (openapi.UnbindCollectorResponseObject, error) {
	err := h.ops.UnbindCollector(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.UnbindCollector204Response{}, nil
}
