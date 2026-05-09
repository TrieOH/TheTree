package sessions

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"IdentityX/internal/shared/xslices"
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

func (repo *sessionRepo) span(ctx context.Context, op string) (context.Context, trace.Span) {
	return repo.tracer.Start(ctx, "SessionRepo."+op)
}

var _ ports.SessionRepository = (*sessionRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.SessionRepository {
	return &sessionRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func mapSessionFromDB(src sqlc.Session) contracts.Session {
	return contracts.Session{
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
	ctx, span := repo.span(ctx, "Create")
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
	created := mapSessionFromDB(sqlcSession)
	span.SetAttributes(
		attribute.String("session.session_id", created.SessionID.String()),
		attribute.String("session.token_id", created.TokenID.String()),
		attribute.Bool("session.created", true),
	)
	span.SetStatus(codes.Ok, "session created")
	return &created, nil
}

func (repo *sessionRepo) GetByID(ctx context.Context, sessionID uuid.UUID) (*contracts.Session, error) {
	ctx, span := repo.span(ctx, "GetByID")
	span.SetAttributes(attribute.String("session_id", sessionID.String()))
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
	return new(mapSessionFromDB(sqlcSession)), nil
}

func (repo *sessionRepo) GetByTokenID(ctx context.Context, tokenID uuid.UUID) (*contracts.Session, error) {
	ctx, span := repo.span(ctx, "GetByTokenID")
	span.SetAttributes(attribute.String("token_id", tokenID.String()))
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
	return new(mapSessionFromDB(sqlcSession)), nil
}

func (repo *sessionRepo) GetByFamilyID(ctx context.Context, familyID uuid.UUID) (*contracts.Session, error) {
	ctx, span := repo.span(ctx, "GetByFamilyID")
	span.SetAttributes(attribute.String("session.family_id", familyID.String()))
	defer span.End()
	sqlcSession, err := repo.queries(ctx).GetSessionByFamilyID(ctx, familyID)
	if err != nil {
		return nil, errx.DB(err, "session")
	}
	span.SetAttributes(attribute.String("session.session_id", sqlcSession.SessionID.String()))
	if sqlcSession.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sqlcSession.ProjectID.String()))
	}
	return new(mapSessionFromDB(sqlcSession)), nil
}

func (repo *sessionRepo) List(ctx context.Context, userID uuid.UUID, userType contracts.UserType) ([]contracts.Session, error) {
	ctx, span := repo.span(ctx, "List")
	span.SetAttributes(attribute.String("entity_id", userID.String()))
	span.SetAttributes(attribute.String("user_type", string(userType)))
	defer span.End()
	sqlcSessions, err := repo.queries(ctx).ListSessions(ctx, sqlc.ListSessionsParams{
		UserType: string(userType),
		UserID:   userID,
	})
	if err != nil {
		return nil, errx.DB(err, "session")
	}
	span.SetAttributes(attribute.Int("sessions.count", len(sqlcSessions)))
	return xslices.MapSlice(sqlcSessions, mapSessionFromDB), nil
}

func (repo *sessionRepo) Update(ctx context.Context, toUpdate contracts.Session, userID uuid.UUID, userType contracts.UserType) error {
	ctx, span := repo.span(ctx, "Update")
	span.SetAttributes(attribute.String("session.user_type", string(toUpdate.UserType)))
	span.SetAttributes(attribute.String("session.token_id", toUpdate.TokenID.String()))
	span.SetAttributes(attribute.String("session.session_id", toUpdate.SessionID.String()))
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
	ctx, span := repo.span(ctx, "RotateToken")
	span.SetAttributes(attribute.String("family_id", familyID.String()))
	span.SetAttributes(attribute.String("new_token_id", newTokenID.String()))
	span.SetAttributes(attribute.String("old_token_id", oldTokenID.String()))
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
	return new(mapSessionFromDB(sqlcSession)), nil
}

func (repo *sessionRepo) MarkRevokedByID(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, userType contracts.UserType) (*contracts.Session, error) {
	ctx, span := repo.span(ctx, "MarkRevokedByID")
	span.SetAttributes(attribute.String("session_id", sessionID.String()))
	span.SetAttributes(attribute.String("user_id", userID.String()))
	defer span.End()
	sqlcSession, err := repo.queries(ctx).RevokeSessionByID(ctx, sqlc.RevokeSessionByIDParams{
		SessionID: sessionID,
		UserType:  string(userType),
		UserID:    userID,
	})
	if err != nil {
		return nil, errx.DB(err, "session")
	}
	return new(mapSessionFromDB(sqlcSession)), nil
}

func (repo *sessionRepo) MarkRevokedByFamilyID(ctx context.Context, familyID uuid.UUID) error {
	ctx, span := repo.span(ctx, "MarkRevokedByFamilyID")
	span.SetAttributes(attribute.String("family_id", familyID.String()))
	defer span.End()
	if err := repo.queries(ctx).RevokeSessionByFamilyID(ctx, familyID); err != nil {
		return errx.DB(err, "session")
	}
	return nil
}

func (repo *sessionRepo) MarkRevokedByFilter(ctx context.Context, filter contracts.Filter) (int, error) {
	ctx, span := repo.span(ctx, "MarkRevokedByFilter")
	span.SetAttributes(attribute.String("user_id", filter.UserID.String()))
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
