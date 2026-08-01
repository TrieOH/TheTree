package authn

import (
	"context"
	"time"

	"github.com/MintzyG/fun"

	"IdentityX/internal/openapi"
	"IdentityX/models"
	"lib/globals"
)

func (h *Handlers) PostSetup(ctx context.Context, req openapi.PostSetupRequestObject) (openapi.PostSetupResponseObject, error) {
	if globals.SetupComplete() {
		return nil, fun.Err("setup already complete").Conflict()
	}
	err := h.ops.Setup(ctx, models.SetupInput{
		Email:    req.Body.Email,
		Password: req.Body.Password,
	})
	if err != nil {
		return nil, err
	}
	tokens, err := h.ops.Login(ctx, models.IDXLoginInput{
		Email:    req.Body.Email,
		Password: req.Body.Password,
	})
	if err != nil {
		return nil, err
	}
	globals.MarkSetupComplete()
	return openapi.PostSetup201JSONResponse{
		Code: 201, Data: tokens, Timestamp: time.Now(), Module: module,
	}, nil
}
