// Package authn implements the StrictServerInterface methods for the authn
// feature. The two setup routes manage the setup flag themselves; every
// other route is covered by the setup guard in the auth dispatch
// middleware.
package authn

import (
	"context"
	"strings"
	"time"

	"IdentityX/internal/handler"
	"IdentityX/internal/services"
	"IdentityX/models"
	"lib/globals"

	"github.com/MintzyG/fun"
)

const module = "IdentityX"

type Handlers struct {
	ops *services.Authn
}

func New(ops *services.Authn) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) GetSetup(_ context.Context, _ handler.GetSetupRequestObject) (handler.GetSetupResponseObject, error) {
	if globals.SetupComplete() {
		return nil, fun.Err("setup already complete").Conflict()
	}
	return handler.GetSetup204Response{}, nil
}

func (h *Handlers) PostSetup(ctx context.Context, req handler.PostSetupRequestObject) (handler.PostSetupResponseObject, error) {
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
	return handler.PostSetup201JSONResponse{
		Code: 201, Data: tokens, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PostRegister(ctx context.Context, req handler.PostRegisterRequestObject) (handler.PostRegisterResponseObject, error) {
	err := h.ops.Register(ctx, models.IDXRegisterInput{
		Email:     req.Body.Email,
		Password:  req.Body.Password,
		ProjectID: req.Params.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return handler.PostRegister201Response{}, nil
}

func (h *Handlers) PostLogin(ctx context.Context, req handler.PostLoginRequestObject) (handler.PostLoginResponseObject, error) {
	tokens, err := h.ops.Login(ctx, models.IDXLoginInput{
		Email:     req.Body.Email,
		Password:  req.Body.Password,
		ProjectID: req.Params.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return handler.PostLogin200JSONResponse{
		Code: 200, Data: tokens, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PostLogout(ctx context.Context, req handler.PostLogoutRequestObject) (handler.PostLogoutResponseObject, error) {
	accessToken, found := strings.CutPrefix(req.Params.Authorization, "Bearer ")
	if !found {
		return nil, fun.ErrUnauthorized("invalid access token")
	}
	err := h.ops.Logout(ctx, models.LogoutInput{
		AccessToken:  accessToken,
		RefreshToken: req.Params.RefreshToken,
	})
	if err != nil {
		return nil, err
	}
	return handler.PostLogout200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PostRefresh(ctx context.Context, req handler.PostRefreshRequestObject) (handler.PostRefreshResponseObject, error) {
	tokens, err := h.ops.Refresh(ctx, req.Params.RefreshToken)
	if err != nil {
		return nil, err
	}
	return handler.PostRefresh200JSONResponse{
		Code: 200, Data: tokens, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetOAuthConnect(ctx context.Context, req handler.GetOAuthConnectRequestObject) (handler.GetOAuthConnectResponseObject, error) {
	url, err := h.ops.OAuthConnect(ctx, string(req.Provider))
	if err != nil {
		return nil, err
	}
	return handler.GetOAuthConnect200JSONResponse{
		Code: 200,
		Data: &struct {
			Url string `json:"url"` //nolint:revive // generated field name
		}{Url: url},
		Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetOAuthCallback(ctx context.Context, req handler.GetOAuthCallbackRequestObject) (handler.GetOAuthCallbackResponseObject, error) {
	if req.Params.Code == "" {
		return nil, fun.ErrBadRequest("missing code")
	}
	tokens, err := h.ops.OAuthCallback(ctx, string(req.Provider), req.Params.Code)
	if err != nil {
		return nil, err
	}
	return handler.GetOAuthCallback201JSONResponse{
		Code: 201, Data: tokens, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetJWKS(ctx context.Context, req handler.GetJWKSRequestObject) (handler.GetJWKSResponseObject, error) {
	jwks, err := h.ops.JWKS(ctx, req.Params.ProjectId)
	if err != nil {
		return nil, err
	}
	keys, _ := jwks["keys"].([]map[string]any)
	return handler.GetJWKS200JSONResponse{Keys: &keys}, nil
}

func (h *Handlers) GetIntrospect(ctx context.Context, _ handler.GetIntrospectRequestObject) (handler.GetIntrospectResponseObject, error) {
	identity, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	return handler.GetIntrospect200JSONResponse{
		Code: 200, Data: identity, Timestamp: time.Now(), Module: module,
	}, nil
}
