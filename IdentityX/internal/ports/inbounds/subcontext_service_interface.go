package inbounds

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type SubContextService interface {
	Add(ctx context.Context, projectID, userID uuid.UUID, data json.RawMessage) error
	Remove(ctx context.Context, projectID, userID uuid.UUID, keys []string) error
	Get(ctx context.Context, projectID, userID uuid.UUID) (json.RawMessage, error)
}
