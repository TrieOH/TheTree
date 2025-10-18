package handler

import (
	"fmt"
	"net/http"

	"GoAuth/internal/models"
	"GoAuth/internal/validation"
	"GoAuth/internal/utils"

	"github.com/spf13/viper"
	resp "github.com/MintzyG/GoResponse/response"
)

// Register godoc
// @Summary Register a new customer
// @Description registers a customer into the system
// @Tags auth
// @Accept json
// @Produce json
// @Param registerInfo body models.RegisterUserRequest true "register request data"
// @Success 201 {string} string "Registered user"
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterUserRequest
	if rs := validation.ValidateWith(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	if rs := h.AuthService.Register(r.Context(), req); rs != nil {
		rs.Send(w)
		return
	}

	resp.Created("Registered user").Send(w)
}

// Login godoc
// @Summary Authenticates a customer
// @Description Autheticates a customer of the system
// @Tags auth
// @Accept json
// @Produce json
// @Param loginInfo body models.LoginUserRequest true "login request data"
// @Success 200 {string} string "Logged in"
// @Header 200 {string} Set-Cookie "access_token cookie for authentication"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Success 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginUserRequest
	if rs := validation.ValidateWith(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	tokens, rs := h.AuthService.Login(r, r.Context(), req)
	if rs != nil {
		rs.Send(w)
		return
	}

	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessTokenString,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshTokenString,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)

	resp.OK("Logged in").Send(w)
}

// Logout godoc
// @Summary Logs out a customer
// @Description Logs out an authenticated customer of the system
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {string} string "Logged out"
// @Header 200 {string} Set-Cookie "clears the access_token cookie"
// @Header 200 {string} Set-Cookie "clears the refresh_token cookie"
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	rs := h.AuthService.Logout(r, r.Context())
	if rs != nil {
		rs.Send(w)
		return
	}

	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)

	resp.OK("Logged out").Send(w)
}

// PublicPing godoc
// @Summary Just replies "pong"
// @Description This route replies pong to any request to test connectivity
// @Description This is not meant to be used as a health check but it can be trusted
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Router /ping/public [post]
func (h *AuthHandler) PublicPing(w http.ResponseWriter, r *http.Request) {
	resp.OK("pong").Send(w)
}

// PrivatePing godoc
// @Summary Just replies "pong to you {EMAIL}"
// @Description This route replies pong with the authenticated user email
// @Description You must be authenticated to use this route
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {string} string
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /ping/private [post]
func (h *AuthHandler) PrivatePing(w http.ResponseWriter, r *http.Request) {
	access_token, err := r.Cookie("access_token")
	if err != nil {
		resp.Unauthorized("missing access_token cookie").Send(w)
		return
	}

	accessClaims, rs := utils.ParseAccessToken(access_token.Value, viper.GetString("JWT_SECRET"))
	if rs != nil {
		rs.Send(w)
		return
	}

	resp.OK("pong to you " + accessClaims.Sub.Email).Send(w)
}

// ListUserSessions godoc
// @Summary Lists all active user sessions
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {string} string
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /sessions [get]
func (h *AuthHandler) ListUserSessions(w http.ResponseWriter, r *http.Request) {
	accessClaims, ok := models.GetAccessClaims(r)
	resp.OK(fmt.Sprintf("%v: %v", ok, accessClaims.Sub.Email)).Send(w)
}
