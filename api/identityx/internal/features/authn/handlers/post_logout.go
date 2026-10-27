package handlers

import (
	"context"
	"strings"
	"time"

	"IdentityX/errs"
	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) PostLogout(ctx context.Context, req openapi.PostLogoutRequestObject) (openapi.PostLogoutResponseObject, error) {
	accessToken, found := strings.CutPrefix(req.Params.Authorization, "Bearer ")
	if !found {
		return nil, errs.ErrInvalidAccessToken
	}
	err := h.ops.Logout(ctx, models.LogoutInput{
		AccessToken:  accessToken,
		RefreshToken: req.Params.RefreshToken,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PostLogout200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
