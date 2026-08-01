// Package authn implements the StrictServerInterface methods for the authn
// feature. The two setup routes manage the setup flag themselves; every
// other route is covered by the setup guard in the auth dispatch
// middleware.
package authn

import (
	"context"
	"strings"
	"time"

	"IdentityX/internal/openapi"
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

func (h *Handlers) GetSetup(_ context.Context, _ openapi.GetSetupRequestObject) (openapi.GetSetupResponseObject, error) {
	if globals.SetupComplete() {
		return nil, fun.Err("setup already complete").Conflict()
	}
	return openapi.GetSetup204Response{}, nil
}

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

func (h *Handlers) PostLogin(ctx context.Context, req openapi.PostLoginRequestObject) (openapi.PostLoginResponseObject, error) {
	tokens, err := h.ops.Login(ctx, models.IDXLoginInput{
		Email:     req.Body.Email,
		Password:  req.Body.Password,
		ProjectID: req.Params.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PostLogin200JSONResponse{
		Code: 200, Data: tokens, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PostLogout(ctx context.Context, req openapi.PostLogoutRequestObject) (openapi.PostLogoutResponseObject, error) {
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
	return openapi.PostLogout200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PostRefresh(ctx context.Context, req openapi.PostRefreshRequestObject) (openapi.PostRefreshResponseObject, error) {
	tokens, err := h.ops.Refresh(ctx, req.Params.RefreshToken)
	if err != nil {
		return nil, err
	}
	return openapi.PostRefresh200JSONResponse{
		Code: 200, Data: tokens, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetOAuthConnect(ctx context.Context, req openapi.GetOAuthConnectRequestObject) (openapi.GetOAuthConnectResponseObject, error) {
	url, err := h.ops.OAuthConnect(ctx, string(req.Provider))
	if err != nil {
		return nil, err
	}
	return openapi.GetOAuthConnect200JSONResponse{
		Code: 200,
		Data: &struct {
			Url string `json:"url"` //nolint:revive // generated field name
		}{Url: url},
		Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetOAuthCallback(ctx context.Context, req openapi.GetOAuthCallbackRequestObject) (openapi.GetOAuthCallbackResponseObject, error) {
	if req.Params.Code == "" {
		return nil, fun.ErrBadRequest("missing code")
	}
	tokens, err := h.ops.OAuthCallback(ctx, string(req.Provider), req.Params.Code)
	if err != nil {
		return nil, err
	}
	return openapi.GetOAuthCallback201JSONResponse{
		Code: 201, Data: tokens, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetJWKS(ctx context.Context, req openapi.GetJWKSRequestObject) (openapi.GetJWKSResponseObject, error) {
	jwks, err := h.ops.JWKS(ctx, req.Params.ProjectId)
	if err != nil {
		return nil, err
	}
	keys, _ := jwks["keys"].([]map[string]any)
	return openapi.GetJWKS200JSONResponse{Keys: &keys}, nil
}

func (h *Handlers) GetIntrospect(ctx context.Context, _ openapi.GetIntrospectRequestObject) (openapi.GetIntrospectResponseObject, error) {
	identity, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.GetIntrospect200JSONResponse{
		Code: 200, Data: identity, Timestamp: time.Now(), Module: module,
	}, nil
}
