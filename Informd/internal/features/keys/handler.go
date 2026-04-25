package keys

import (
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/xslices"
	"net/http"
	"time"

	_ "Informd/internal/shared/contracts"

	"github.com/MintzyG/FastUtilitiesNet"
	"github.com/MintzyG/FastUtilitiesNet/bind"
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
		r.Get("/api-keys", h.List)
		r.Post("/api-keys", h.Create)
		r.Delete("/api-keys/{id}", h.Revoke)
	})
}

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

type APIKeyResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Prefix    string     `json:"prefix"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

type CreateAPIKeyResponse struct {
	APIKeyResponse
	Key string `json:"key"` // only returned once
}

// Create godoc
// @Summary Create an API key
// @Description Creates a new API key for the given project. The raw key is returned once and never stored.
// @Tags api_keys
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param project_id path string true "Project ID"
// @Param request body CreateAPIKeyRequest true "API key details"
// @Success 201 {object} CreateAPIKeyResponse "API key created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /api-keys [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	var payload CreateAPIKeyRequest
	if err := bind.Body(req).Bind(&payload); err != nil {
		fun.Error(err).Send(w)
		return
	}

	rawKey, apiKey, err := h.commands.Create(r.Context(), payload.Name)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(CreateAPIKeyResponse{
		APIKeyResponse: APIKeyResponse{
			ID:        apiKey.ID,
			Name:      apiKey.Name,
			Prefix:    apiKey.KeyPrefix,
			CreatedAt: apiKey.CreatedAt,
			RevokedAt: apiKey.RevokedAt,
		},
		Key: rawKey,
	}).Send(w)
}

// List godoc
// @Summary List API keys
// @Description Lists all API keys for the given project (raw keys are never returned)
// @Tags api_keys
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param project_id path string true "Project ID"
// @Success 200 {array} APIKeyResponse "API keys retrieved successfully"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /api-keys [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	keys, err := h.queries.List(r.Context())
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	out := xslices.MapSlice(keys, func(k contracts.APIKey) APIKeyResponse {
		return APIKeyResponse{
			ID:        k.ID,
			Name:      k.Name,
			Prefix:    k.KeyPrefix,
			CreatedAt: k.CreatedAt,
			RevokedAt: k.RevokedAt,
		}
	})

	fun.OK().WithData(out).Send(w)
}

// Revoke godoc
// @Summary Revoke an API key
// @Description Revokes the given API key, immediately invalidating it
// @Tags api_keys
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param project_id path string true "Project ID"
// @Param id path string true "API key ID"
// @Success 200 {object} object "Key revoked"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /api-keys/{id} [delete]
func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)

	keyID, err := req.Path("id").UUID()
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	if err = h.commands.RevokeAPIKey(r.Context(), keyID); err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("key revoked").Send(w)
}
