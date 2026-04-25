package namespaces

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
		r.Get("/namespaces", h.List)
		r.Post("/namespaces", h.Create)
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

// List godoc
// @Summary List namespaces
// @Description Lists all namespaces owned by the authenticated user
// @Tags namespaces
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Success 200 {array} contracts.Namespace "Namespaces retrieved successfully"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /namespaces [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	namespaces, err := h.queries.List(r.Context())
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(namespaces).Send(w)
}
