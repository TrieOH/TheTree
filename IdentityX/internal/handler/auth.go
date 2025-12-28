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

	if err := h.AuthService.Register(r.Context(), req); err != nil {
		ErrToResp(err).Send(w)
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

	tokens, err := h.AuthService.Login(r, r.Context(), req)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	accessToken, err := utils.ParseAccessToken(tokens.AccessTokenString, utils.GoAuthPublicKey)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	refreshToken, err := utils.ParseRefreshToken(tokens.RefreshTokenString, utils.GoAuthPublicKey)
	if err != nil {
		ErrToResp(err).Send(w)
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
	err := h.AuthService.Logout(r.Context())
	if err != nil {
		ErrToResp(err).Send(w)
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
// @Param Cookie header string true "Cookie: refresh_token=yyy"
// @Param registerInfo body models.RegisterUserRequest true "register request data"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Success 200 {string} string "Refreshed tokens"
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		resp.Unauthorized("error getting refresh token").AddTrace(err).Send(w)
		return
	}

	if refreshTokenCookie.Value == "" {
		resp.Unauthorized("missing refresh token value").Send(w)
		return
	}

	var data models.RefreshData
	data.RefreshCookie = refreshTokenCookie
	data.Agent = r.UserAgent()
	data.IP = utils.GetClientIP(r)

	tokens, err := h.AuthService.Refresh(data, r.Context())
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	accessToken, err := utils.ParseAccessToken(tokens.AccessTokenString, utils.GoAuthPublicKey)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	refreshToken, err := utils.ParseRefreshToken(tokens.RefreshTokenString, utils.GoAuthPublicKey)
	if err != nil {
		ErrToResp(err).Send(w)
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

	resp.OK("Refreshed tokens").Send(w)
}

// Me godoc
// @Summary Sends session info to user
// @Description sends the whole access info and the refresh expiry time
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Header 200 {string} Set-Cookie "access_token cookie for authentication"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Success 200 {object} map[string]any
// @Failure 500 {object} models.ErrorResponse
// @Router /sessions/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	access, err := models.GetAccessClaims(r.Context())
	if err != nil {
		resp.InternalServerError("Failed to get access claims").AddTrace(err).Send(w)
		return
	}
	refresh, err := models.GetRefreshClaims(r.Context())
	if err != nil {
		resp.InternalServerError("Failed to get refresh claims").AddTrace(err).Send(w)
		return
	}

	resp.OK().WithData(map[string]interface{}{
		"access":              access,
		"refresh_expire_date": refresh.ExpiresAt,
	}).Send(w)
}

// JWKS godoc
// @Summary Exposes the public key using a JWKS
// @Description Lets users verify the tokens using the public key through a JWKS
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Router /.well-known/jwks.json [get]
func (h *AuthHandler) JWKS(w http.ResponseWriter, _ *http.Request) {
	jwks := map[string]any{
		"keys": []any{utils.PublicKeyToJWK(utils.GoAuthPublicKey)},
	}

	resp.OK().WithData(jwks).Send(w)
}
