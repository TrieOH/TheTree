package account

import (
	"IdentityX/internal/shared/contracts"
	"net/http"
	"time"

	_ "IdentityX/internal/shared/contracts"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
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

type TokenParam struct {
	Token string `fun_query:"token,required"`
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handler,
	jwt func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(httprate.Limit(5, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByRealIP)))

		r.Post("/account/forgot-password", h.ForgotPassword)
		r.With(jwt).Post("/account/verify/resend", h.ResendVerificationEmail)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.WithParams[TokenParam](true))
			r.Post("/account/reset-password", h.ResetPassword)
			r.With(jwt).Post("/account/verify", h.Verify)
		})
	})
}

// Verify godoc
// @Summary Verify user email
// @Description Verifies a user's email address using a verification token.
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification Token"
// @Success 200 {string} string "User verified successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Missing or invalid token"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/verify [get]
func (handler *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	token := middlewares.QueryParams[TokenParam](r).Token
	err := handler.accounts.Verify(r.Context(), token)
	if fun.Bail(w, err) {
		return
	}
	fun.OK("user verified, please refresh").Send(w)
}

// ResendVerificationEmail godoc
// @Summary Resend verification email
// @Description Resends the email verification link to the currently authenticated user.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {string} string "Verification email resent successfully"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/verify/resend [post]
func (handler *Handler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	err := handler.accounts.ResendVerificationEmail(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.OK("verification email resent successfully").Send(w)
}

// ForgotPassword godoc
// @Summary Initiate password reset
// @Description Sends a password reset email to the user if the account exists.
// @Tags auth
// @Accept json
// @Produce json
// @Param forgotInfo body contracts.ForgotPasswordRequest true "User email and optional project ID"
// @Success 200 {string} string "Forgot password email sent"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/forgot-password [post]
func (handler *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	var payload contracts.ForgotPasswordRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err := handler.accounts.ForgotPassword(r.Context(), payload.ToInput())
	if fun.Bail(w, err) {
		return
	}
	fun.OK("forgot password email sent").Send(w)
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Resets the user's password using a valid reset token.
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Reset password token"
// @Param resetInfo body contracts.ResetPasswordRequest true "New password information"
// @Success 200 {string} string "Password reset successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input or token"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: Invalid or expired token"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/reset-password [post]
func (handler *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	token := middlewares.QueryParams[TokenParam](r).Token
	var payload contracts.ResetPasswordRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err := handler.accounts.ResetPassword(r.Context(), payload.ToInput(token))
	if fun.Bail(w, err) {
		return
	}
	fun.OK("password reset successfully").Send(w)
}
