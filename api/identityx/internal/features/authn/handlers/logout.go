package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"
	"strings"

	"github.com/MintzyG/fun"
)

// Logout godoc
// @Summary logs out a user from IDX
// @Description This route is disabled until setup is complete
// @Tags authn
// @ID authn_logout
// @Accept json
// @Produce json
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /auth/logout [post]
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	authorization := req.Header("Authorization").String()
	accessToken, found := strings.CutPrefix("Bearer ", authorization)
	if !found {
		fun.Error(fun.ErrUnauthorized("Invalid access token"))
		return
	}
	refreshToken := req.Header("refresh_token").String()
	err := h.commands.Logout(r.Context(), models.LogoutInput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
