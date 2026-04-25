package forms

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

func NewFormsHandler(
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
		r.Get("/projects/{project_id}/forms", h.List)
		r.Post("/projects/{project_id}/forms", h.Create)
	})
}

type CreateFormRequest struct {
	Title string `json:"title" validate:"required"`
}

// Create godoc
// @Summary Create a form
// @Description Creates a form in the given project.
// @Tags forms
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param project_id path string true "Project ID"
// @Param request body CreateFormRequest true "Form title"
// @Success 201 {object} contracts.Form "Form created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /projects/{project_id}/forms [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)

	projectID, err := req.Path("project_id").UUID()
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	var payload CreateFormRequest
	if err = bind.Body(req).Bind(&payload); err != nil {
		fun.Error(err).Send(w)
		return
	}

	form, err := h.commands.Create(r.Context(), payload.Title, projectID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(form).Send(w)
}

// List godoc
// @Summary List Forms
// @Description Lists all Forms for the given project
// @Tags forms
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param project_id path string true "Project ID"
// @Success 200 {array} contracts.Form "Forms retrieved successfully"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /projects/{project_id}/forms [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	forms, err := h.queries.List(r.Context(), projectID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(forms).Send(w)
}
