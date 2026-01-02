package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/ports/outbound"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type sessionRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

var _ outbound.SessionRepository = (*sessionRepo)(nil)

func NewSessionRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) outbound.SessionRepository {
	return &sessionRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func mapSessionFromDB(dst *session.Session, src *sqlc.Session) {
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

func (r sessionRepo) Create(ctx context.Context, new session.Session) (*session.Session, error) {
	ctx, span := r.tracer.Start(ctx, "SessionRepo.Create",
		trace.WithAttributes(
			attribute.String("session.user_id", new.UserID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := r.q.CreateUserSession(ctx, sqlc.CreateUserSessionParams{
		UserID:    new.UserID,
		IssuedAt:  new.IssuedAt,
		UserAgent: new.UserAgent,
		UserIp:    new.UserIp,
		ExpiresAt: new.ExpiresAt,
		ProjectID: new.ProjectID,
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

	var created session.Session
	mapSessionFromDB(&created, &sqlcSession)

	span.SetAttributes(
		attribute.String("session.session_id", created.SessionID.String()),
		attribute.String("session.token_id", created.TokenID.String()),
		attribute.Bool("session.created", true),
	)
	span.SetStatus(codes.Ok, "session created")

	return &created, nil
}

func (r sessionRepo) GetById(ctx context.Context, sessionID uuid.UUID) (*session.Session, error) {
	ctx, span := r.tracer.Start(ctx, "SessionRepo.GetById",
		trace.WithAttributes(
			attribute.String("session_id", sessionID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := r.q.GetUserSessionById(ctx, sessionID)

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

	var sess session.Session
	mapSessionFromDB(&sess, &sqlcSession)

	return &sess, nil
}

func (r sessionRepo) GetByTokenID(ctx context.Context, tokenID uuid.UUID) (*session.Session, error) {
	ctx, span := r.tracer.Start(ctx, "SessionRepo.GetByTokenID",
		trace.WithAttributes(
			attribute.String("token_id", tokenID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := r.q.GetUserSessionByTokenId(ctx, tokenID)

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

	var sess session.Session
	mapSessionFromDB(&sess, &sqlcSession)

	return &sess, nil
}

func (r sessionRepo) List(ctx context.Context, userID uuid.UUID) ([]session.Session, error) {
	ctx, span := r.tracer.Start(ctx, "SessionRepo.List",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
		),
	)
	defer span.End()

	sqlcSessions, err := r.q.ListUserSessions(ctx, userID)

	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.Int("sessions.count", len(sqlcSessions)))

	sessions := make([]session.Session, 0, len(sqlcSessions))
	for _, sqlcSession := range sqlcSessions {
		var sess session.Session
		mapSessionFromDB(&sess, &sqlcSession)
		sessions = append(sessions, sess)
	}

	return sessions, nil
}

func (r sessionRepo) Update(ctx context.Context, updated session.Session) error {
	ctx, span := r.tracer.Start(ctx, "SessionRepo.Update",
		trace.WithAttributes(
			attribute.String("session.user_id", updated.UserID.String()),
			attribute.String("session.token_id", updated.TokenID.String()),
			attribute.String("session.session_id", updated.SessionID.String()),
		),
	)
	defer span.End()

	err := r.q.UpdateUserSession(ctx, sqlc.UpdateUserSessionParams{
		SessionID: updated.SessionID,
		IssuedAt:  updated.IssuedAt,
		UserAgent: updated.UserAgent,
		UserIp:    updated.UserIp,
		ExpiresAt: updated.ExpiresAt,
		TokenID:   updated.TokenID,
	})

	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (r sessionRepo) DeleteByFilter(ctx context.Context, filter session.Filter) ([]session.Session, error) {
	ctx, span := r.tracer.Start(ctx, "SessionRepo.DeleteByFilter",
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

	sqlcSessions, err := r.q.DeleteSessionsByFilter(ctx, sqlc.DeleteSessionsByFilterParams{
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

	sessions := make([]session.Session, 0, len(sqlcSessions))
	for _, sqlcSession := range sqlcSessions {
		var sess session.Session
		mapSessionFromDB(&sess, &sqlcSession)
		sessions = append(sessions, sess)
	}

	return sessions, nil
}
