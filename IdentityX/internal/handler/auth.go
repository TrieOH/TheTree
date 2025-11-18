package handler

import (
	"GoAuth/internal/utils"
	"net/http"
	"time"

	"GoAuth/internal/models"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/FastUtilitiesNet/validation"
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
	if rs := validation.ValidateInto(r, &req); rs != nil {
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
// @Description Authenticates a customer of the system
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
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	tokens, rs := h.AuthService.Login(r, r.Context(), req)
	if rs != nil {
		rs.Send(w)
		return
	}

	accessToken, rs := utils.ParseAccessToken(tokens.AccessTokenString, utils.GoAuthPublicKey)
	if rs != nil {
		rs.Send(w)
		return
	}

	refreshToken, rs := utils.ParseRefreshToken(tokens.RefreshTokenString, utils.GoAuthPublicKey)
	if rs != nil {
		rs.Send(w)
		return
	}

	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessTokenString,
		Path:     "/",
		MaxAge:   int(time.Until(accessToken.ExpiresAt.Time).Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshTokenString,
		Path:     "/",
		MaxAge:   int(time.Until(refreshToken.ExpiresAt.Time).Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)

	resp.OK("Logged in").WithData(map[string]interface{}{
		"access_token_claims":  accessToken,
		"refresh_token_claims": refreshToken,
	}).Send(w)
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
		SameSite: http.SameSiteStrictMode,
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)

	resp.OK("Logged out").Send(w)
}

// Refresh godoc
// @Summary Refreshes the user token pair
// @Description Creates a new token pair from a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param registerInfo body models.RegisterUserRequest true "register request data"
// @Header 200 {string} Set-Cookie "access_token cookie for authentication"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Success 200 {string} string "Refreshed tokens"
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	tokens, rs := h.AuthService.Refresh(r, r.Context())
	if rs != nil {
		rs.Send(w)
		return
	}

	accessToken, rs := utils.ParseAccessToken(tokens.AccessTokenString, utils.GoAuthPublicKey)
	if rs != nil {
		rs.Send(w)
		return
	}

	refreshToken, rs := utils.ParseRefreshToken(tokens.RefreshTokenString, utils.GoAuthPublicKey)
	if rs != nil {
		rs.Send(w)
		return
	}

	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessTokenString,
		Path:     "/",
		MaxAge:   int(time.Until(accessToken.ExpiresAt.Time).Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshTokenString,
		Path:     "/",
		MaxAge:   int(time.Until(refreshToken.ExpiresAt.Time).Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)

	resp.OK("Refreshed tokens").WithData(map[string]interface{}{
		"access_token_claims":  accessToken,
		"refresh_token_claims": refreshToken,
	}).Send(w)
}
