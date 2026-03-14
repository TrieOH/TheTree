package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"encoding/json"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
)

type SessionHandler struct {
	sessions inbounds.SessionService
	redis    outbounds.RedisCacheService
}

func NewSessionHandler(uc inbounds.SessionService, redis outbounds.RedisCacheService) *SessionHandler {
	return &SessionHandler{sessions: uc, redis: redis}
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
	ctx := r.Context()
	sessionID, rs := getUUID(r, "session_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	svcCookie, err := r.Cookie("svc_session")
	if err != nil || svcCookie.Value == "" {
		resp.Unauthorized().WithMsg("missing svc_session cookie").Send(w)
		return
	}

	key := "svc_session:" + svcCookie.Value
	data, found, err := handler.redis.GetAny(ctx, key)
	if err != nil || !found {
		resp.Unauthorized().WithMsg("invalid service session").Send(w)
		return
	}

	bytesData, ok := data.([]byte)
	if !ok {
		resp.Unauthorized().WithMsg("invalid session type").Send(w)
		return
	}

	// Inline unmarshal
	var snapshot dto.MeResponse
	if err := json.Unmarshal(bytesData, &snapshot); err != nil {
		_ = handler.redis.Delete(ctx, key)
		resp.Unauthorized().WithMsg("failed to unmarshal session").Send(w)
		return
	}

	err = handler.sessions.RevokeByID(ctx, sessionID, snapshot.AccessClaims.Sub.SessionID)
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
		resp.FromError(fail.New(errx.AuthMissingAccessCookie).Trace(err.Error())).Send(w)
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
	ctx := r.Context()

	svcCookie, err := r.Cookie("svc_session")
	if err != nil || svcCookie.Value == "" {
		resp.Unauthorized().WithMsg("missing svc_session cookie").Send(w)
		return
	}

	key := "svc_session:" + svcCookie.Value
	data, found, err := handler.redis.GetAny(ctx, key)
	if err != nil || !found {
		resp.Unauthorized().WithMsg("invalid service session").Send(w)
		return
	}

	bytesData, ok := data.([]byte)
	if !ok {
		resp.Unauthorized().WithMsg("invalid session type").Send(w)
		return
	}

	var snapshot authz.ServiceSnapshot
	if err := json.Unmarshal(bytesData, &snapshot); err != nil {
		_ = handler.redis.Delete(ctx, key)
		resp.Unauthorized().WithMsg("failed to unmarshal session").Send(w)
		return
	}

	resp.OK().WithData(snapshot).Send(w)
}
