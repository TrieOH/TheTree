package api_keys

import (
	"IdentityX/internal/platform/middlewares"
	"net/http"

	_ "IdentityX/contracts"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	apiKeys CommandService
}

func NewHandler(
	apiKeys CommandService,
) *Handler {
	return &Handler{apiKeys: apiKeys}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handler,
	jwt func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwt)
		r.Use(middlewares.ClientOnly())
		r.Post("/projects/{project_id}/api-keys/rotate", h.RotateApiKey)
		r.Delete("/projects/{project_id}/api-keys", h.RevokeApiKey)
	})
}

// RotateApiKey godoc
// @Summary Rotate API key for a project
// @Description Generates a new API key for the project and invalidates the previous one. This is the only time the full key will be shown.
// @Tags api_keys
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {string} string ""
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized"
// @Failure 404 {object} contracts.ErrorResponse "Project not found"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/api-keys/rotate [post]
func (handler *Handler) RotateApiKey(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	apiKey, err := handler.apiKeys.Rotate(r.Context(), projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, apiKey)
}

// RevokeApiKey godoc
// @Summary Revoke API key for a project
// @Description Deletes the API key for the project, disabling programmatic access.
// @Tags api_keys
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized"
// @Failure 404 {object} contracts.ErrorResponse "Project not found"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/api-keys [delete]
func (handler *Handler) RevokeApiKey(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = handler.apiKeys.Revoke(r.Context(), projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
