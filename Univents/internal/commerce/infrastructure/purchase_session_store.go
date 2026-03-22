package infrastructure

import (
	"context"
	"encoding/json"
	"time"
	"univents/internal/commerce/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionStore struct {
	rdb *redis.Client
}

func NewSessionStore(rdb *redis.Client) *SessionStore {
	return &SessionStore{rdb: rdb}
}

func (s *SessionStore) Save(ctx context.Context, session domain.PurchaseSession) error {
	b, err := json.Marshal(session)
	if err != nil {
		return err
	}
	ttl := time.Until(session.ExpiresAt) + 2*time.Minute
	return s.rdb.Set(ctx, domain.PurchaseSessionKey(session.UserID, session.SessionID), b, ttl).Err()
}

func (s *SessionStore) Load(ctx context.Context, userID, sessionID uuid.UUID) (*domain.PurchaseSession, error) {
	b, err := s.rdb.Get(ctx, domain.PurchaseSessionKey(userID, sessionID)).Bytes()
	if err != nil {
		return nil, err
	}
	var session domain.PurchaseSession
	if err = json.Unmarshal(b, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionStore) Delete(ctx context.Context, userID, sessionID uuid.UUID) error {
	return s.rdb.Del(ctx, domain.PurchaseSessionKey(userID, sessionID)).Err()
}
