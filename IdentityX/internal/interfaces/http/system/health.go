package system

import (
	"IdentityX/internal/shared/authz"
	"net/http"

	_ "IdentityX/internal/shared/contracts"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type HealthResponse struct {
	Status  string    `json:"status" example:"ok"`
	Service string    `json:"service" example:"univents-api"`
	UserID  uuid.UUID `json:"user_id,omitempty" example:"some-uuid"`
}

// Health godoc
// @Summary Health check
// @Description Returns service health status
// @Tags system
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (handler *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	response := HealthResponse{
		Status:  "ok",
		Service: "IdentityXAPI",
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
func (handler *Handler) ProtectedHealth(w http.ResponseWriter, r *http.Request) {
	sub, err := authz.RequirePrincipal(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	response := HealthResponse{
		Status:  "ok",
		Service: "IdentityXAPI",
		UserID:  sub.UserID,
	}

	resp.OK("ok").WithData(response).Send(w)
}
