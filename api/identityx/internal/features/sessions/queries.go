package sessions

import (
	"IdentityX/contracts"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/ports"
	"context"
	"lib/database"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type QueryService struct {
	sessions ports.SessionRepository
	logger   *zap.Logger
	tracer   trace.Tracer
	txRunner database.TxRunner
}

func NewQueryService(
	sessions ports.SessionRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	txRunner database.TxRunner,
) *QueryService {
	return &QueryService{
		sessions: sessions,
		logger:   logger,
		tracer:   tracer,
		txRunner: txRunner,
	}
}

// List handles the business logic for listing all sessions for the authenticated user.
func (uc *QueryService) List(ctx context.Context) ([]contracts.Session, error) {
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
