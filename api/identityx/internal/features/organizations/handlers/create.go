package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// Create godoc
// @Summary Create an organization
// @Description Creates and organization with the subject as the owner, only allows IDX Clients
// @Tags organizations
// @ID organizations_create
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateOrganizationRequest true "Organization creation data"
// @Success 200 {object} models.Organization "Organization data"
// @Failure 401 {object} fun.Response "Unauthorized"
// @Failure 404 {object} fun.Response "Bad Request"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /organizations [post]
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	var payload models.CreateOrganizationRequest
	if fun.BailInto(w, fun.From(r), &payload) {
		return
	}
	org, err := h.commands.Create(r.Context(), payload.ToInput())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, org, http.StatusCreated)
}
