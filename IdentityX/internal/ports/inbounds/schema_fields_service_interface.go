package inbounds

import (
	"GoAuth/internal/domain/field"
	"context"
)

type SchemaFieldsService interface {
	Create(ctx context.Context, in SchemaFieldInput) (CreateFieldsResult, error)
	EditField(ctx context.Context, in EditFieldInput) (*field.Field, error)
	DeleteField(ctx context.Context, in DeleteFieldInput) error
	SetFieldOptions(ctx context.Context, in SetFieldOptionsInput) ([]field.Option, error)
	DeleteFieldOption(ctx context.Context, in DeleteFieldOptionInput) error
}
