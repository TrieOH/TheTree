package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// ListMineArchived godoc
// @Summary Lists direct archived forms you own
// @Description Gets a list of all the archived forms you created directly in your user i.e. not namespaced
// @Tags forms
// @ID forms_listminearchived
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.Form}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/archived [get]
func (h *Handlers) ListMineArchived(w http.ResponseWriter, r *http.Request) {
	forms, err := h.queries.ListArchivedForms(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, forms)
}
