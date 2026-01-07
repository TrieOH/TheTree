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
	SetOptions(ctx context.Context, options []field.Option) error
	SetRequiredRules(ctx context.Context, rules []field.RequiredRule) error
	SetVisibilityRules(ctx context.Context, rules []field.VisibilityRule) error
	Delete(ctx context.Context, fieldID uuid.UUID) error
}
