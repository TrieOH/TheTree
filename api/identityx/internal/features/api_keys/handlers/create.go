package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// Create godoc
// @Summary Create an api key in a project
// @Tags apikeys
// @ID apikeys_create
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path uuid.UUID true "Project ID"
// @Param request body models.CreateApiKeyRequest true "Api key creation data"
// @Success 200 {object} fun.Response{data=models.CreateApiKeyResponse} "Api key data"
// @Failure 401 {object} fun.Response "Unauthorized"
// @Failure 404 {object} fun.Response "Bad Request"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /projects/{project_id}/api_keys [post]
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
	var payload models.CreateApiKeyRequest
	if fun.BailInto(w, req, &payload) {
		return
	}
	key, rawKey, err := h.commands.Create(r.Context(), payload.ToInput(&projectID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, models.CreateApiKeyResponse{
		Key:    key,
		RawKey: rawKey,
	}, http.StatusCreated)
}
