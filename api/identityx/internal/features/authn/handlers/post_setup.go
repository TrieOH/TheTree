package handlers

import (
	"context"
	"time"

	"IdentityX/errs"
	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) PostSetup(ctx context.Context, req openapi.PostSetupRequestObject) (openapi.PostSetupResponseObject, error) {
	done, err := h.setup.Check(ctx)
	if err != nil {
		return nil, err
	}
	if done {
		return nil, errs.ErrSetupComplete
	}
	err = h.ops.Setup(ctx, models.SetupInput{
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
	h.setup.MarkComplete()
	return openapi.PostSetup201JSONResponse{
		Code: 201, Data: tokens, Timestamp: time.Now(), Module: module,
	}, nil
}
