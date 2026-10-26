package handlers

import (
	"context"

	"IdentityX/errs"
	"IdentityX/internal/openapi"
)

func (h *Handlers) GetSetup(ctx context.Context, _ openapi.GetSetupRequestObject) (openapi.GetSetupResponseObject, error) {
	done, err := h.setup.Check(ctx)
	if err != nil {
		return nil, err
	}
	if done {
		return nil, errs.ErrSetupComplete
	}
	return openapi.GetSetup204Response{}, nil
}
