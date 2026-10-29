package authn

import (
	"context"
	"strings"

	"IdentityX/errs"
	"IdentityX/internal/features/authn/service"
	"IdentityX/internal/openapi"
	"IdentityX/internal/setup"
	"IdentityX/models"
)

type Handlers struct {
	ops   *service.Operations
	setup *setup.State
}

func NewHandlers(ops *service.Operations, st *setup.State) *Handlers {
	return &Handlers{ops: ops, setup: st}
}

func (h *Handlers) GetIntrospect(ctx context.Context, _ openapi.GetIntrospectRequestObject) (openapi.GetIntrospectResponseObject, error) {
	identity, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.GetIntrospect200JSONResponse(*identity), nil
}

func (h *Handlers) GetJWKS(ctx context.Context, req openapi.GetJWKSRequestObject) (openapi.GetJWKSResponseObject, error) {
	jwks, err := h.ops.JWKS(ctx, req.Params.ProjectId)
	if err != nil {
		return nil, err
	}
	keys, _ := jwks["keys"].([]map[string]any)
	return openapi.GetJWKS200JSONResponse{Keys: &keys}, nil
}

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

func (h *Handlers) PostVerifyEmail(ctx context.Context, req openapi.PostVerifyEmailRequestObject) (openapi.PostVerifyEmailResponseObject, error) {
	err := h.ops.VerifyEmail(ctx, models.VerifyEmailInput{Token: req.Body.Token})
	if err != nil {
		return nil, err
	}
	return openapi.PostVerifyEmail204Response{}, nil
}

func (h *Handlers) PostResendVerification(ctx context.Context, req openapi.PostResendVerificationRequestObject) (openapi.PostResendVerificationResponseObject, error) {
	err := h.ops.ResendVerification(ctx, models.ResendVerificationInput{
		Email:     req.Body.Email,
		ProjectID: req.Body.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PostResendVerification204Response{}, nil
}

func (h *Handlers) PostForgotPassword(ctx context.Context, req openapi.PostForgotPasswordRequestObject) (openapi.PostForgotPasswordResponseObject, error) {
	err := h.ops.ForgotPassword(ctx, models.ForgotPasswordInput{
		Email:     req.Body.Email,
		ProjectID: req.Body.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PostForgotPassword204Response{}, nil
}

func (h *Handlers) PostResetPassword(ctx context.Context, req openapi.PostResetPasswordRequestObject) (openapi.PostResetPasswordResponseObject, error) {
	err := h.ops.ResetPassword(ctx, models.ResetPasswordInput{
		Token:       req.Body.Token,
		NewPassword: req.Body.Password,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PostResetPassword204Response{}, nil
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
	return openapi.PostLogin200JSONResponse(*tokens), nil
}

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
	return openapi.PostLogout204Response{}, nil
}

func (h *Handlers) PostRefresh(ctx context.Context, req openapi.PostRefreshRequestObject) (openapi.PostRefreshResponseObject, error) {
	tokens, err := h.ops.Refresh(ctx, req.Params.RefreshToken)
	if err != nil {
		return nil, err
	}
	return openapi.PostRefresh200JSONResponse(*tokens), nil
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
	return openapi.PostSetup201JSONResponse(*tokens), nil
}
