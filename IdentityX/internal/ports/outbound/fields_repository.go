package outbound

import (
	"GoAuth/internal/domain/field"
	"context"

	"github.com/google/uuid"
)

type FieldsRepository interface {
	Create(ctx context.Context, field field.Field) (*field.Field, error)
	Update(ctx context.Context, field field.Field) error
	SetOptions(ctx context.Context, options []field.Option) error
	SetRequiredRules(ctx context.Context, required []field.RequiredRule) error
	SetVisibilityRules(ctx context.Context, visibilityRules []field.VisibilityRule) error
	Delete(ctx context.Context, FieldID uuid.UUID) error
}
