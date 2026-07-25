package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/MintzyG/fun/middlewares"
)

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID := middlewares.QueryParams[models.ProjectIDQueryParam](r)
	var payload models.IDXLoginRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	tokens, err := h.commands.Login(r.Context(), payload.ToInput(projectID.ProjectID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, tokens)
}
