package repo

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/models"
	"GoAuth/internal/sqlc"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type SessionRepo interface {
	Create(ctx context.Context, session models.Session) (*models.Session, error)
	GetById(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	GetByTokenID(ctx context.Context, tokenID uuid.UUID) (*models.Session, error)
	List(ctx context.Context, userID uuid.UUID) ([]models.Session, error)
	Update(ctx context.Context, session models.Session) error
	DeleteByFilter(ctx context.Context, filter models.SessionFilter) ([]models.Session, error)
}

type sessionRepo struct {
	q   *sqlc.Queries
	log *zap.Logger
}

func NewSessionRepo(q *sqlc.Queries, log *zap.Logger) SessionRepo {
	return &sessionRepo{
		q:   q,
		log: log,
	}
}

func mapSessionFromDB(dst *models.Session, src *sqlc.Session) {
	dst.SessionID = src.SessionID
	dst.ProjectID = src.ProjectID
	dst.UserID = src.UserID
	dst.TokenID = src.TokenID
	dst.IssuedAt = src.IssuedAt
	dst.UserAgent = src.UserAgent
	dst.UserIp = src.UserIp
	dst.ExpiresAt = src.ExpiresAt
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
	dst.UserType = src.UserType
}

func (s sessionRepo) Create(ctx context.Context, session models.Session) (*models.Session, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "SessionRepo.Create",
		trace.WithAttributes(
			attribute.String("session.user_id", session.UserID.String()),
			attribute.String("session.token_id", session.TokenID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := s.q.CreateUserSession(ctx, sqlc.CreateUserSessionParams{
		UserID:    session.UserID,
		IssuedAt:  session.IssuedAt,
		UserAgent: session.UserAgent,
		UserIp:    session.UserIp,
		ExpiresAt: session.ExpiresAt,
		ProjectID: session.ProjectID,
	})

	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return nil, sqlErr
	}

	span.SetAttributes(attribute.String("session.user_type", sqlcSession.UserType))
	if sqlcSession.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sqlcSession.ProjectID.String()))
	}

	mapSessionFromDB(&session, &sqlcSession)

	span.SetAttributes(
		attribute.String("session.session_id", session.SessionID.String()),
		attribute.Bool("session.created", true),
	)
	span.SetStatus(codes.Ok, "session created")

	return &session, nil
}

func (s sessionRepo) GetById(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "SessionRepo.GetById",
		trace.WithAttributes(
			attribute.String("session_id", sessionID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := s.q.GetUserSessionById(ctx, sessionID)

	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(
		attribute.String("session.user_id", sqlcSession.UserID.String()),
		attribute.String("session.token_id", sqlcSession.TokenID.String()),
		attribute.String("session.user_type", sqlcSession.UserType),
	)

	if sqlcSession.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sqlcSession.ProjectID.String()))
	}

	var session models.Session
	mapSessionFromDB(&session, &sqlcSession)

	return &session, nil
}

func (s sessionRepo) GetByTokenID(ctx context.Context, tokenID uuid.UUID) (*models.Session, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "SessionRepo.GetByTokenID",
		trace.WithAttributes(
			attribute.String("token_id", tokenID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := s.q.GetUserSessionByTokenId(ctx, tokenID)

	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(
		attribute.String("session.session_id", sqlcSession.SessionID.String()),
		attribute.String("session.user_id", sqlcSession.UserID.String()),
		attribute.String("session.user_type", sqlcSession.UserType),
	)

	if sqlcSession.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sqlcSession.ProjectID.String()))
	}

	var session models.Session
	mapSessionFromDB(&session, &sqlcSession)

	return &session, nil
}

func (s sessionRepo) List(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "SessionRepo.List",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
		),
	)
	defer span.End()

	sqlcSessions, err := s.q.ListUserSessions(ctx, userID)

	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.Int("sessions.count", len(sqlcSessions)))

	sessions := make([]models.Session, 0, len(sqlcSessions))
	for _, sqlcSession := range sqlcSessions {
		var session models.Session
		mapSessionFromDB(&session, &sqlcSession)
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s sessionRepo) Update(ctx context.Context, session models.Session) error {
	ctx, span := GoAuthRepoTracer.Start(ctx, "SessionRepo.Update",
		trace.WithAttributes(
			attribute.String("session.user_id", session.UserID.String()),
			attribute.String("session.token_id", session.TokenID.String()),
		),
	)
	defer span.End()

	err := s.q.UpdateUserSession(ctx, sqlc.UpdateUserSessionParams{
		SessionID: session.SessionID,
		IssuedAt:  session.IssuedAt,
		UserAgent: session.UserAgent,
		UserIp:    session.UserIp,
		ExpiresAt: session.ExpiresAt,
		TokenID:   session.TokenID,
	})

	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (s sessionRepo) DeleteByFilter(ctx context.Context, filter models.SessionFilter) ([]models.Session, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "SessionRepo.DeleteByFilter",
		trace.WithAttributes(
			attribute.String("session.user_id", filter.UserID.String()),
		),
	)

	if filter.SessionID != nil {
		span.SetAttributes(attribute.String("session.session_id", filter.SessionID.String()))
	}
	if filter.TokenID != nil {
		span.SetAttributes(attribute.String("session.with_token_id", filter.TokenID.String()))
	}
	if filter.ExcludeID != nil {
		span.SetAttributes(attribute.String("session.with_exclude_id", filter.ExcludeID.String()))
	}
	if filter.ExpiredBefore != nil {
		span.SetAttributes(attribute.String("session.with_expired", filter.ExpiredBefore.String()))
	}

	defer span.End()

	sqlcSessions, err := s.q.DeleteSessionsByFilter(ctx, sqlc.DeleteSessionsByFilterParams{
		UserID:        filter.UserID,
		SessionID:     filter.SessionID,
		ExcludeID:     filter.ExcludeID,
		TokenID:       filter.TokenID,
		ExpiredBefore: filter.ExpiredBefore,
	})

	if err != nil {
		sessionErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sessionErr)
		return nil, sessionErr
	}

	span.SetAttributes(attribute.Int("sessions.deleted", len(sqlcSessions)))

	sessions := make([]models.Session, 0, len(sqlcSessions))
	for _, sqlcSession := range sqlcSessions {
		var session models.Session
		mapSessionFromDB(&session, &sqlcSession)
		sessions = append(sessions, session)
	}

	return sessions, nil
}
