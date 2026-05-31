package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// ListArchivedForms godoc
// @Summary Lists the namespace archived forms
// @Description Gets a list of all the archived forms created in the namespace
// @Tags namespaces
// @ID namespaces_listarchivedforms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.Form}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms [get]
func (h *Handlers) ListArchivedForms(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	forms, err := h.queries.ListArchivedForms(r.Context(), namespaceID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, forms)
}
