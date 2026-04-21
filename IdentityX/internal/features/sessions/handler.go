package sessions

import (
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"IdentityX/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	sessions CommandService
	redis    ports.RedisCacheService
}

func NewHandler(
	sessions CommandService,
	redis ports.RedisCacheService,
) *Handler {
	return &Handler{
		sessions: sessions,
		redis:    redis,
	}
}

// ListUserSessions godoc
// @Summary Lists all active user sessions
// @Description Retrieves a list of all active sessions for the authenticated user.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} contracts.Session "List of active user sessions"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /sessions [get]
func (handler *Handler) ListUserSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := handler.sessions.List(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().
		WithData(sessions).
		Send(w)
}

type MeResponse struct {
	RefreshExpireDate *jwt.NumericDate       `json:"refresh_expire_date"`
	AccessClaims      contracts.AccessClaims `json:"access"`
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
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid session ID or trying to revoke current session"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 404 {object} contracts.ErrorResponse "Not Found: Session not found"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /sessions/{session_id} [delete]
func (handler *Handler) RevokeUserSessionByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID, rs := validation.GetUUID(r, "session_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	err := handler.sessions.RevokeByID(ctx, sessionID)
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
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /sessions/others [delete]
func (handler *Handler) RevokeOtherSessions(w http.ResponseWriter, r *http.Request) {
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
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /sessions [delete]
func (handler *Handler) RevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	err := handler.sessions.RevokeAll(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("revoked sessions").Send(w)
}

// Me godoc
// @Summary Sends current session information to user
// @Description Returns details about the current access and refresh security, including their expiry times.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} authz.Principal "Current session information"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /sessions/me [get]
func (handler *Handler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(principal).Send(w)
}
