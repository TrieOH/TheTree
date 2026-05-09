package sessions

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/feature_deps"
	"IdentityX/internal/shared/ports"
	"context"
	"errors"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	sessions ports.SessionRepository
	keys     ports.KeysRepository
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewCommandService(deps feature_deps.SessionCommandDeps) *CommandService {
	return errx.MustProvide(&CommandService{
		sessions: deps.Sessions,
		keys:     deps.Keys,
		logger:   deps.Logger,
		tracer:   deps.Tracer,
		tx:       deps.Tx,
	})
}

// RevokeByID handles the business logic for revoking a specific session for the authenticated user.
// It ensures that the user is not revoking the current session.
func (uc *CommandService) RevokeByID(ctx context.Context, sessionID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "SessionService.RevokeByID")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	if principal.Method == authz.AuthMethodApiKey {
		return errors.New("sessions are not revocable through api keys")
	}

	if *principal.SessionID == sessionID {
		return fun.ErrForbidden("cannot revoke the currently active session")
	}

	var userType contracts.UserType
	if principal.ProjectID == nil {
		userType = contracts.UserTypeClient
	} else {
		userType = contracts.UserTypeProject
	}

	var sess *contracts.Session
	sess, err = uc.sessions.MarkRevokedByID(ctx, principal.UserID, sessionID, userType)
	if fun.Is(err, fun.CodeNotFound) {
		return fun.ErrUnauthorized("session not found or revoked")
	}
	if err != nil {
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
func (uc *CommandService) RevokeOthers(ctx context.Context) error {
	ctx, span := uc.tracer.Start(ctx, "SessionService.RevokeOthers")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	if principal.Method == authz.AuthMethodApiKey {
		return errors.New("sessions are not revocable through api keys")
	}

	if principal.ProjectID != nil {
		span.SetAttributes(attribute.String("user.project_id", principal.ProjectID.String()))
	}

	var userType contracts.UserType
	if principal.ProjectID == nil {
		userType = contracts.UserTypeClient
	} else {
		userType = contracts.UserTypeProject
	}

	revokedCount, err := uc.sessions.MarkRevokedByFilter(ctx, contracts.Filter{
		UserType:  userType,
		UserID:    principal.UserID,
		ExcludeID: principal.SessionID,
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

	var userType contracts.UserType
	if principal.ProjectID == nil {
		userType = contracts.UserTypeClient
	} else {
		userType = contracts.UserTypeProject
	}

	revokedCount, err := uc.sessions.MarkRevokedByFilter(ctx, contracts.Filter{
		UserType: userType,
		UserID:   principal.UserID,
	})
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.Int("sessions.revoked.count", revokedCount))

	return nil
}
