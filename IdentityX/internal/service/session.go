package service

import (
	"GoAuth/internal/models"
	"GoAuth/internal/repository"
	"context"
	"log"
	"net/http"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
)

func (s *AuthService) ListUserSessions(r *http.Request, ctx context.Context) ([]repository.UserSession, *resp.Response) {
	accessClaims, err := models.GetAccessClaims(r)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	sessions, err := s.queries.ListUserSessions(ctx, accessClaims.Sub.ID)
	if err != nil {
		return nil, resp.InternalServerError("error listing user sessions").WithTracePrefix("database-error").AddTrace(err)
	}
	return sessions, nil
}

func (s *AuthService) RevokeUserSessionByID(r *http.Request, ctx context.Context, sessionId string) *resp.Response {
	accessClaims, err := models.GetAccessClaims(r)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	refreshClaims, err := models.GetRefreshClaims(r)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	jti, err := uuid.Parse(refreshClaims.ID)
	if err != nil {
		return resp.InternalServerError().AddTrace("failed to parse refresh jti", err.Error())
	}

	sid, err := uuid.Parse(sessionId)
	if err != nil {
		return resp.InternalServerError().AddTrace("failed to parse session id", err.Error())
	}

	if accessClaims.Sub.SessionID == sid {
		return resp.BadRequest("can't revoke a currently active session, please logout instead")
	}

	revokedSession, err := s.queries.RevokeUserSessionById(ctx, repository.RevokeUserSessionByIdParams{
		SessionID: sid,
		TokenID:   jti,
		UserID:    accessClaims.Sub.ID,
	})

	if err != nil {
		return resp.InternalServerError("error revoking user session").WithTracePrefix("database-error").AddTrace(err)
	}

	err = s.queries.BlacklistToken(ctx, repository.BlacklistTokenParams{
		TokenID:   revokedSession.TokenID,
		ExpiresAt: revokedSession.ExpiresAt,
	})

	if err != nil {
		return resp.InternalServerError().AddTrace("failed to blacklist user refresh token", err.Error())
	}

	return nil
}

func (s *AuthService) RevokeOtherSessions(r *http.Request, ctx context.Context) *resp.Response {
	accessClaims, err := models.GetAccessClaims(r)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	refreshClaims, err := models.GetRefreshClaims(r)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	jti, err := uuid.Parse(refreshClaims.ID)
	if err != nil {
		return resp.InternalServerError().AddTrace("failed to parse refresh jti", err.Error())
	}

	revokedSessions, err := s.queries.RevokeOtherSessions(ctx, repository.RevokeOtherSessionsParams{
		TokenID: jti,
		UserID:  accessClaims.Sub.ID,
	})

	if err != nil {
		return resp.InternalServerError("error revoking user sessions").WithTracePrefix("database-error").AddTrace(err)
	}

	tokenIDs := make([]uuid.UUID, len(revokedSessions))
	expiresAt := make([]time.Time, len(revokedSessions))

	for i, session := range revokedSessions {
		tokenIDs[i] = session.TokenID
		expiresAt[i] = session.ExpiresAt
	}

	blacklistedTokens, err := s.queries.BlacklistManyTokens(ctx, repository.BlacklistManyTokensParams{
		Column1: tokenIDs,
		Column2: expiresAt,
	})

	if len(blacklistedTokens) != len(tokenIDs) {
		log.Println(blacklistedTokens)
		log.Println(tokenIDs)
	}

	if err != nil {
		return resp.InternalServerError().AddTrace("failed to blacklist other user tokens", err.Error())
	}

	return nil
}

func (s *AuthService) RevokeAllSessions(r *http.Request, ctx context.Context) *resp.Response {
	accessClaims, err := models.GetAccessClaims(r)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	revokedSessions, err := s.queries.RevokeAllSessions(ctx, accessClaims.Sub.ID)

	if err != nil {
		return resp.InternalServerError("error revoking user sessions").WithTracePrefix("database-error").AddTrace(err)
	}

	tokenIDs := make([]uuid.UUID, len(revokedSessions))
	expiresAt := make([]time.Time, len(revokedSessions))

	for i, session := range revokedSessions {
		tokenIDs[i] = session.TokenID
		expiresAt[i] = session.ExpiresAt
	}

	blacklistedTokens, err := s.queries.BlacklistManyTokens(ctx, repository.BlacklistManyTokensParams{
		Column1: tokenIDs,
		Column2: expiresAt,
	})

	if len(blacklistedTokens) != len(tokenIDs) {
		log.Println(blacklistedTokens)
		log.Println(tokenIDs)
	}

	if err != nil {
		return resp.InternalServerError().AddTrace("failed to blacklist user tokens", err.Error())
	}

	return nil
}
