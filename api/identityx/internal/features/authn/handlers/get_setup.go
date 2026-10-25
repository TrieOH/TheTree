package handlers

import (
	"context"

	"github.com/MintzyG/fun"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetSetup(ctx context.Context, _ openapi.GetSetupRequestObject) (openapi.GetSetupResponseObject, error) {
	done, err := h.setup.Check(ctx)
	if err != nil {
		return nil, err
	}
	if done {
		return nil, fun.Err("setup already complete").Conflict()
	}
	return openapi.GetSetup204Response{}, nil
}
