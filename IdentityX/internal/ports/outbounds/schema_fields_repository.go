package outbounds

import (
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"context"

	"github.com/google/uuid"
)

type SchemaFieldsRepository interface {
	Create(ctx context.Context, toCreate field.Field) (*field.Field, error)
	Update(ctx context.Context, toUpdate field.Field) error
	GetByVersionID(ctx context.Context, schemaVersionID uuid.UUID) ([]field.Field, error)
	ListFromSchema(ctx context.Context, schemaID uuid.UUID) ([]field.Field, error)
	ListFromVersion(ctx context.Context, schemaID, versionID uuid.UUID) ([]field.Field, error)
	Delete(ctx context.Context, fieldID uuid.UUID) error
	CloneFromTo(ctx context.Context, fromVersionID, toVersionID uuid.UUID) error
	DiffVersionsState(ctx context.Context, fromVersionID, toVersionID uuid.UUID) (bool, error)
	DiffVersionsFullState(ctx context.Context, fromVersionID, toVersionID uuid.UUID) (schema.DiffResult, error)

	CreateOption(ctx context.Context, option field.Option) (*field.Option, error)
	CreateVisibilityRule(ctx context.Context, rule field.VisibilityRule) (*field.VisibilityRule, error)
	CreateRequiredRule(ctx context.Context, rule field.RequiredRule) (*field.RequiredRule, error)

	GetByVersionIDWithRelations(ctx context.Context, schemaVersionID uuid.UUID) ([]field.Field, error)
	ListFromVersionWithRelations(ctx context.Context, schemaID, versionID uuid.UUID) ([]field.Field, error)

	CreateBatch(ctx context.Context, toCreate []field.Field) error
	CreateOptionsBatch(ctx context.Context, options []field.Option) error
	CreateVisibilityRulesBatch(ctx context.Context, rules []field.VisibilityRule) error
	CreateRequiredRulesBatch(ctx context.Context, rules []field.RequiredRule) error

	GetOptionsByFieldIDs(ctx context.Context, fieldIDs []uuid.UUID) ([]field.Option, error)
	GetVisibilityRulesByFieldIDs(ctx context.Context, fieldIDs []uuid.UUID) ([]field.VisibilityRule, error)
	GetRequiredRulesByFieldIDs(ctx context.Context, fieldIDs []uuid.UUID) ([]field.RequiredRule, error)
}
