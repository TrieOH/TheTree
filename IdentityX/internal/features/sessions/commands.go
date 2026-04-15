package sessions

import (
	"IdentityX/internal/features/tokens"
	"IdentityX/internal/platform/database"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	sessions ports.SessionRepository
	tokens   tokens.CommandService
	logger   *zap.Logger
	tracer   trace.Tracer
	txRunner database.TxRunner
}

func NewCommandService(
	sessions ports.SessionRepository,
	tokens tokens.CommandService,
	logger *zap.Logger,
	tracer trace.Tracer,
	txRunner database.TxRunner,
) *CommandService {
	return &CommandService{
		sessions: sessions,
		tokens:   tokens,
		logger:   logger,
		tracer:   tracer,
		txRunner: txRunner,
	}
}

// List handles the business logic for listing all sessions for the authenticated user.
func (uc *CommandService) List(ctx context.Context) ([]contracts.Session, error) {
	ctx, span := uc.tracer.Start(ctx, "SessionService.List")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var identityType contracts.IdentityType
	if principal.ProjectID == nil {
		identityType = contracts.ClientIdentity
	} else {
		identityType = contracts.ProjectIdentity
	}

	sessions, err := uc.sessions.List(ctx, principal.UserID, identityType)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("sessions.count", len(sessions)))

	return sessions, nil
}

// RevokeByID handles the business logic for revoking a specific session for the authenticated user.
// It ensures that the user is not revoking the current session.
func (uc *CommandService) RevokeByID(ctx context.Context, sessionID uuid.UUID, currentSessionID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "SessionService.RevokeByID")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	if currentSessionID == sessionID {
		return fail.New(errx.SessionSelfRevokeForbidden).RecordCtx(ctx)
	}

	var identityType contracts.IdentityType
	if principal.ProjectID == nil {
		identityType = contracts.ClientIdentity
	} else {
		identityType = contracts.ProjectIdentity
	}

	var sess *contracts.Session
	sess, err = uc.sessions.MarkRevokedByID(ctx, principal.UserID, sessionID, identityType)
	if fail.Is(err, errx.SQLNotFound) {
		return fail.New(errx.SessionNotFound).RecordCtx(ctx)
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
func (uc *CommandService) RevokeOthers(ctx context.Context, accessToken string) error {
	ctx, span := uc.tracer.Start(ctx, "SessionService.RevokeOthers")
	defer span.End()

	claims, err := uc.tokens.VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return err
	}

	currentSessionID := claims.Sub.SessionID

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var identityType contracts.IdentityType
	if principal.ProjectID == nil {
		identityType = contracts.ClientIdentity
	} else {
		identityType = contracts.ProjectIdentity
	}

	revokedCount, err := uc.sessions.MarkRevokedByFilter(ctx, contracts.Filter{
		IdentityType: identityType,
		EntityID:     principal.UserID,
		ExcludeID:    &currentSessionID,
	})
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.Int("sessions.revoked.count", revokedCount))
	return nil
}

// RevokeAll handles the business logic for revoking all sessions for the authenticated user.
func (uc *CommandService) RevokeAll(ctx context.Context) error {
	ctx, span := uc.tracer.Start(ctx, "SessionService.RevokeAll")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var identityType contracts.IdentityType
	if principal.ProjectID == nil {
		identityType = contracts.ClientIdentity
	} else {
		identityType = contracts.ProjectIdentity
	}

	revokedCount, err := uc.sessions.MarkRevokedByFilter(ctx, contracts.Filter{
		IdentityType: identityType,
		EntityID:     principal.UserID,
	})
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.Int("sessions.revoked.count", revokedCount))

	return nil
}
