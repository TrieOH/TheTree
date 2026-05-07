package keys

import (
	"net/http"
	"time"

	_ "Informd/internal/shared/contracts"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
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
		r.Post("/api-keys", h.Create)
		r.Post("/api-keys/bulk", h.BulkGet)
		r.Delete("/api-keys/{id}", h.Revoke)
	})
}

type CreateAPIKeyRequest struct {
	Name string `json:"name" validate:"required"`
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
	if bind.BailInto(w, req, &payload) {
		return
	}
	rawKey, apiKey, err := h.commands.Create(r.Context(), payload.Name)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, CreateAPIKeyResponse{
		APIKeyResponse: APIKeyResponse{
			ID:        apiKey.ID,
			Name:      apiKey.Name,
			Prefix:    apiKey.KeyPrefix,
			CreatedAt: apiKey.CreatedAt,
			RevokedAt: apiKey.RevokedAt,
		},
		Key: rawKey,
	}, http.StatusCreated)
}

type BulkGetRequest struct {
	IDs []uuid.UUID `json:"ids" validate:"required"`
}

// BulkGet godoc
// @Summary Bulk get api keys
// @Description Returns a list of api keys by their IDs. IDs should be obtained via a SpiceDB lookup on the client side.
// @Tags api_keys
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body BulkGetRequest true "APIKey IDs"
// @Success 200 {array} contracts.Form "Forms retrieved successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /api_keys/bulk [post]
func (h *Handler) BulkGet(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	var payload BulkGetRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	forms, err := h.queries.BulkGet(r.Context(), payload.IDs)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, forms)
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
	if fun.Bail(w, err) {
		return
	}
	err = h.commands.RevokeAPIKey(r.Context(), keyID)
	if fun.Bail(w, err) {
		return
	}
	fun.OK("key revoked").Send(w)
}
