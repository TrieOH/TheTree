package service

import (
	"context"
	"net/http"
	"GoAuth/internal/repository"
	"GoAuth/internal/models"
	"github.com/google/uuid"
	resp "github.com/MintzyG/GoResponse/response"
)

func (s *AuthService) ListUserSessions(ctx context.Context) ([]repository.UserSession, *resp.Response) {
	sessions, err := s.queries.ListUserSessions(ctx)
	if err != nil {
		return nil, resp.InternalServerError("error listing user sessions").WithTracePrefix("database-error").AddTrace(err)
	}
	return sessions, nil
}

func (s *AuthService) RevokeUserSession(r *http.Request, ctx context.Context, session_id string) *resp.Response {
	refresh_claims, err := models.GetRefreshClaims(r)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	jti, err := uuid.Parse(refresh_claims.ID)
	if err != nil {
		return resp.InternalServerError().AddTrace("failed to parse refresh jti", err.Error())
	}

	sid, err := uuid.Parse(session_id)
	if err != nil {
		return resp.InternalServerError().AddTrace("failed to parse session id", err.Error())
	}

	token_id, err := s.queries.RevokeUserSession(ctx, repository.RevokeUserSessionParams{
		SessionID: sid,
		TokenID: jti,
	})

	if err != nil {
		return resp.InternalServerError("error revoking user session").WithTracePrefix("database-error").AddTrace(err)
	}

	err = s.queries.BlacklistToken(ctx, repository.BlacklistTokenParams{
		TokenID: token_id,
		AccessJti: refresh_claims.Sub.AccessJTI,
		ExpiresAt: refresh_claims.ExpiresAt.Time,
	})

	if err != nil {
		return resp.InternalServerError().AddTrace("failed to blacklist user refresh token", err.Error())
	}

	return nil
}
