package projects

import (
	"net/http"

	_ "Informd/internal/shared/contracts"

	"github.com/MintzyG/FastUtilitiesNet"
	"github.com/MintzyG/FastUtilitiesNet/bind"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewProjectHandler(
	commands *CommandService,
	queries *QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handler,
	jwt func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwt)
		r.Get("/projects", h.List)
		r.Post("/projects", h.Create)
	})
}

type CreateProjectRequest struct {
	Name string `json:"name"`
}

// Create godoc
// @Summary Create a project
// @Description Creates a new project for the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body CreateProjectRequest true "Project details"
// @Success 201 {object} contracts.Project "Project created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /projects [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)

	var payload CreateProjectRequest
	if err := bind.Body(req).Bind(&payload); err != nil {
		fun.Error(err).Send(w)
		return
	}

	project, err := h.commands.Create(r.Context(), payload.Name)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(project).Send(w)
}

// List godoc
// @Summary List projects
// @Description Lists all projects owned by the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Success 200 {array} contracts.Project "Projects retrieved successfully"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /projects [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	projects, err := h.queries.List(r.Context())
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(projects).Send(w)
}
