package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// GetPlatformProfile godoc
// @Summary Get a platform-scoped actor's profile
// @Tags profiles
// @ID profiles_getplatformprofile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param actor_id path string true "Actor ID"
// @Success 200 {object} fun.Response{data=models.ActorProfile}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /actors/{actor_id}/profile [get]
func (h *Handlers) GetPlatformProfile(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	actorID, err := req.Path("actor_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	profile, err := h.queries.GetPlatformProfile(r.Context(), actorID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, profile)
}
