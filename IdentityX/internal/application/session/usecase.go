package session

import (
	"GoAuth/internal/adapters/observability/tracing"
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/authz"
	"GoAuth/internal/domain/revoked_refreshes"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/ports/outbound"
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	usecaseTracer = otel.Tracer("auth_usecase")
)

type UseCase struct {
	sessions outbound.SessionRepository
	refresh  outbound.RevokedRefreshTokenRepository
}

func New(
	sessions outbound.SessionRepository,
	refresh outbound.RevokedRefreshTokenRepository,
) *UseCase {
	return &UseCase{
		sessions: sessions,
		refresh:  refresh,
	}
}

// ListUserSessions handles the business logic for listing all sessions for the authenticated user.
func (uc *UseCase) ListUserSessions(ctx context.Context) ([]OutputSession, error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.ListUserSessions")
	defer span.End()

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	tracing.AnnotatePrincipal(span, principal)

	sessions, err := uc.sessions.List(ctx, principal.UserID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("sessions.count", len(sessions)))

	return OutputSessionSliceFromSessionSlice(sessions), nil
}

// RevokeUserSessionByID handles the business logic for revoking a specific session for the authenticated user.
// It ensures that the user is not revoking the current session.
func (uc *UseCase) RevokeUserSessionByID(ctx context.Context, sessionId string) error {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.RevokeUserSessionByID")
	defer span.End()

	sid, err := uuid.Parse(sessionId)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid session id").WithID(apierr.SessionInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return err
	}

	tracing.AnnotatePrincipal(span, principal)

	if principal.SessionID == sid {
		apiErr := apierr.ErrForbidden.WithMsg("cannot revoke the currently active session").WithID(apierr.SessionSelfRevokeForbidden)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	revokedSessions, err := uc.sessions.DeleteByFilter(ctx, session.Filter{
		UserID:    principal.UserID,
		SessionID: &sid,
	})
	if err != nil {
		return err
	}

	if len(revokedSessions) == 0 {
		return apierr.ErrNotFound.WithMsg("session not found").WithID(apierr.SessionNotFound)
	}

	sess := revokedSessions[0]
	if err := uc.refresh.Revoke(ctx, revoked_refreshes.RevokedRefreshToken{
		TokenID:   sess.TokenID,
		ExpiresAt: sess.ExpiresAt,
	}); err != nil {
		return err
	}

	return nil
}

// RevokeOtherSessions handles the business logic for revoking all sessions for the authenticated user except for the current one.
func (uc *UseCase) RevokeOtherSessions(ctx context.Context) error {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.RevokeOtherSessions")
	defer span.End()

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return err
	}

	tracing.AnnotatePrincipal(span, principal)

	revokedSessions, err := uc.sessions.DeleteByFilter(ctx, session.Filter{
		UserID:    principal.UserID,
		ExcludeID: &principal.SessionID,
	})
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.Int("sessions.deleted.count", len(revokedSessions)))

	tokenIDs := make([]uuid.UUID, len(revokedSessions))
	expireAts := make([]time.Time, len(revokedSessions))
	for i, sess := range revokedSessions {
		tokenIDs[i] = sess.TokenID
		expireAts[i] = sess.ExpiresAt
	}

	err = uc.refresh.RevokeMany(ctx, tokenIDs, expireAts)
	if err != nil {
		span.SetAttributes(attribute.Bool("sessions.revoked.success", false))
		return err
	}

	span.SetAttributes(attribute.Bool("sessions.revoked.success", true))
	return nil
}

// RevokeAllSessions handles the business logic for revoking all sessions for the authenticated user.
func (uc *UseCase) RevokeAllSessions(ctx context.Context) error {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.RevokeAllSessions")
	defer span.End()

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return err
	}

	tracing.AnnotatePrincipal(span, principal)

	revokedSessions, err := uc.sessions.DeleteByFilter(ctx, session.Filter{
		UserID: principal.UserID,
	})
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.Int("sessions.deleted.count", len(revokedSessions)))

	tokenIDs := make([]uuid.UUID, len(revokedSessions))
	expireAts := make([]time.Time, len(revokedSessions))
	for i, sess := range revokedSessions {
		tokenIDs[i] = sess.TokenID
		expireAts[i] = sess.ExpiresAt
	}

	err = uc.refresh.RevokeMany(ctx, tokenIDs, expireAts)
	if err != nil {
		span.SetAttributes(attribute.Bool("sessions.revoked.success", false))
		return err
	}

	span.SetAttributes(attribute.Bool("sessions.revoked.success", true))
	return nil
}

// Me returns the principal of the authenticated user.
func (uc *UseCase) Me(ctx context.Context) (*authz.Principal, error) {
	return authz.RequirePrincipal(ctx)
}
