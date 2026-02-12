package handlers

import (
	"net/http"
	"univents/internal/adapters/http/dto"
	"univents/internal/domain/authz"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type UniventsHandler struct{}

func NewUniventsHandler() *UniventsHandler {
	return &UniventsHandler{}
}

// Health godoc
// @Summary Health check
// @Description Returns service health status
// @Tags system
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /health [get]
func (handler *UniventsHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := dto.HealthResponse{
		Status:  "ok",
		Service: "univents-api",
	}

	resp.OK("ok").WithData(response).Send(w)
}

// ProtectedHealth godoc
// @Summary Health check
// @Description Returns service health status and user id
// @Tags system
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /protected/health [get]
func (handler *UniventsHandler) ProtectedHealth(w http.ResponseWriter, r *http.Request) {
	sub, err := authz.RequireSubject(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	response := dto.HealthResponse{
		Status:  "ok",
		Service: "univents-api",
		UserID:  sub.ID,
	}

	resp.OK("ok").WithData(response).Send(w)
}
