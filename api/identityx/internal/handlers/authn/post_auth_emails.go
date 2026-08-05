package authn

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) PostVerifyEmail(ctx context.Context, req openapi.PostVerifyEmailRequestObject) (openapi.PostVerifyEmailResponseObject, error) {
	err := h.ops.VerifyEmail(ctx, models.VerifyEmailInput{Token: req.Body.Token})
	if err != nil {
		return nil, err
	}
	return openapi.PostVerifyEmail200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PostResendVerification(ctx context.Context, req openapi.PostResendVerificationRequestObject) (openapi.PostResendVerificationResponseObject, error) {
	err := h.ops.ResendVerification(ctx, models.ResendVerificationInput{
		Email:     req.Body.Email,
		ProjectID: req.Body.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PostResendVerification200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PostForgotPassword(ctx context.Context, req openapi.PostForgotPasswordRequestObject) (openapi.PostForgotPasswordResponseObject, error) {
	err := h.ops.ForgotPassword(ctx, models.ForgotPasswordInput{
		Email:     req.Body.Email,
		ProjectID: req.Body.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PostForgotPassword200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) PostResetPassword(ctx context.Context, req openapi.PostResetPasswordRequestObject) (openapi.PostResetPasswordResponseObject, error) {
	err := h.ops.ResetPassword(ctx, models.ResetPasswordInput{
		Token:       req.Body.Token,
		NewPassword: req.Body.Password,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PostResetPassword200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
