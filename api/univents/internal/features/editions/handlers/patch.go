package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) Patch(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.PatchEditionRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	edition, err := h.ops.Patch(r.Context(), payload.ToInput(editionID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, edition)
}
