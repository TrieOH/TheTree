package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// ResponseCount godoc
// @Summary Gets the number of responses of the form
// @Tags forms
// @ID forms_responsecount
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param form_id path string true "Form ID"
// @Success 200 {object} fun.Response "Form updated successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/responses/count [get]
func (h *Handlers) ResponseCount(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	count, err := h.queries.GetResponseCount(r.Context(), formID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, map[string]int{"count": count})
}
