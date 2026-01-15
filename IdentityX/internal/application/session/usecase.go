package session

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbound"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	usecaseTracer = otel.Tracer("auth_usecase")
)

type UseCase struct {
	sessions outbound.SessionRepository
}

var _ inbounds.SessionService = (*UseCase)(nil)

func New(
	sessions outbound.SessionRepository,
) inbounds.SessionService {
	return &UseCase{
		sessions: sessions,
	}
}

// List handles the business logic for listing all sessions for the authenticated user.
func (uc *UseCase) List(ctx context.Context) ([]inbounds.OutputSession, error) {
	ctx, span := usecaseTracer.Start(ctx, "SessionService.List")
	defer span.End()

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	sessions, err := uc.sessions.List(ctx, principal.UserID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("sessions.count", len(sessions)))

	return inbounds.OutputSessionSliceFromSessionSlice(sessions), nil
}

// RevokeByID handles the business logic for revoking a specific session for the authenticated user.
// It ensures that the user is not revoking the current session.
func (uc *UseCase) RevokeByID(ctx context.Context, sessionID uuid.UUID) error {
	ctx, span := usecaseTracer.Start(ctx, "SessionService.RevokeByID")
	defer span.End()

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return apierr.FromService(span, err)
	}

	if principal.SessionID == sessionID {
		return apierr.FromService(span, inbounds.ErrRevokeCurrentSession{})
	}

	sess, err := uc.sessions.MarkRevokedByID(ctx, principal.UserID, sessionID)
	if apierr.IsNotFound(err) {
		return apierr.FromService(span, inbounds.ErrSessionNotFound{})
	} else if err != nil {
		return err
	}

	span.SetAttributes(
		attribute.String("session.id", sess.SessionID.String()),
	)

	if sess.RevokedAt != nil {
		span.SetAttributes(attribute.String("session.revoked_at", sess.RevokedAt.String()))
	}

	return nil
}

// RevokeOthers handles the business logic for revoking all sessions for the authenticated user except for the current one.
func (uc *UseCase) RevokeOthers(ctx context.Context) error {
	ctx, span := usecaseTracer.Start(ctx, "SessionService.RevokeOthers")
	defer span.End()

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return apierr.FromService(span, err)
	}

	revokedCount, err := uc.sessions.MarkRevokedByFilter(ctx, session.Filter{
		UserID:    principal.UserID,
		ExcludeID: &principal.SessionID,
	})
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.Int("sessions.revoked.count", revokedCount))
	return nil
}

// RevokeAll handles the business logic for revoking all sessions for the authenticated user.
func (uc *UseCase) RevokeAll(ctx context.Context) error {
	ctx, span := usecaseTracer.Start(ctx, "SessionService.RevokeAll")
	defer span.End()

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return apierr.FromService(span, err)
	}

	revokedCount, err := uc.sessions.MarkRevokedByFilter(ctx, session.Filter{
		UserID: principal.UserID,
	})
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.Int("sessions.revoked.count", revokedCount))

	return nil
}

// Me returns the principal of the authenticated user.
func (uc *UseCase) Me(ctx context.Context) (*inbounds.PrincipalOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SessionService.Me")
	defer span.End()
	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}
	return inbounds.PrincipalToPrincipalOutput(*principal), nil
}
