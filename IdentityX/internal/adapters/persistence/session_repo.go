package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/ports/outbounds"
	"context"
	"time"

	"github.com/MintzyG/fail"
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
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbounds.SessionRepository = (*sessionRepo)(nil)

func NewSessionRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) outbounds.SessionRepository {
	return &sessionRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func mapSessionFromDB(dst *session.Session, src *sqlc.Session) {
	dst.SessionID = src.SessionID
	dst.IdentityID = src.IdentityID
	dst.FamilyID = src.FamilyID
	dst.ProjectID = src.ProjectID
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

func mapSessionIdentityFromDB(dst *session.Identity, src *sqlc.Identity) {
	dst.ID = src.ID
	dst.IdentityType = session.IdentityType(src.Type)
	dst.EntityID = src.EntityID
	dst.CreatedAt = src.CreatedAt
}

func (repo *sessionRepo) Create(ctx context.Context, toCreate session.Session) (*session.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.Create",
		trace.WithAttributes(
			attribute.String("session.identity_id", toCreate.IdentityID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).CreateUserSession(ctx, sqlc.CreateUserSessionParams{
		IdentityID: toCreate.IdentityID,
		IssuedAt:   toCreate.IssuedAt,
		UserAgent:  toCreate.UserAgent,
		UserIp:     toCreate.UserIP,
		ExpiresAt:  toCreate.ExpiresAt,
		ProjectID:  toCreate.ProjectID,
	})

	if err != nil {
		return nil, fail.From(err)
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
		return nil, fail.From(err)
	}

	span.SetAttributes(
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
		return nil, fail.From(err)
	}

	span.SetAttributes(
		attribute.String("session.session_id", sqlcSession.SessionID.String()),
		attribute.String("session.user_type", sqlcSession.UserType),
	)

	if sqlcSession.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sqlcSession.ProjectID.String()))
	}

	var sess session.Session
	mapSessionFromDB(&sess, &sqlcSession)

	return &sess, nil
}

func (repo *sessionRepo) GetByFamilyID(ctx context.Context, familyID uuid.UUID) (*session.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.GetByFamilyID",
		trace.WithAttributes(
			attribute.String("session.family_id", familyID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).GetSessionByFamilyID(ctx, familyID)
	if err != nil {
		return nil, fail.From(err)
	}

	span.SetAttributes(attribute.String("session.session_id", sqlcSession.SessionID.String()))
	if sqlcSession.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sqlcSession.ProjectID.String()))
	}

	var sess session.Session
	mapSessionFromDB(&sess, &sqlcSession)
	return &sess, nil
}

func (repo *sessionRepo) List(ctx context.Context, entityID uuid.UUID, identityType session.IdentityType) ([]session.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.List",
		trace.WithAttributes(
			attribute.String("entity_id", entityID.String()),
			attribute.String("identity_type", string(identityType)),
		),
	)
	defer span.End()

	sqlcSessions, err := repo.queries(ctx).ListSessions(ctx, sqlc.ListSessionsParams{
		Type:     sqlc.IdentityType(identityType),
		EntityID: entityID,
	})

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

func (repo *sessionRepo) Update(ctx context.Context, toUpdate session.Session, entityID uuid.UUID, identityType session.IdentityType) error {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.Update",
		trace.WithAttributes(
			attribute.String("session.identity_id", toUpdate.IdentityID.String()),
			attribute.String("session.token_id", toUpdate.TokenID.String()),
			attribute.String("session.session_id", toUpdate.SessionID.String()),
		),
	)
	defer span.End()

	err := repo.queries(ctx).UpdateSession(ctx, sqlc.UpdateSessionParams{
		SessionID: toUpdate.SessionID,
		Type:      sqlc.IdentityType(identityType),
		EntityID:  entityID,
		IssuedAt:  toUpdate.IssuedAt,
		UserAgent: toUpdate.UserAgent,
		UserIp:    toUpdate.UserIP,
		ExpiresAt: toUpdate.ExpiresAt,
		TokenID:   toUpdate.TokenID,
	})

	if err != nil {
		return fail.From(err)
	}

	return nil
}

func (repo *sessionRepo) RotateToken(ctx context.Context, familyID uuid.UUID, newTokenID uuid.UUID, oldTokenID uuid.UUID, expiresAt time.Time) (*session.Session, error) {
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
		return nil, fail.From(err)
	}

	span.SetAttributes(
		attribute.String("session.session_id", sqlcSession.SessionID.String()),
	)

	var rotatedSession session.Session
	mapSessionFromDB(&rotatedSession, &sqlcSession)
	return &rotatedSession, nil
}

func (repo *sessionRepo) MarkRevokedByID(ctx context.Context, entityID uuid.UUID, sessionID uuid.UUID, identityType session.IdentityType) (*session.Session, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.MarkRevokedByID",
		trace.WithAttributes(
			attribute.String("session_id", sessionID.String()),
			attribute.String("entity_id", entityID.String()),
		),
	)
	defer span.End()

	sqlcSession, err := repo.queries(ctx).RevokeSessionByID(ctx, sqlc.RevokeSessionByIDParams{
		SessionID: sessionID,
		Type:      sqlc.IdentityType(identityType),
		EntityID:  entityID,
	})
	if err != nil {
		return nil, fail.From(err)
	}

	var revokedSession session.Session
	mapSessionFromDB(&revokedSession, &sqlcSession)
	return &revokedSession, nil
}

func (repo *sessionRepo) MarkRevokedByFamilyID(ctx context.Context, familyID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.MarkRevokedByFamilyID",
		trace.WithAttributes(
			attribute.String("family_id", familyID.String()),
		),
	)
	defer span.End()

	if err := repo.queries(ctx).RevokeSessionByFamilyID(ctx, familyID); err != nil {
		return fail.From(err)
	}

	return nil
}

func (repo *sessionRepo) MarkRevokedByFilter(ctx context.Context, filter session.Filter) (int, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.MarkRevokedByFilter",
		trace.WithAttributes(
			attribute.String("user_id", filter.EntityID.String()),
		),
	)
	defer span.End()

	var err error
	var revokeType string
	var sqlcSessions []sqlc.Session
	if filter.ExcludeID != nil {
		revokeType = "other"
		sqlcSessions, err = repo.queries(ctx).RevokeOtherSessions(ctx, sqlc.RevokeOtherSessionsParams{
			Type:      sqlc.IdentityType(filter.IdentityType),
			EntityID:  filter.EntityID,
			SessionID: *filter.ExcludeID,
		})
	} else {
		revokeType = "all"
		sqlcSessions, err = repo.queries(ctx).RevokeAllSessions(ctx, sqlc.RevokeAllSessionsParams{
			Type:     sqlc.IdentityType(filter.IdentityType),
			EntityID: filter.EntityID,
		})
	}

	if err != nil {
		return 0, fail.From(err)
	}

	span.SetAttributes(attribute.Int("revoke.count", len(sqlcSessions)))
	span.SetAttributes(attribute.String("revoke.type", revokeType))

	return len(sqlcSessions), nil
}

func (repo *sessionRepo) CreateIdentity(ctx context.Context, identityType session.IdentityType, entityID uuid.UUID) (*session.Identity, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.CreateIdentity",
		trace.WithAttributes(
			attribute.String("entity_id", entityID.String()),
			attribute.String("identity_type", string(identityType)),
		),
	)
	defer span.End()

	sqlcIdentity, err := repo.queries(ctx).CreateSessionIdentity(ctx, sqlc.CreateSessionIdentityParams{
		Type:     sqlc.IdentityType(identityType),
		EntityID: entityID,
	})
	if err != nil {
		return nil, fail.From(err)
	}

	span.SetAttributes(attribute.String("identity_id", sqlcIdentity.ID.String()))

	var foundIdentity session.Identity
	mapSessionIdentityFromDB(&foundIdentity, &sqlcIdentity)
	return &foundIdentity, nil
}

func (repo *sessionRepo) GetIdentityByEntityIDAndType(ctx context.Context, entityID uuid.UUID, identityType session.IdentityType) (*session.Identity, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.GetByEntityIDAndType",
		trace.WithAttributes(
			attribute.String("entity_id", entityID.String()),
			attribute.String("identity_type", string(identityType)),
		),
	)
	defer span.End()

	sqlcIdentity, err := repo.queries(ctx).GetSessionIdentityByEntityIDAndType(ctx, sqlc.GetSessionIdentityByEntityIDAndTypeParams{
		Type:     sqlc.IdentityType(identityType),
		EntityID: entityID,
	})
	if err != nil {
		return nil, fail.From(err)
	}

	span.SetAttributes(attribute.String("identity_id", sqlcIdentity.ID.String()))

	var foundIdentity session.Identity
	mapSessionIdentityFromDB(&foundIdentity, &sqlcIdentity)
	return &foundIdentity, nil
}

func (repo *sessionRepo) GetIdentityByIDAndType(ctx context.Context, identityID uuid.UUID, identityType session.IdentityType) (*session.Identity, error) {
	ctx, span := repo.tracer.Start(ctx, "SessionRepo.GetByIDAndType",
		trace.WithAttributes(
			attribute.String("identity_id", identityID.String()),
			attribute.String("identity_type", string(identityType)),
		),
	)
	defer span.End()

	sqlcIdentity, err := repo.queries(ctx).GetSessionIdentityByIDAndType(ctx, sqlc.GetSessionIdentityByIDAndTypeParams{
		Type: sqlc.IdentityType(identityType),
		ID:   identityID,
	})
	if err != nil {
		return nil, fail.From(err)
	}

	var foundIdentity session.Identity
	mapSessionIdentityFromDB(&foundIdentity, &sqlcIdentity)
	return &foundIdentity, nil
}
