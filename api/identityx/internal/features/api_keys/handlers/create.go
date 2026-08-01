package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateAPIKeyRequest
	if fun.BailInto(w, req, &payload) {
		return
	}
	key, rawKey, err := h.ops.Create(r.Context(), payload.ToInput(&projectID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, models.CreateAPIKeyResponse{
		Key:    key,
		RawKey: rawKey,
	}, http.StatusCreated)
}
