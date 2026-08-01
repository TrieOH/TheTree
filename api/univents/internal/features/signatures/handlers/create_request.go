package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) CreateRequest(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateSignatureRequestRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	request, err := h.ops.CreateRequest(r.Context(), payload.ToInput(editionID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, request, http.StatusCreated)
}
