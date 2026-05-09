package system

import (
	"net/http"
	"univents/internal/shared/authz"

	_ "univents/internal/shared/contracts"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

type UniventsHandler struct{}

func NewUniventsHandler() *UniventsHandler {
	return &UniventsHandler{}
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
func (handler *UniventsHandler) Health(w http.ResponseWriter, _ *http.Request) {
	response := HealthResponse{
		Status:  "ok",
		Service: "univents-api",
	}

	fun.OK("ok").WithData(response).Send(w)
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
func (handler *UniventsHandler) ProtectedHealth(w http.ResponseWriter, r *http.Request) {
	sub, err := authz.RequireSubject(r.Context())
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	response := HealthResponse{
		Status:  "ok",
		Service: "univents-api",
		UserID:  sub.ID,
	}

	fun.OK("ok").WithData(response).Send(w)
}
