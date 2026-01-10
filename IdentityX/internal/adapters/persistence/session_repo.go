package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/ports/outbound"
	"context"
	"database/sql"
	"time"

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

func (repo *sessionRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(txKeyValue).(*sql.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
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
	dst.UserIP = src.UserIp
	dst.RevokedAt = src.RevokedAt
	dst.ExpiresAt = src.ExpiresAt
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
	dst.UserType = src.UserType
}

func (repo *sessionRepo) Create(ctx context.Context, toCreate session.Session) (*session.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.Create",
		trace.WithAttributes(
			attribute.String("session.user_id", toCreate.UserID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).CreateUserSession(ctx, sqlc.CreateUserSessionParams{
		UserID:    toCreate.UserID,
		IssuedAt:  toCreate.IssuedAt,
		UserAgent: toCreate.UserAgent,
		UserIp:    toCreate.UserIP,
		ExpiresAt: toCreate.ExpiresAt,
		ProjectID: toCreate.ProjectID,
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

func (repo *sessionRepo) GetByID(ctx context.Context, sessionID uuid.UUID) (*session.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.GetByID",
		trace.WithAttributes(
			attribute.String("session_id", sessionID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).GetUserSessionByID(ctx, sessionID)

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

func (repo *sessionRepo) GetByTokenID(ctx context.Context, tokenID uuid.UUID) (*session.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.GetByTokenID",
		trace.WithAttributes(
			attribute.String("token_id", tokenID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).GetUserSessionByTokenID(ctx, tokenID)

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

func (repo *sessionRepo) List(ctx context.Context, userID uuid.UUID) ([]session.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.List",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
		),
	)
	defer span.End()

	sqlcSessions, err := repo.queries(ctx).ListUserSessions(ctx, userID)

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

func (repo *sessionRepo) Update(ctx context.Context, toUpdate session.Session) error {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.Update",
		trace.WithAttributes(
			attribute.String("session.user_id", toUpdate.UserID.String()),
			attribute.String("session.token_id", toUpdate.TokenID.String()),
			attribute.String("session.session_id", toUpdate.SessionID.String()),
		),
	)
	defer span.End()

	err := repo.queries(ctx).UpdateUserSession(ctx, sqlc.UpdateUserSessionParams{
		SessionID: toUpdate.SessionID,
		IssuedAt:  toUpdate.IssuedAt,
		UserAgent: toUpdate.UserAgent,
		UserIp:    toUpdate.UserIP,
		ExpiresAt: toUpdate.ExpiresAt,
		TokenID:   toUpdate.TokenID,
	})

	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (repo *sessionRepo) RotateToken(ctx context.Context, oldTokenID uuid.UUID, newTokenID uuid.UUID, expiresAt time.Time) (*session.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.RotateToken",
		trace.WithAttributes(
			attribute.String("old_token_id", oldTokenID.String()),
			attribute.String("new_token_id", newTokenID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).RotateSessionToken(ctx, sqlc.RotateSessionTokenParams{
		ExpiresAt:  expiresAt,
		NewTokenID: newTokenID,
		OldTokenID: oldTokenID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(
		attribute.String("session.session_id", sqlcSession.SessionID.String()),
		attribute.String("session.user_id", sqlcSession.UserID.String()),
	)

	var rotatedSession session.Session
	mapSessionFromDB(&rotatedSession, &sqlcSession)
	return &rotatedSession, nil
}

func (repo *sessionRepo) MarkRevokedByID(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (*session.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.MarkRevokedByID",
		trace.WithAttributes(
			attribute.String("session_id", sessionID.String()),
			attribute.String("user_id", userID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).RevokeSessionByID(ctx, sqlc.RevokeSessionByIDParams{
		SessionID: sessionID,
		UserID:    userID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	var revokedSession session.Session
	mapSessionFromDB(&revokedSession, &sqlcSession)
	return &revokedSession, nil
}

func (repo *sessionRepo) MarkRevokedByFilter(ctx context.Context, filter session.Filter) (int, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.MarkRevokedByFilter",
		trace.WithAttributes(
			attribute.String("user_id", filter.UserID.String()),
		),
	)
	defer span.End()

	var err error
	var revokeType string
	var sqlcSessions []sqlc.Session
	if filter.ExcludeID != nil {
		revokeType = "other"
		sqlcSessions, err = repo.queries(ctx).RevokeOtherSessions(ctx, sqlc.RevokeOtherSessionsParams{
			UserID:    filter.UserID,
			SessionID: *filter.ExcludeID,
		})
	} else {
		revokeType = "all"
		sqlcSessions, err = repo.queries(ctx).RevokeAllSessions(ctx, filter.UserID)
	}

	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return 0, sqlcErr
	}

	span.SetAttributes(attribute.Int("revoke.count", len(sqlcSessions)))
	span.SetAttributes(attribute.String("revoke.type", revokeType))

	return len(sqlcSessions), nil
}
