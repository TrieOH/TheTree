package handler

import (
	"Informd/contracts"
	"Informd/internal/features/forms/commands"
	"Informd/internal/features/forms/queries"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handler,
	anyAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(anyAuth)
		r.Post("/forms", h.Create)
		r.With(middlewares.WithParams[contracts.BulkGetParams]()).Post("/forms/bulk", h.BulkGet)
		r.Post("/namespaces/{namespace_id}/forms", h.CreateInWorkspace)
		r.Post("/forms/{form_id}/steps", h.Create)
	})
}

// Create godoc
// @Summary Create a form
// @Description Creates a form not namespaced.
// @Tags forms
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body contracts.CreateFormRequest true "Form title"
// @Success 201 {object} contracts.Form "Form created successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	var payload contracts.CreateFormRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	form, err := h.commands.Create(r.Context(), payload.Title, nil)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, form, http.StatusCreated)
}

// CreateInWorkspace godoc
// @Summary Create a form
// @Description Creates a form in the given namespace.
// @Tags forms
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param namespace_id path string true "Namespace ID"
// @Param request body contracts.CreateFormRequest true "Form title"
// @Success 201 {object} contracts.Form "Form created successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms [post]
func (h *Handler) CreateInWorkspace(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload contracts.CreateFormRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	form, err := h.commands.Create(r.Context(), payload.Title, &namespaceID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, form, http.StatusCreated)
}

// BulkGet godoc
// @Summary Bulk get forms
// @Description Returns a list of forms by their IDs. IDs should be obtained via a SpiceDB lookup on the client side.
// @Tags forms
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body contracts.BulkGetRequest true "Form IDs"
// @Success 200 {array} contracts.Form "Forms retrieved successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/bulk [post]
func (h *Handler) BulkGet(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	params := middlewares.QueryParams[contracts.BulkGetParams](r)
	var payload contracts.BulkGetRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	forms, err := h.queries.BulkGet(r.Context(), payload.IDs, params)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, forms)
}

// CreateStep godoc
// @Summary Create a step
// @Description Creates a step on a form.
// @Tags steps
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body contracts.CreateStepRequest true "Form title"
// @Param form_id path string true "Form ID"
// @Success 201 {object} contracts.Step "Form created successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/steps [post]
func (h *Handler) CreateStep(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload contracts.CreateStepRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	form, err := h.commands.CreateStep(r.Context(), formID, payload)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, form, http.StatusCreated)
}
