package sessions

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ ports.SessionRepository = (*sessionRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.SessionRepository {
	return &sessionRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func mapSessionFromDB(src *sqlc.Session) *contracts.Session {
	return &contracts.Session{
		SessionID: src.SessionID,
		ProjectID: src.ProjectID,
		UserID:    src.UserID,
		UserType:  contracts.UserType(src.UserType),
		FamilyID:  src.FamilyID,
		TokenID:   src.TokenID,
		IssuedAt:  src.IssuedAt,
		UserAgent: src.UserAgent,
		UserIP:    src.UserIp,
		RevokedAt: src.RevokedAt,
		ExpiresAt: src.ExpiresAt,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}

func (repo *sessionRepo) Create(ctx context.Context, toCreate contracts.Session) (*contracts.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.Create")
	defer span.End()

	sqlcSession, err := repo.queries(ctx).CreateUserSession(ctx, sqlc.CreateUserSessionParams{
		IssuedAt:  toCreate.IssuedAt,
		UserAgent: toCreate.UserAgent,
		UserIp:    toCreate.UserIP,
		ExpiresAt: toCreate.ExpiresAt,
		ProjectID: toCreate.ProjectID,
		UserID:    toCreate.UserID,
	})

	if err != nil {
		return nil, errx.DB(err, "session")
	}

	span.SetAttributes(attribute.String("session.user_type", sqlcSession.UserType))
	if sqlcSession.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sqlcSession.ProjectID.String()))
	}

	created := mapSessionFromDB(&sqlcSession)

	span.SetAttributes(
		attribute.String("session.session_id", created.SessionID.String()),
		attribute.String("session.token_id", created.TokenID.String()),
		attribute.Bool("session.created", true),
	)
	span.SetStatus(codes.Ok, "session created")

	return created, nil
}

func (repo *sessionRepo) GetByID(ctx context.Context, sessionID uuid.UUID) (*contracts.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.GetByID",
		trace.WithAttributes(
			attribute.String("session_id", sessionID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).GetUserSessionByID(ctx, sessionID)

	if err != nil {
		return nil, errx.DB(err, "session")
	}

	span.SetAttributes(
		attribute.String("session.token_id", sqlcSession.TokenID.String()),
		attribute.String("session.user_type", sqlcSession.UserType),
	)

	if sqlcSession.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sqlcSession.ProjectID.String()))
	}

	sess := mapSessionFromDB(&sqlcSession)
	return sess, nil
}

func (repo *sessionRepo) GetByTokenID(ctx context.Context, tokenID uuid.UUID) (*contracts.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.GetByTokenID",
		trace.WithAttributes(
			attribute.String("token_id", tokenID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).GetUserSessionByTokenID(ctx, tokenID)

	if err != nil {
		return nil, errx.DB(err, "session")
	}

	span.SetAttributes(
		attribute.String("session.session_id", sqlcSession.SessionID.String()),
		attribute.String("session.user_type", sqlcSession.UserType),
	)

	if sqlcSession.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sqlcSession.ProjectID.String()))
	}

	sess := mapSessionFromDB(&sqlcSession)
	return sess, nil
}

func (repo *sessionRepo) GetByFamilyID(ctx context.Context, familyID uuid.UUID) (*contracts.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.GetByFamilyID",
		trace.WithAttributes(
			attribute.String("session.family_id", familyID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).GetSessionByFamilyID(ctx, familyID)
	if err != nil {
		return nil, errx.DB(err, "session")
	}

	span.SetAttributes(attribute.String("session.session_id", sqlcSession.SessionID.String()))
	if sqlcSession.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sqlcSession.ProjectID.String()))
	}

	sess := mapSessionFromDB(&sqlcSession)
	return sess, nil
}

func (repo *sessionRepo) List(ctx context.Context, userID uuid.UUID, userType contracts.UserType) ([]contracts.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.List",
		trace.WithAttributes(
			attribute.String("entity_id", userID.String()),
			attribute.String("user_type", string(userType)),
		),
	)
	defer span.End()

	sqlcSessions, err := repo.queries(ctx).ListSessions(ctx, sqlc.ListSessionsParams{
		UserType: string(userType),
		UserID:   userID,
	})

	if err != nil {
		return nil, errx.DB(err, "session")
	}

	span.SetAttributes(attribute.Int("sessions.count", len(sqlcSessions)))

	sessions := make([]contracts.Session, 0, len(sqlcSessions))
	for _, sqlcSession := range sqlcSessions {
		sess := mapSessionFromDB(&sqlcSession)
		sessions = append(sessions, *sess)
	}

	return sessions, nil
}

func (repo *sessionRepo) Update(ctx context.Context, toUpdate contracts.Session, userID uuid.UUID, userType contracts.UserType) error {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.Update",
		trace.WithAttributes(
			attribute.String("session.user_type", string(toUpdate.UserType)),
			attribute.String("session.token_id", toUpdate.TokenID.String()),
			attribute.String("session.session_id", toUpdate.SessionID.String()),
		),
	)
	defer span.End()

	err := repo.queries(ctx).UpdateSession(ctx, sqlc.UpdateSessionParams{
		SessionID: toUpdate.SessionID,
		UserType:  string(userType),
		UserID:    userID,
		IssuedAt:  toUpdate.IssuedAt,
		UserAgent: toUpdate.UserAgent,
		UserIp:    toUpdate.UserIP,
		ExpiresAt: toUpdate.ExpiresAt,
		TokenID:   toUpdate.TokenID,
	})

	if err != nil {
		return errx.DB(err, "session")
	}

	return nil
}

func (repo *sessionRepo) RotateToken(ctx context.Context, familyID uuid.UUID, newTokenID uuid.UUID, oldTokenID uuid.UUID, expiresAt time.Time) (*contracts.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.RotateToken",
		trace.WithAttributes(
			attribute.String("family_id", familyID.String()),
			attribute.String("new_token_id", newTokenID.String()),
			attribute.String("old_token_id", oldTokenID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).RotateSessionToken(ctx, sqlc.RotateSessionTokenParams{
		ExpiresAt:  expiresAt,
		NewTokenID: newTokenID,
		OldTokenID: oldTokenID,
		FamilyID:   familyID,
	})
	if err != nil {
		return nil, errx.DB(err, "session")
	}

	span.SetAttributes(
		attribute.String("session.session_id", sqlcSession.SessionID.String()),
	)

	rotatedSession := mapSessionFromDB(&sqlcSession)
	return rotatedSession, nil
}

func (repo *sessionRepo) MarkRevokedByID(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, userType contracts.UserType) (*contracts.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.MarkRevokedByID",
		trace.WithAttributes(
			attribute.String("session_id", sessionID.String()),
			attribute.String("user_id", userID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).RevokeSessionByID(ctx, sqlc.RevokeSessionByIDParams{
		SessionID: sessionID,
		UserType:  string(userType),
		UserID:    userID,
	})
	if err != nil {
		return nil, errx.DB(err, "session")
	}

	revokedSession := mapSessionFromDB(&sqlcSession)
	return revokedSession, nil
}

func (repo *sessionRepo) MarkRevokedByFamilyID(ctx context.Context, familyID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.MarkRevokedByFamilyID",
		trace.WithAttributes(
			attribute.String("family_id", familyID.String()),
		),
	)
	defer span.End()

	if err := repo.queries(ctx).RevokeSessionByFamilyID(ctx, familyID); err != nil {
		return errx.DB(err, "session")
	}

	return nil
}

func (repo *sessionRepo) MarkRevokedByFilter(ctx context.Context, filter contracts.Filter) (int, error) {
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
			UserType:  string(filter.UserType),
			UserID:    filter.UserID,
			SessionID: *filter.ExcludeID,
		})
	} else {
		revokeType = "all"
		sqlcSessions, err = repo.queries(ctx).RevokeAllSessions(ctx, sqlc.RevokeAllSessionsParams{
			UserType: string(filter.UserType),
			UserID:   filter.UserID,
		})
	}

	if err != nil {
		return 0, errx.DB(err, "session")
	}

	span.SetAttributes(attribute.Int("revoke.count", len(sqlcSessions)))
	span.SetAttributes(attribute.String("revoke.type", revokeType))

	return len(sqlcSessions), nil
}
