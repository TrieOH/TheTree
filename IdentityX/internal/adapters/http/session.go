package http

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/application/session"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

type SessionHandler struct {
	uc *session.UseCase
}

func NewSessionHandler(uc *session.UseCase) *SessionHandler {
	return &SessionHandler{uc: uc}
}

// ListUserSessions godoc
// @Summary Lists all active user sessions
// @Description Retrieves a list of all active sessions for the authenticated user.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} dto.SessionResponse "List of active user sessions"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /sessions [get]
func (sh *SessionHandler) ListUserSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := sh.uc.ListUserSessions(r.Context())
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK().
		WithData(dto.SessionResponseSliceFromSessionOutputSlice(sessions)).
		Send(w)
}

// RevokeUserSessionByID godoc
// @Summary Revokes a user session by ID
// @Description Revokes a specific user session by its ID, provided it's not the current session.
// @Tags auth
// @Accept json
// @Produce json
// @Param session_id path string true "ID of the session to be invalidated"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} object "Session revoked successfully"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid session ID or trying to revoke current session"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 404 {object} ErrorResponse "Not Found: Session not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /sessions/{session_id} [delete]
func (sh *SessionHandler) RevokeUserSessionByID(w http.ResponseWriter, r *http.Request) {
	err := sh.uc.RevokeUserSessionByID(r.Context(), chi.URLParam(r, "session_id"))
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK("revoked session").Send(w)
}

// RevokeOtherSessions godoc
// @Summary Revokes all user sessions except the current one
// @Description Invalidates all active sessions for the authenticated user, except for the one currently in use.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} object "Other sessions revoked successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /sessions/others [delete]
func (sh *SessionHandler) RevokeOtherSessions(w http.ResponseWriter, r *http.Request) {
	err := sh.uc.RevokeOtherSessions(r.Context())
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK("revoked sessions").Send(w)
}

// RevokeAllSessions godoc
// @Summary Revokes all user sessions
// @Description Invalidates all active sessions for the authenticated user, including the current one.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} object "All sessions revoked successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /sessions [delete]
func (sh *SessionHandler) RevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	err := sh.uc.RevokeAllSessions(r.Context())
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK("revoked sessions").Send(w)
}

// Me godoc
// @Summary Sends current session information to user
// @Description Returns details about the current access and refresh tokens, including their expiry times.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} object{access=object,refresh_expire_date=string} "Current session information"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /sessions/me [get]
func (sh *SessionHandler) Me(w http.ResponseWriter, r *http.Request) {
	principal, err := sh.uc.Me(r.Context())
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK().WithData(map[string]interface{}{
		"access":              principal.AccessClaims,
		"refresh_expire_date": principal.RefreshClaims.ExpiresAt,
	}).Send(w)
}
