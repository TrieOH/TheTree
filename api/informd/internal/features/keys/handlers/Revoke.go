package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// Revoke godoc
// @Summary Revoke an API key
// @Description Revokes the given API key, immediately invalidating it
// @Tags api_keys
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param project_id path string true "Project ID"
// @Param id path string true "API key ID"
// @Success 200 {object} fun.Response "Key revoked"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /api-keys/{id} [delete]
func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	keyID, err := req.Path("id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = h.commands.RevokeAPIKey(r.Context(), keyID)
	if fun.Bail(w, err) {
		return
	}
	fun.OK("key revoked").Send(w)
}
