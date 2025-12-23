package service

import (
	"GoAuth/internal/models"
	"context"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
)

func (s *AuthService) ListUserSessions(ctx context.Context) ([]models.Session, *resp.Response) {
	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		return nil, resp.InternalServerError().AddTrace(err)
	}

	sessions, err := s.sessionRepo.List(ctx, accessClaims.Sub.ID)
	if err != nil {
		return nil, resp.InternalServerError("error listing user sessions").WithTracePrefix("database-error").AddTrace(err)
	}
	return sessions, nil
}

func (s *AuthService) RevokeUserSessionByID(ctx context.Context, sessionId string) *resp.Response {
	sid, err := uuid.Parse(sessionId)
	if err != nil {
		return resp.InternalServerError().AddTrace("failed to parse session id", err.Error())
	}

	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	if accessClaims.Sub.SessionID == sid {
		return resp.BadRequest("can't revoke a currently active session, please logout instead")
	}

	revokedSession, err := s.sessionRepo.DeleteByFilter(ctx, models.SessionFilter{
		UserID:    accessClaims.Sub.ID,
		SessionID: &sid,
	})

	if err != nil {
		return resp.InternalServerError("error revoking user session").WithTracePrefix("database-error").AddTrace(err)
	}

	err = s.revokedRefreshTokensRepo.Revoke(ctx, models.RefreshBlacklist{
		TokenID:   revokedSession[0].TokenID,
		ExpiresAt: revokedSession[0].ExpiresAt,
	})

	if err != nil {
		return resp.InternalServerError().AddTrace("failed to blacklist user refresh token", err.Error())
	}

	return nil
}

func (s *AuthService) RevokeOtherSessions(ctx context.Context) *resp.Response {
	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	revokedSessions, err := s.sessionRepo.DeleteByFilter(ctx, models.SessionFilter{
		UserID:    accessClaims.Sub.ID,
		ExcludeID: &accessClaims.Sub.SessionID,
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

	err = s.revokedRefreshTokensRepo.RevokeMany(ctx, tokenIDs, expiresAt)

	if err != nil {
		return resp.InternalServerError().AddTrace("failed to blacklist other user tokens", err.Error())
	}

	return nil
}

func (s *AuthService) RevokeAllSessions(ctx context.Context) *resp.Response {
	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		return resp.InternalServerError().AddTrace(err)
	}

	revokedSessions, err := s.sessionRepo.DeleteByFilter(ctx, models.SessionFilter{
		UserID: accessClaims.Sub.ID,
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

	err = s.revokedRefreshTokensRepo.RevokeMany(ctx, tokenIDs, expiresAt)

	if err != nil {
		return resp.InternalServerError().AddTrace("failed to blacklist user tokens", err.Error())
	}

	return nil
}
