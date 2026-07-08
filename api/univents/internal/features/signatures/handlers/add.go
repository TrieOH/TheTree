package handlers

import (
	"net/http"
	"univents/contracts"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) Add(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload contracts.AddSignatureRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	sig, err := h.commands.Add(r.Context(), payload.ToInput(editionID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, sig)
}
