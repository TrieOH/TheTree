package inbounds

import (
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type SchemaService interface {
	Draft(ctx context.Context, in SchemaServiceInput) (*SchemaOutput, error)
	Publish(ctx context.Context, in SchemaServiceInput) error
	GetByID(ctx context.Context, in SchemaServiceInput) (*SchemaOutput, error)
	GetVerbose(ctx context.Context, in SchemaServiceInput) (*SchemaVerboseOutput, error)
	GetIDsFromProjectID(ctx context.Context, projectID uuid.UUID) ([]uuid.UUID, error)
	List(ctx context.Context, projectID uuid.UUID) ([]SchemaOutput, error)
	GetLatestForm(ctx context.Context, in SchemaServiceInput) (*FormOutput, error)
	GetFormByVersion(ctx context.Context, in SchemaServiceInput, versionNumber int) (*FormOutput, error)
	CheckSchemaCompatibility(ctx context.Context, userID, projectID uuid.UUID) (bool, error)
	GetUpgradeForm(ctx context.Context) ([]FormResponse, error)
	ValidateAndConstructMetadata(ctx context.Context, projectID uuid.UUID, schemaType schema.Type, flowID string, customFields *json.RawMessage) (*json.RawMessage, error)
	ValidateFields(ctx context.Context, custom map[string]any, fieldDefs map[string]field.Field, registerFields []field.Field) (map[string]any, error)
	UpdateMetadata(ctx context.Context, customFields *json.RawMessage) error
}
