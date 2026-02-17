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
	SetVisibilityRules(ctx context.Context, in SetVisibilityRulesInput) ([]field.VisibilityRule, error)
	EditVisibilityRule(ctx context.Context, in EditVisibilityRuleInput) (*field.VisibilityRule, error)
	DeleteVisibilityRule(ctx context.Context, in DeleteVisibilityRuleInput) error
	SetRequiredRules(ctx context.Context, in SetRequiredRulesInput) ([]field.RequiredRule, error)
	EditRequiredRule(ctx context.Context, in EditRequiredRuleInput) (*field.RequiredRule, error)
	DeleteRequiredRule(ctx context.Context, in DeleteRequiredRuleInput) error
	BatchUpdateFields(ctx context.Context, in BatchUpdateFieldsInput) (BatchUpdateFieldsResult, error)
}
