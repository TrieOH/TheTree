package handlers

import (
	"net/http"
	"payssage/models"

	"github.com/MintzyG/fun"
)

// Create godoc
// @Summary Create an organization
// @Description Creates an organization with the subject as the owner, only allows IDX Clients
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
	var payload models.CreateOrganizationRequest
	if fun.BailInto(w, fun.From(r), &payload) {
		return
	}
	org, err := h.ops.Create(r.Context(), payload.ToInput())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, org, http.StatusCreated)
}
