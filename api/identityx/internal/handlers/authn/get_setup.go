package authn

import (
	"context"

	"github.com/MintzyG/fun"

	"IdentityX/internal/openapi"
	"IdentityX/internal/setup"
)

func (h *Handlers) GetSetup(_ context.Context, _ openapi.GetSetupRequestObject) (openapi.GetSetupResponseObject, error) {
	if setup.Complete() {
		return nil, fun.Err("setup already complete").Conflict()
	}
	return openapi.GetSetup204Response{}, nil
}
