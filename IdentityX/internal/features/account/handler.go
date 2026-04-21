package account

import (
	"IdentityX/internal/shared/validation"
	"net/http"

	_ "IdentityX/internal/shared/contracts"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
)

type Handler struct {
	accounts CommandService
}

func NewHandler(
	accounts CommandService,
) *Handler {
	return &Handler{
		accounts: accounts,
	}
}

// Verify godoc
// @Summary Verify user email
// @Description Verifies a user's email address using a verification token.
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification Token"
// @Success 200 {object} object "User verified successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Missing or invalid token"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/verify [get]
func (handler *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	token, rs := validation.GetString(r, "token")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.accounts.Verify(ctx, token)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("user verified, please refresh").Send(w)
}

// ResendVerificationEmail godoc
// @Summary Resend verification email
// @Description Resends the email verification link to the currently authenticated user.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} object "Verification email resent successfully"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/verify/resend [post]
func (handler *Handler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := handler.accounts.ResendVerificationEmail(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("verification email resent successfully").Send(w)
}

type ForgotPasswordRequest struct {
	Email     string     `json:"email" validate:"required,email"`
	ProjectID *uuid.UUID `json:"project_id"`
}

// ForgotPassword godoc
// @Summary Initiate password reset
// @Description Sends a password reset email to the user if the account exists.
// @Tags auth
// @Accept json
// @Produce json
// @Param forgotInfo body ForgotPasswordRequest true "User email and optional project ID"
// @Success 200 {object} object "Forgot password email sent"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/forgot-password [post]
func (handler *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := ForgotPasswordInput{
		Email:     req.Email,
		ProjectID: req.ProjectID,
	}

	ctx := r.Context()
	err := handler.accounts.ForgotPassword(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("forgot password email sent").Send(w)
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,passwd,min=8,max=72"`
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Resets the user's password using a valid reset token.
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Reset password token"
// @Param resetInfo body ResetPasswordRequest true "New password information"
// @Success 200 {object} object "Password reset successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input or token"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: Invalid or expired token"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/reset-password [post]
func (handler *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	token, rs := validation.GetString(r, "token")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req ResetPasswordRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := ResetPasswordInput{
		Token:       token,
		NewPassword: req.NewPassword,
	}

	ctx := r.Context()
	err := handler.accounts.ResetPassword(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("password reset successfully").Send(w)
}
