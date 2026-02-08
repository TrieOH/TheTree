package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/apierr"
	"GoAuth/internal/ports/inbounds"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
)

type SessionHandler struct {
	sessions inbounds.SessionService
}

func NewSessionHandler(uc inbounds.SessionService) *SessionHandler {
	return &SessionHandler{sessions: uc}
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
func (handler *SessionHandler) ListUserSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := handler.sessions.List(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
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
func (handler *SessionHandler) RevokeUserSessionByID(w http.ResponseWriter, r *http.Request) {
	sessionID, rs := getUUID(r, "session_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	accessToken, err := r.Cookie("access_token")
	if err != nil {
		resp.FromError(fail.New(apierr.AuthMissingAccessCookie).Trace(err.Error())).Send(w)
		return
	}

	err = handler.sessions.RevokeByID(r.Context(), sessionID, accessToken.Value)
	if err != nil {
		resp.FromError(err).Send(w)
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
func (handler *SessionHandler) RevokeOtherSessions(w http.ResponseWriter, r *http.Request) {
	accessToken, err := r.Cookie("access_token")
	if err != nil {
		resp.FromError(fail.New(apierr.AuthMissingAccessCookie).Trace(err.Error())).Send(w)
		return
	}

	err = handler.sessions.RevokeOthers(r.Context(), accessToken.Value)
	if err != nil {
		resp.FromError(err).Send(w)
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
func (handler *SessionHandler) RevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	err := handler.sessions.RevokeAll(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
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
// @Success 200 {object} dto.MeResponse "Current session information"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /sessions/me [get]
func (handler *SessionHandler) Me(w http.ResponseWriter, r *http.Request) {
	accessToken, err := r.Cookie("access_token")
	if err != nil {
		resp.FromError(fail.New(apierr.AuthMissingAccessCookie).Trace(err.Error())).Send(w)
		return
	}

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		resp.FromError(fail.New(apierr.AuthMissingRefreshCookie).Trace(err.Error())).Send(w)
		return
	}

	me, err := handler.sessions.Me(r.Context(), accessToken.Value, refreshToken.Value)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.MeOutputToMeResponse(*me)).Send(w)
}
