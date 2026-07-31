package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"
	"strings"

	"github.com/MintzyG/fun"
)

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
	refreshToken := req.Header("Refresh-Token").String()
	err := h.commands.Logout(r.Context(), models.LogoutInput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
