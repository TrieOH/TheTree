package outbound

import (
	"GoAuth/internal/domain/field"
	"context"

	"github.com/google/uuid"
)

type SchemaFieldsRepository interface {
	Create(ctx context.Context, toCreate field.Field) (*field.Field, error)
	Update(ctx context.Context, toUpdate field.Field) error
	GetByVersionID(ctx context.Context, schemaVersionID uuid.UUID) ([]field.Field, error)
	List(ctx context.Context, schemaID uuid.UUID) ([]field.Field, error)
	Delete(ctx context.Context, fieldID uuid.UUID) error
	CloneFromTo(ctx context.Context, fromVersionID, toVersionID uuid.UUID) error
	DiffVersionsState(ctx context.Context, fromVersionID, toVersionID uuid.UUID) (bool, error)
}
