package handlers

import (
	"context"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) PostLogin(ctx context.Context, req openapi.PostLoginRequestObject) (openapi.PostLoginResponseObject, error) {
	tokens, err := h.ops.Login(ctx, models.IDXLoginInput{
		Email:     req.Body.Email,
		Password:  req.Body.Password,
		ProjectID: req.Params.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PostLogin200JSONResponse(*tokens), nil
}
