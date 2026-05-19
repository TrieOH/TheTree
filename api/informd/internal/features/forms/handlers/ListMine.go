package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// ListMine godoc
// @Summary Lists direct forms you own
// @Description Gets a list of all the forms you created directly in your user i.e. not namespaced
// @Tags forms
// @ID forms_listmine
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.Form}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms [get]
func (h *Handlers) ListMine(w http.ResponseWriter, r *http.Request) {
	forms, err := h.queries.ListForms(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, forms)
}
