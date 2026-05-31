package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// ResponseCount godoc
// @Summary Gets the number of responses of the form
// @Tags namespaces
// @ID namespaces_formresponsecount
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param namespace_id path string true "Namespace ID"
// @Param form_id path string true "Form ID"
// @Success 200 {object} fun.Response "Form updated successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms/{form_id}/responses/count [get]
func (h *Handlers) ResponseCount(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	count, err := h.queries.GetFormResponseCount(r.Context(), namespaceID, formID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, map[string]int{"count": count})
}
