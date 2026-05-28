package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// List godoc
// @Summary list steps in a form
// @Tags steps
// @ID steps_list
// @Accept json
// @Produce json
// @Param form_id path string true "Form ID"
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/steps [get]
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	steps, err := h.queries.List(r.Context(), formID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, steps)
}
