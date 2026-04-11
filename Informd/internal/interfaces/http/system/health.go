package system

import (
	"TrieForms/internal/shared/authz"
	"net/http"

	_ "TrieForms/internal/shared/contracts"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
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
func (handler *SystemHandler) Health(w http.ResponseWriter, _ *http.Request) {
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
// @Failure 401 {object} contracts.ErrorResponse
// @Router /protected/health [get]
func (handler *SystemHandler) ProtectedHealth(w http.ResponseWriter, r *http.Request) {
	sub, err := authz.RequireSubject(r.Context())
	if err != nil {
		resp.Error(err).Send(w)
		return
	}

	response := HealthResponse{
		Status:  "ok",
		Service: "univents-api",
		UserID:  sub.ID,
	}

	resp.OK("ok").WithData(response).Send(w)
}
