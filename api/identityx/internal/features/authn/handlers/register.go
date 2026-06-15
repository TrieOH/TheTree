package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"lib/telemetry"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/MintzyG/fun/middlewares"
	"go.uber.org/zap"
)

// Register godoc
// @Summary registers a user to IDX
// @Description This route is disabled until setup is complete
// @Tags authn
// @ID authn_register
// @Accept json
// @Produce json
// @Param project_id query uuid.UUID false "Project ID"
// @Param request body models.IDXRegisterRequest true "register details"
// @Success 201 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /auth/register [post]
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID := middlewares.QueryParams[models.ProjectIDQueryParam](r)
	telemetry.DLog().Info("Login", zap.Any("projectID", projectID.ProjectID))
	var payload models.IDXRegisterRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err := h.commands.Register(r.Context(), payload.ToInput(projectID.ProjectID))
	if fun.Bail(w, err) {
		return
	}
	fun.Created().Send(w)
}
