package session

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	usecaseTracer = otel.Tracer("auth_usecase")
)

type UseCase struct {
	sessions outbounds.SessionRepository
	tx       inbounds.TxRunner
}

var _ inbounds.SessionService = (*UseCase)(nil)

func New(
	sessions outbounds.SessionRepository,
	tx inbounds.TxRunner,
) inbounds.SessionService {
	return &UseCase{
		sessions: sessions,
		tx:       tx,
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

	var identityType session.IdentityType
	if principal.ProjectID == nil {
		identityType = session.ClientIdentity
	} else {
		identityType = session.ProjectIdentity
	}

	sessions, err := uc.sessions.List(ctx, principal.UserID, identityType)
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
		return fail.New(apierr.SessionSelfRevokeForbidden)
	}

	var identityType session.IdentityType
	if principal.ProjectID == nil {
		identityType = session.ClientIdentity
	} else {
		identityType = session.ProjectIdentity
	}

	var sess *session.Session
	sess, err = uc.sessions.MarkRevokedByID(ctx, principal.UserID, sessionID, identityType)
	if fail.Is(err, apierr.SQLNotFound) {
		return fail.New(apierr.SessionNotFound)
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

	var identityType session.IdentityType
	if principal.ProjectID == nil {
		identityType = session.ClientIdentity
	} else {
		identityType = session.ProjectIdentity
	}

	revokedCount, err := uc.sessions.MarkRevokedByFilter(ctx, session.Filter{
		IdentityType: identityType,
		EntityID:     principal.UserID,
		ExcludeID:    &principal.SessionID,
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

	var identityType session.IdentityType
	if principal.ProjectID == nil {
		identityType = session.ClientIdentity
	} else {
		identityType = session.ProjectIdentity
	}

	revokedCount, err := uc.sessions.MarkRevokedByFilter(ctx, session.Filter{
		IdentityType: identityType,
		EntityID:     principal.UserID,
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
