package sessions

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/security"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"
	"errors"

	"github.com/MintzyG/fail/v3"
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
	txRunner database.TxRunner
}

func NewCommandService(
	sessions ports.SessionRepository,
	keys ports.KeysRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	txRunner database.TxRunner,
) *CommandService {
	return &CommandService{
		sessions: sessions,
		keys:     keys,
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

	var userType contracts.UserType
	if principal.ProjectID == nil {
		userType = contracts.UserTypeClient
	} else {
		userType = contracts.UserTypeProject
	}

	sessions, err := uc.sessions.List(ctx, principal.UserID, userType)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("sessions.count", len(sessions)))

	return sessions, nil
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
		return errors.New("sessions are not revocable through api security")
	}

	if *principal.SessionID == sessionID {
		return fail.New(errx.SessionSelfRevokeForbidden).RecordCtx(ctx)
	}

	var userType contracts.UserType
	if principal.ProjectID == nil {
		userType = contracts.UserTypeClient
	} else {
		userType = contracts.UserTypeProject
	}

	var sess *contracts.Session
	sess, err = uc.sessions.MarkRevokedByID(ctx, principal.UserID, sessionID, userType)
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
func (uc *CommandService) RevokeOthers(ctx context.Context, accessTokenStr string) error {
	ctx, span := uc.tracer.Start(ctx, "SessionService.RevokeOthers")
	defer span.End()

	accessToken := &contracts.AccessClaims{}
	_, err := security.ParseJWTUnverified[*contracts.AccessClaims](accessTokenStr, accessToken)
	if err != nil {
		return err
	}

	var keyPair *contracts.Pair
	if accessToken.Sub.ProjectID != nil {
		span.SetAttributes(attribute.String("user.project_id", accessToken.Sub.ProjectID.String()))
		keyPair, err = uc.keys.GetActiveProjectSigningKey(ctx, *accessToken.Sub.ProjectID)
		if err != nil {
			return err
		}
	} else {
		keyPair, err = uc.keys.GetActiveGoAuthSigningKey(ctx)
		if err != nil {
			return err
		}
	}

	claims, err := security.VerifyAccessToken(ctx, accessTokenStr, keyPair)
	if err != nil {
		return err
	}

	currentSessionID := claims.Sub.SessionID

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
		UserType:  userType,
		UserID:    principal.UserID,
		ExcludeID: &currentSessionID,
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
