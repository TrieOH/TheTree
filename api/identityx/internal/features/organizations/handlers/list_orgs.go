package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// ListOrgs godoc
// @Summary Lists the organizations you have access to
// @Description Lists the organizations you have access to, this includes your own and the ones you were invited too
// @Tags namespaces
// @ID namespaces_listorgs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.Organization}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /organizations [get]
func (h *Handlers) ListOrgs(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	namespaces, err := h.queries.ListOrgs(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, namespaces)
}
