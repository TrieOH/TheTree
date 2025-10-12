package handler

import (
	"net/http"

	"GoAuth/internal/models"
	"GoAuth/internal/validation"

	resp "github.com/MintzyG/GoResponse/response"
)

// Register godoc
// @Summary Register a new user
// @Description registers a client into the system
// @Tags auth
// @Accept json
// @Produce json
// @Param registerInfo body models.RegisterUserRequest true "Register request data"
// @Success 201 {string} string "Registered user"
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterUserRequest
	if rs := validation.ValidateWith(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	if rs := h.AuthService.Register(r.Context(), req); rs != nil {
		rs.Send(w)
		return
	}

	resp.Created("Registered user").Send(w)
}
