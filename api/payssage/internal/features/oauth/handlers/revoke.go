package handlers

import (
	"net/http"
	"payssage/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) Revoke(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	provider, err := req.Path("provider").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	var payload models.RevokeRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	if !payload.Flow.IsValid() {
		fun.BadRequest("invalid flow, either collector or seller").Send(w)
		return
	}
	if fun.Bail(w, h.commands.Revoke(r.Context(), payload.ToInput(provider))) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
