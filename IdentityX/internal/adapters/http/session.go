package http

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/application/session"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type SessionHandler struct {
	uc *session.UseCase
}

func NewSessionHandler(uc *session.UseCase) *SessionHandler {
	return &SessionHandler{uc: uc}
}

// ListUserSessions godoc
// @Summary Lists all active user sessions
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} dto.SessionResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
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
// @Summary Revokes a user session if it isn't the current one
// @Tags auth
// @Accept json
// @Produce json
// @Param session_id path string true "ID of the session to be invalidated"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {string} string "revoked session"
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /sessions/{session_id} [delete]
func (sh *SessionHandler) RevokeUserSessionByID(w http.ResponseWriter, r *http.Request) {
	err := sh.uc.RevokeUserSessionByID(r.Context(), r.PathValue("session_id"))
	if err != nil {
		ErrToResp(err).Send(w)
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
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
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
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {string} string "revoked sessions"
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
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
// @Summary Sends session info to user
// @Description sends the whole access info and the refresh expiry time
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Header 200 {string} Set-Cookie "access_token cookie for authentication"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Success 200 {object} map[string]any
// @Failure 500 {object} domain.ErrorResponse
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
