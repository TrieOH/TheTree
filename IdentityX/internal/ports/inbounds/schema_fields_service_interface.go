package inbounds

import (
	"context"
)

type SchemaFieldsService interface {
	Create(ctx context.Context, in SchemaFieldInput) ([]OutputField, error)
}
