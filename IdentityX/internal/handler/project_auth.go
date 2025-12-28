package handler

import (
	"GoAuth/internal/utils"
	"net/http"
	"time"

	"GoAuth/internal/models"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/FastUtilitiesNet/validation"
)

// ProjectRegister godoc
// @Summary Register a new user into a client project
// @Description registers a user into the specified project
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to register user"
// @Param registerInfo body models.RegisterProjectUserRequest true "register project user request data"
// @Success 201 {string} string "Registered user"
// @Failure 500 {object} models.ErrorResponse
// @Router /projects/{project_id}/register [post]
func (h *AuthHandler) ProjectRegister(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("project_id")
	if projectId == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	var req models.RegisterProjectUserRequest
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	if err := h.AuthService.RegisterProjectUser(r.Context(), projectId, req); err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.Created("Registered user").Send(w)
}

// ProjectLogin godoc
// @Summary Authenticates a user into a client project
// @Description Authenticates a user into the specified client project
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to login user"
// @Param loginInfo body models.LoginProjectUserRequest true "login project user request data"
// @Success 200 {string} string "Logged in"
// @Header 200 {string} Set-Cookie "access_token cookie for authentication"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Success 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /projects/{project_id}/login [post]
func (h *AuthHandler) ProjectLogin(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("project_id")
	if projectId == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	var req models.LoginProjectUserRequest
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	tokens, err := h.AuthService.LoginProjectUser(r, r.Context(), projectId, req)
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
		HttpOnly: false,
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
