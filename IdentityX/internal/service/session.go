package service

import (
	"log"
	"time"
	"context"
	"net/http"
	"GoAuth/internal/repository"
	"GoAuth/internal/models"
	"github.com/google/uuid"
	resp "github.com/MintzyG/GoResponse/response"
)

func (s *AuthService) ListUserSessions(r *http.Request, ctx context.Context) ([]repository.UserSession, *resp.Response) {
	access_claims, err := models.GetAccessClaims(r)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	sessions, err := s.queries.ListUserSessions(ctx, access_claims.Sub.ID)
	if err != nil {
		return nil, resp.InternalServerError("error listing user sessions").WithTracePrefix("database-error").AddTrace(err)
	}
	return sessions, nil
}

func (s *AuthService) RevokeUserSessionByID(r *http.Request, ctx context.Context, session_id string) *resp.Response {
	access_claims, err := models.GetAccessClaims(r)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

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

	revoked_session, err := s.queries.RevokeUserSessionById(ctx, repository.RevokeUserSessionByIdParams{
		SessionID: sid,
		TokenID: jti,
		UserID: access_claims.Sub.ID,
	})

	if err != nil {
		return resp.InternalServerError("error revoking user session").WithTracePrefix("database-error").AddTrace(err)
	}

	err = s.queries.BlacklistToken(ctx, repository.BlacklistTokenParams{
		TokenID: revoked_session.TokenID,
		ExpiresAt: revoked_session.ExpiresAt,
	})

	if err != nil {
		return resp.InternalServerError().AddTrace("failed to blacklist user refresh token", err.Error())
	}

	return nil
}

func (s *AuthService) RevokeOtherSessions(r *http.Request, ctx context.Context) *resp.Response {
	access_claims, err := models.GetAccessClaims(r)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	refresh_claims, err := models.GetRefreshClaims(r)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	jti, err := uuid.Parse(refresh_claims.ID)
	if err != nil {
		return resp.InternalServerError().AddTrace("failed to parse refresh jti", err.Error())
	}

	revoked_sessions, err := s.queries.RevokeOtherSessions(ctx, repository.RevokeOtherSessionsParams{
		TokenID: jti,
		UserID: access_claims.Sub.ID,
	})

	if err != nil {
		return resp.InternalServerError("error revoking user sessions").WithTracePrefix("database-error").AddTrace(err)
	}

	tokenIDs := make([]uuid.UUID, len(revoked_sessions))
	expiresAt := make([]time.Time, len(revoked_sessions))

        for i, session := range revoked_sessions {
		tokenIDs[i] = session.TokenID
		expiresAt[i] = session.ExpiresAt
	}

	blacklisted_tokens, err := s.queries.BlacklistManyTokens(ctx, repository.BlacklistManyTokensParams{
		Column1: tokenIDs,
		Column2: expiresAt,
	})

	if len(blacklisted_tokens) != len(tokenIDs) {
		log.Println(blacklisted_tokens)
		log.Println(tokenIDs)
	}

	if err != nil {
		return resp.InternalServerError().AddTrace("failed to blacklist other user tokens", err.Error())
	}

	return nil
}
