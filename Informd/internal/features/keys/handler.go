package keys

import (
	"TrieForms/internal/shared/validation"
	"net/http"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewApiKeysHandler(
	commands *CommandService,
	queries *QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
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
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /projects/{project_id}/keys [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs == nil {
		rs.Send(w)
		return
	}

	var req CreateAPIKeyRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.Error(err).Send(w)
		return
	}

	rawKey, apiKey, err := h.commands.Create(r.Context(), req.Name, projectID)
	if err != nil {
		resp.Error(err).Send(w)
		return
	}

	resp.Created().WithData(CreateAPIKeyResponse{
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
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /projects/{project_id}/keys [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs == nil {
		rs.Send(w)
		return
	}

	keys, err := h.queries.List(r.Context(), projectID)
	if err != nil {
		resp.Error(err).Send(w)
		return
	}

	out := make([]APIKeyResponse, 0, len(keys))
	for _, k := range keys {
		out = append(out, APIKeyResponse{
			ID:        k.ID,
			Name:      k.Name,
			Prefix:    k.KeyPrefix,
			CreatedAt: k.CreatedAt,
			RevokedAt: k.RevokedAt,
		})
	}

	resp.OK().WithData(out).Send(w)
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
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /projects/{project_id}/keys/{id} [delete]
func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs == nil {
		rs.Send(w)
		return
	}

	keyID, rs := validation.GetUUID(r, "id")
	if rs != nil {
		rs.Send(w)
		return
	}

	if err := h.commands.RevokeAPIKey(r.Context(), projectID, keyID); err != nil {
		resp.Error(err).Send(w)
		return
	}

	resp.OK("key revoked").Send(w)
}
