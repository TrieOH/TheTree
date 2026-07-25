package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) Setup(w http.ResponseWriter, r *http.Request) {
	if globals.SetupComplete() {
		fun.Forbidden("setup already complete").Send(w)
		return
	}
	req := fun.From(r)
	ctx := r.Context()
	var payload models.IDXLoginRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err := h.commands.Setup(ctx, payload.ToSetupInput())
	if fun.Bail(w, err) {
		return
	}
	tokens, err := h.commands.Login(ctx, payload.ToInput(nil))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, tokens, http.StatusCreated)
	globals.MarkSetupComplete()
}
