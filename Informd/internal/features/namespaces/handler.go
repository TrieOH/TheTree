package namespaces

import (
	"Informd/internal/shared/contracts"
	"net/http"

	_ "Informd/internal/shared/contracts"

	"github.com/MintzyG/FastUtilitiesNet"
	"github.com/MintzyG/FastUtilitiesNet/bind"
	"github.com/MintzyG/FastUtilitiesNet/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewHandler(
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
		r.Post("/namespaces", h.Create)
		r.With(middlewares.WithParams[contracts.BulkGetParams]()).Post("/namespaces/bulk", h.BulkGet)
	})
}

type CreateNamespaceRequest struct {
	Name string `json:"name"`
}

// Create godoc
// @Summary Create a namespace
// @Description Creates a new namespace for the authenticated user
// @Tags namespaces
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body CreateNamespaceRequest true "Project details"
// @Success 201 {object} contracts.Namespace "Namespace created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /namespaces [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)

	var payload CreateNamespaceRequest
	if err := bind.Body(req).Bind(&payload); err != nil {
		fun.Error(err).Send(w)
		return
	}

	namespace, err := h.commands.Create(r.Context(), payload.Name)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(namespace).Send(w)
}

type BulkGetRequest struct {
	IDs []uuid.UUID `json:"ids" validate:"required"`
}

// BulkGet godoc
// @Summary Bulk get namespaces
// @Description Returns a list of namespaces by their IDs. IDs should be obtained via a SpiceDB lookup on the client side.
// @Tags forms
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body BulkGetRequest true "Namespace IDs"
// @Success 200 {array} contracts.Form "Namespaces retrieved successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /namespaces/bulk [post]
func (h *Handler) BulkGet(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)

	var payload BulkGetRequest
	if err := bind.Body(req).Bind(&payload); err != nil {
		fun.Error(err).Send(w)
		return
	}

	forms, err := h.queries.BulkGet(r.Context(), payload.IDs)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(forms).Send(w)
}
