package service

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/models"
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (s *AuthService) ListUserSessions(ctx context.Context) ([]models.Session, error) {
	ctx, span := GoAuthServiceTracer.Start(ctx, "AuthService.ListUserSessions")
	defer span.End()

	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	annotateAccessClaims(span, accessClaims)

	sessions, err := s.sessionRepo.List(ctx, accessClaims.Sub.ID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("sessions.count", len(sessions)))

	return sessions, nil
}

func (s *AuthService) RevokeUserSessionByID(ctx context.Context, sessionId string) error {
	ctx, span := GoAuthServiceTracer.Start(ctx, "AuthService.RevokeUserSessionByID")
	defer span.End()

	sid, err := uuid.Parse(sessionId)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid session id").WithID(apierr.SessionInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return err
	}

	annotateAccessClaims(span, accessClaims)

	if accessClaims.Sub.SessionID == sid {
		apiErr := apierr.ErrForbidden.WithMsg("cannot revoke the currently active session").WithID(apierr.SessionSelfRevokeForbidden)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	revokedSessions, err := s.sessionRepo.DeleteByFilter(ctx, models.SessionFilter{
		UserID:    accessClaims.Sub.ID,
		SessionID: &sid,
	})
	if err != nil {
		return err
	}

	if len(revokedSessions) == 0 {
		return apierr.ErrNotFound.WithMsg("session not found").WithID(apierr.SessionNotFound)
	}

	session := revokedSessions[0]

	if err := s.revokedRefreshTokensRepo.Revoke(ctx, models.RevokedRefreshToken{
		TokenID:   session.TokenID,
		ExpiresAt: session.ExpiresAt,
	}); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RevokeOtherSessions(ctx context.Context) error {
	ctx, span := GoAuthServiceTracer.Start(ctx, "AuthService.RevokeOtherSessions")
	defer span.End()

	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return err
	}

	annotateAccessClaims(span, accessClaims)

	revokedSessions, err := s.sessionRepo.DeleteByFilter(ctx, models.SessionFilter{
		UserID:    accessClaims.Sub.ID,
		ExcludeID: &accessClaims.Sub.SessionID,
	})
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.Int("sessions.deleted.count", len(revokedSessions)))

	tokenIDs := make([]uuid.UUID, len(revokedSessions))
	expiresAt := make([]time.Time, len(revokedSessions))

	for i, session := range revokedSessions {
		tokenIDs[i] = session.TokenID
		expiresAt[i] = session.ExpiresAt
	}

	err = s.revokedRefreshTokensRepo.RevokeMany(ctx, tokenIDs, expiresAt)
	if err != nil {
		span.SetAttributes(attribute.Bool("sessions.revoked.success", false))
		return err
	}

	span.SetAttributes(attribute.Bool("sessions.revoked.success", true))
	return nil
}

func (s *AuthService) RevokeAllSessions(ctx context.Context) error {
	ctx, span := GoAuthServiceTracer.Start(ctx, "AuthService.RevokeAllSessions")
	defer span.End()

	accessClaims, err := models.GetAccessClaims(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return err
	}

	annotateAccessClaims(span, accessClaims)

	revokedSessions, err := s.sessionRepo.DeleteByFilter(ctx, models.SessionFilter{
		UserID: accessClaims.Sub.ID,
	})
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.Int("sessions.deleted.count", len(revokedSessions)))

	tokenIDs := make([]uuid.UUID, len(revokedSessions))
	expiresAt := make([]time.Time, len(revokedSessions))

	for i, session := range revokedSessions {
		tokenIDs[i] = session.TokenID
		expiresAt[i] = session.ExpiresAt
	}

	err = s.revokedRefreshTokensRepo.RevokeMany(ctx, tokenIDs, expiresAt)
	if err != nil {
		span.SetAttributes(attribute.Bool("sessions.revoked.success", false))
		return err
	}

	span.SetAttributes(attribute.Bool("sessions.revoked.success", true))
	return nil
}
