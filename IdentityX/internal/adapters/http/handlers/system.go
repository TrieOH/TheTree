package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/domain/authz"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// Health godoc
// @Summary Health check
// @Description Returns service health status
// @Tags system
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /health [get]
func (handler *SystemHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := dto.HealthResponse{
		Status:  "ok",
		Service: "univents-api",
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
// @Success 200 {object} dto.HealthResponse
// @Failure 401 {object} ErrorResponse
// @Router /protected/health [get]
func (handler *SystemHandler) ProtectedHealth(w http.ResponseWriter, r *http.Request) {
	sub, err := authz.RequirePrincipal(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	response := dto.HealthResponse{
		Status:  "ok",
		Service: "univents-api",
		UserID:  sub.UserID,
	}

	resp.OK("ok").WithData(response).Send(w)
}
