package system

import (
	"TrieForms/internal/shared/authz"
	"encoding/json"
	"net/http"
	"strings"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
)

type SystemHandler struct {
	gaClient *goauth.Client
}

func NewSystemHandler(gaClient *goauth.Client) *SystemHandler {
	return &SystemHandler{
		gaClient: gaClient,
	}
}

type HealthResponse struct {
	Status  string    `json:"status" example:"ok"`
	Service string    `json:"service" example:"system-api"`
	UserID  uuid.UUID `json:"user_id,omitempty" example:"some-uuid"`
}

// Health godoc
// @Summary Health check
// @Description Returns service health status
// @Tags system
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (handler *SystemHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "ok",
		Service: "forms-api",
	}

	resp.OK("ok").WithData(response).Send(w)
}

// ProtectedHealth godoc
// @Summary Protected health check
// @Description Returns service health status and authenticated user id
// @Tags system
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Success 200 {object} HealthResponse
// @Failure 401 {object} swag.ErrorResponse
// @Router /protected/health [get]
func (handler *SystemHandler) ProtectedHealth(w http.ResponseWriter, r *http.Request) {
	sub, err := authz.RequireSubject(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	response := HealthResponse{
		Status:  "ok",
		Service: "univents-api",
		UserID:  sub.ID,
	}

	resp.OK("ok").WithData(response).Send(w)
}

// Exchange godoc
// @Summary      Exchange global access token for service session
// @Description  Validates a global access token and creates a service session payload using a snapshot builder. Returns the session ID, TTL, and claims. Frontend sets the session cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer token, e.g., 'Bearer <token>'"
// @Success      200  {object}  goauth.SessionResult "Service session created"
// @Failure      400  {object}  object  "Bad request / empty token"
// @Failure      401  {object}  object  "Unauthorized / invalid token"
// @Failure      500  {object}  object  "Internal server error"
// @Router       /auth/exchange [post]
func (handler *SystemHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		resp.Unauthorized().WithMsg("missing bearer token").Send(w)
		return
	}

	accessToken := strings.TrimSpace(authHeader[7:])
	if accessToken == "" {
		resp.Unauthorized().WithMsg("empty bearer token").Send(w)
		return
	}

	sessionRes, err := handler.gaClient.ExchangeAndCreateSession(ctx, accessToken, BuildServiceSnapshot)
	if err != nil {
		resp.Unauthorized("failed to create service session: " + err.Error()).Send(w)
		return
	}

	resp.OK("service session created").WithData(sessionRes).Send(w)
}

// SnapshotPayload is the typed payload stored in the service session
type SnapshotPayload struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

// BuildServiceSnapshot builds a typed snapshot for the service session
func BuildServiceSnapshot(claims *goauth.AccessClaims) ([]byte, error) {
	payload := SnapshotPayload{
		UserID: claims.Sub.ID,
		Email:  claims.Sub.Email,
	}
	return json.Marshal(payload)
}

// UnmarshalSnapshot unmarshals the session bytes into a typed payload
func UnmarshalSnapshot(data []byte) (*SnapshotPayload, error) {
	if data == nil {
		return nil, fail.New(goauth.SDKUnknownErrorID).WithArgs("session data is nil")
	}

	var payload SnapshotPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fail.New(goauth.SDKUnknownErrorID).WithArgs("failed to unmarshal session: " + err.Error())
	}

	return &payload, nil
}
