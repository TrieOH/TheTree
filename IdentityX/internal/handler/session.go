package handler

import (
	"net/http"
	resp "github.com/MintzyG/GoResponse/response"
)

// ListUserSessions godoc
// @Summary Lists all active user sessions
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} repository.UserSession
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /sessions [get]
func (h *AuthHandler) ListUserSessions(w http.ResponseWriter, r *http.Request) {
	sessions, rs := h.AuthService.ListUserSessions(r.Context())
	if rs != nil {
		rs.Send(w)
		return
	}

	resp.OK().WithData(sessions).Send(w)
}

// RevokeUserSession godoc
// @Summary Revokes a user session if it isn't the current one
// @Tags auth
// @Accept json
// @Produce json
// @Param session_id path string true "ID of the session to be invalidated"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {string} string "revoked session"
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /sessions/{session_id} [delete]
func (h *AuthHandler) RevokeUserSession(w http.ResponseWriter, r *http.Request) {
	rs := h.AuthService.RevokeUserSession(r, r.Context(), r.PathValue("session_id"))
	if rs != nil {
		rs.Send(w)
		return
	}

	resp.OK("revoked session").Send(w)
}

// RevokeOtherSessions godoc
// @Summary Revokes all user sessions but the current one
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {string} string "revoked sessions"
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /sessions/others [delete]
func (h *AuthHandler) RevokeOtherSessions(w http.ResponseWriter, r *http.Request) {
	rs := h.AuthService.RevokeOtherSessions(r, r.Context())
	if rs != nil {
		rs.Send(w)
		return
	}

	resp.OK("revoked sessions").Send(w)
}
