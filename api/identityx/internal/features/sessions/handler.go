package sessions

import (
	"IdentityX/internal/shared/authz"
	"net/http"

	_ "IdentityX/models"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	commands CommandService
	queries  QueryService
}

func NewHandler(
	commands CommandService,
	queries QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handler,
	jwt func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwt)
		r.Get("/sessions", h.List)
		r.Get("/sessions/me", h.Me)
		r.Delete("/sessions/{session_id}", h.RevokeByID)
		r.Delete("/sessions/others", h.RevokeOthers)
		r.Delete("/sessions", h.RevokeAll)
	})
}

// List godoc
// @Summary Lists all active user sessions
// @Description Retrieves a list of all active sessions for the authenticated user.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} models.Session "List of active user sessions"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /sessions [get]
func (handler *Handler) List(w http.ResponseWriter, r *http.Request) {
	sessions, err := handler.queries.List(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, sessions)
}

// RevokeByID godoc
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
func (handler *Handler) RevokeByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	sessionID, err := req.Path("session_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = handler.commands.RevokeByID(r.Context(), sessionID)
	if fun.Bail(w, err) {
		return
	}
	fun.OK("revoked session").Send(w)
}

// RevokeOthers godoc
// @Summary Revokes all user sessions except the current one
// @Description Invalidates all active sessions for the authenticated user, except for the one currently in use.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} object "Other sessions revoked successfully"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /sessions/others [delete]
func (handler *Handler) RevokeOthers(w http.ResponseWriter, r *http.Request) {
	err := handler.commands.RevokeOthers(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.OK("revoked sessions").Send(w)
}

// RevokeAll godoc
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
func (handler *Handler) RevokeAll(w http.ResponseWriter, r *http.Request) {
	err := handler.commands.RevokeAll(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.OK("revoked sessions").Send(w)
}

// Me godoc
// @Summary Sends current session information to user
// @Description Returns details about the current access and refresh tokens, including their expiry times.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} authz.Principal "Current session information"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /sessions/me [get]
func (handler *Handler) Me(w http.ResponseWriter, r *http.Request) {
	principal, err := authz.RequirePrincipal(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, principal)
}
