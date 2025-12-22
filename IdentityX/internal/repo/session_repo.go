package repo

import (
	"GoAuth/internal/models"
	"GoAuth/internal/sqlc"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type SessionRepo interface {
	Create(ctx context.Context, session models.Session) (*models.Session, error)
	GetById(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	GetByTokenId(ctx context.Context, tokenID uuid.UUID) (*models.Session, error)
	List(ctx context.Context, userID uuid.UUID) ([]models.Session, error)
	Update(ctx context.Context, session models.Session) (*models.Session, error)
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

func copySessionFromDB(dst *models.Session, src *sqlc.UserSession) error {
	return copier.Copy(dst, src)
}

func (s sessionRepo) Create(ctx context.Context, session models.Session) (*models.Session, error) {
	if session.UserID == uuid.Nil {
		return nil, errors.New("invalid user id")
	}

	sqlcSession, err := s.q.CreateUserSession(ctx, sqlc.CreateUserSessionParams{
		TokenID:   session.TokenID,
		UserID:    session.UserID,
		IssuedAt:  session.IssuedAt,
		UserAgent: session.UserAgent,
		UserIp:    session.UserIp,
		ExpiresAt: session.ExpiresAt,
		ProjectID: session.ProjectID,
	})

	if err != nil {
		return nil, err
	}

	if err = copySessionFromDB(&session, &sqlcSession); err != nil {
		s.log.Error(
			"failed to copy session",
			zap.Error(err),
			zap.String("session_id", sqlcSession.SessionID.String()),
		)
		return nil, fmt.Errorf("failed to copy session properly: %w", err)
	}

	return &session, nil
}

func (s sessionRepo) GetById(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	if sessionID == uuid.Nil {
		return nil, errors.New("session id is required for GetById")
	}

	sqlcSession, err := s.q.GetUserSessionById(ctx, sessionID)

	if err != nil {
		return nil, err
	}

	var session models.Session
	if err = copier.Copy(&session, sqlcSession); err != nil {
		s.log.Error(
			"failed to copy session",
			zap.Error(err),
			zap.String("session_id", sqlcSession.SessionID.String()),
		)
		return nil, fmt.Errorf("failed to copy session: %w", err)
	}

	return &session, nil
}

func (s sessionRepo) GetByTokenId(ctx context.Context, tokenID uuid.UUID) (*models.Session, error) {
	if tokenID == uuid.Nil {
		return nil, errors.New("token id is required for GetByTokenId")
	}

	sqlcSession, err := s.q.GetUserSessionByTokenId(ctx, tokenID)

	if err != nil {
		return nil, err
	}

	var session models.Session
	if err = copier.Copy(&session, sqlcSession); err != nil {
		s.log.Error(
			"failed to copy session",
			zap.Error(err),
			zap.String("session_id", sqlcSession.SessionID.String()),
		)
		return nil, fmt.Errorf("failed to copy session: %w", err)
	}

	return &session, nil
}

func (s sessionRepo) List(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user_id is required for list")
	}

	sqlcSessions, err := s.q.ListUserSessions(ctx, userID)

	if err != nil {
		return nil, err
	}

	var sessions []models.Session
	for _, sqlcSession := range sqlcSessions {
		var session models.Session
		if err = copier.Copy(&session, sqlcSession); err != nil {
			s.log.Error(
				"failed to copy session",
				zap.Error(err),
				zap.String("session_id", sqlcSession.SessionID.String()),
			)
			return nil, fmt.Errorf("failed to copy session of ID{%v}: %w", sqlcSession.SessionID, err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s sessionRepo) Update(ctx context.Context, session models.Session) (*models.Session, error) {
	if session.SessionID == uuid.Nil {
		return nil, errors.New("session_id is required for update")
	}

	sqlcSession, err := s.q.UpdateUserSession(ctx, sqlc.UpdateUserSessionParams{
		SessionID: session.SessionID,
		IssuedAt:  session.IssuedAt,
		UserAgent: session.UserAgent,
		UserIp:    session.UserIp,
		ExpiresAt: session.ExpiresAt,
		TokenID:   session.TokenID,
	})

	if err != nil {
		return nil, err
	}

	if err = copySessionFromDB(&session, &sqlcSession); err != nil {
		s.log.Error(
			"failed to copy session",
			zap.Error(err),
			zap.String("session_id", sqlcSession.SessionID.String()),
		)
		return nil, fmt.Errorf("failed to copy session of ID{%v}: %w", sqlcSession.SessionID, err)
	}

	return &session, nil
}

func (s sessionRepo) DeleteByFilter(ctx context.Context, filter models.SessionFilter) ([]models.Session, error) {
	if filter.UserID == uuid.Nil {
		return nil, errors.New("user_id is required for session deletion")
	}

	sqlcSessions, err := s.q.DeleteSessionsByFilter(ctx, sqlc.DeleteSessionsByFilterParams{
		UserID:        filter.UserID,
		SessionID:     filter.SessionID,
		ExcludeID:     filter.ExcludeID,
		TokenID:       filter.TokenID,
		ExpiredBefore: filter.ExpiredBefore,
	})

	if err != nil {
		return nil, err
	}

	var sessions []models.Session
	for _, sqlcSession := range sqlcSessions {
		var session models.Session
		if err = copier.Copy(&session, sqlcSession); err != nil {
			s.log.Error(
				"failed to copy session",
				zap.Error(err),
				zap.String("session_id", sqlcSession.SessionID.String()),
			)
			return nil, fmt.Errorf("failed to copy session of ID{%v}: %w", sqlcSession.SessionID, err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}
