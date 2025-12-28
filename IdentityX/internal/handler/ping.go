package handler

import (
	"net/http"

	"GoAuth/internal/utils"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// PublicPing godoc
// @Summary Just replies "pong"
// @Description This route replies pong to any request to test connectivity
// @Description This is not meant to be used as a health check, but it can be trusted
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Router /ping/public [post]
func (h *AuthHandler) PublicPing(w http.ResponseWriter, _ *http.Request) {
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
	accessToken, err := r.Cookie("access_token")
	if err != nil {
		resp.Unauthorized("missing access_token cookie").Send(w)
		return
	}

	accessClaims, err := utils.ParseAccessToken(accessToken.Value, utils.GoAuthPublicKey)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK("pong to you " + accessClaims.Sub.Email).Send(w)
}
