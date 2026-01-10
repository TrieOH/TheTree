package auth

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

// validateAndConstructMetadata validates custom fields against a schema and returns structured metadata
func (uc *UseCase) validateAndConstructMetadata(
	ctx context.Context,
	span trace.Span,
	projectID uuid.UUID,
	schemaType, flowID string,
	customFields *json.RawMessage,
) (*json.RawMessage, error) {
	var err error
	var registerSchema *schema.Schema
	registerSchema, err = uc.schemas.FindByFlowIDAndType(ctx, flowID, schemaType, projectID)
	if err != nil {
		return nil, err
	}

	if registerSchema.CurrentVersionID == nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("schema has no published version").WithID(apierr.SchemaNoPublishedVersion)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	if registerSchema.Status == schema.StatusDraft {
		apiErr := apierr.ErrBadRequest.WithMsg("can't register to a draft schema").WithID(apierr.ProjectUserRegisterOnSchemaDraft)
		apierr.RecordSystemError(span, apiErr)
		return nil, apiErr
	}

	var registerVersion *schema.Version
	registerVersion, err = uc.versions.GetCurrent(ctx, registerSchema.ID)
	if err != nil {
		return nil, err
	}

	if registerVersion.Status == schema.VersionStatusDraft {
		apiErr := apierr.ErrBadRequest.WithMsg("can't register to a draft schema version").WithID(apierr.ProjectUserRegisterOnSchemaVersionDraft)
		apierr.RecordSystemError(span, apiErr)
		return nil, apiErr
	}

	if registerVersion.ID != *registerSchema.CurrentVersionID {
		apiErr := apierr.ErrInternal.WithMsg("schema version and retrieved version mismatch").WithID(apierr.SchemaVersionMismatch)
		apierr.RecordSystemError(span, apiErr)
		return nil, apiErr
	}

	var registerFields []field.Field
	registerFields, err = uc.fields.GetByVersionID(ctx, registerVersion.ID)
	if err != nil {
		return nil, err
	}

	fieldDefs := make(map[string]field.Field)
	for _, f := range registerFields {
		fieldDefs[f.Key] = f
	}

	var custom map[string]any

	if customFields == nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("the schema custom fields are required on a schema register").WithID(apierr.RequestMissingSchemaCustomFields)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	if err := json.Unmarshal(*customFields, &custom); err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid custom fields JSON").WithID(apierr.RequestInvalidJSON).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	validated := make(map[string]any)
	for key, value := range custom {
		f, ok := fieldDefs[key]
		if !ok {
			apiErr := apierr.ErrInvalidInput.WithMsg("unknown custom field").WithID(apierr.FieldNotDefinedInSchema).WithCause(errors.New("unknown field: " + key))
			apierr.RecordDomainError(span, apiErr)
			return nil, apiErr
		}

		if !validateFieldType(f.Type, value) {
			apiErr := apierr.ErrInvalidInput.WithMsg("invalid field type").WithID(apierr.FieldTypeMismatch).WithCause(fmt.Errorf("field %q expects %s, got %T", key, f.Type, value))
			apierr.RecordDomainError(span, apiErr)
			return nil, apiErr
		}

		validated[key] = value
	}

	for _, f := range registerFields {
		if !f.Required {
			continue
		}

		if _, ok := validated[f.Key]; !ok {
			apiErr := apierr.ErrInvalidInput.WithMsg("missing required field").WithID(apierr.FieldRequiredMissing).WithCause(errors.New("missing field: " + f.Key))
			apierr.RecordDomainError(span, apiErr)
			return nil, apiErr
		}
	}

	metadata := make(map[string]any)
	schemaPayload := make(map[string]any)

	schemaPayload["schema_id"] = registerSchema.ID.String()
	schemaPayload["schema_version_id"] = registerVersion.ID.String()

	for k, v := range validated {
		schemaPayload[k] = v
	}

	flowMap := map[string]any{
		flowID: schemaPayload,
	}

	metadata[schemaType] = flowMap

	marshalledMetadata, err := json.Marshal(metadata)
	if err != nil {
		apiErr := apierr.ErrInternal.WithID(apierr.SystemInternalError).WithCause(err)
		apierr.RecordSystemError(span, apiErr)
		return nil, apiErr
	}

	rawMetadata := json.RawMessage(marshalledMetadata)
	return &rawMetadata, nil
}
