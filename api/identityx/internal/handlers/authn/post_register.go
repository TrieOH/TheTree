package authn

import (
	"context"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) PostRegister(ctx context.Context, req openapi.PostRegisterRequestObject) (openapi.PostRegisterResponseObject, error) {
	err := h.ops.Register(ctx, models.IDXRegisterInput{
		Email:     req.Body.Email,
		Password:  req.Body.Password,
		ProjectID: req.Params.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PostRegister201Response{}, nil
}
