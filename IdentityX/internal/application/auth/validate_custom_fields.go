package auth

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

// validateAndConstructMetadata validates custom fields against a schema and returns structured metadata
func (uc *UseCase) validateAndConstructMetadata(
	ctx context.Context,
	span trace.Span,
	projectID uuid.UUID,
	schemaType schema.Type,
	flowID string,
	customFields *json.RawMessage,
) (*json.RawMessage, error) {
	var ok bool
	var err error
	var registerSchema *schema.Schema

	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	if registerSchema, err = schemas.FindByFlowIDAndType(ctx, flowID, schemaType, projectID); err != nil {
		return nil, err
	}
	if err = registerSchema.CanRegister(); err != nil {
		return nil, apierr.FromService(span, err)
	}
	var registerVersion *version.Version
	if registerVersion, err = versions.GetCurrent(ctx, registerSchema.ID); err != nil {
		return nil, err
	}
	if err = registerVersion.CanRegister(); err != nil {
		return nil, apierr.FromService(span, err)
	}
	if ok = registerSchema.IsVersion(registerVersion.ID); !ok {
		return nil, apierr.FromService(span, schema.ErrSchemaVersionMismatch{})
	}

	var registerFields []field.Field
	if registerFields, err = fields.GetByVersionID(ctx, registerVersion.ID); err != nil {
		return nil, err
	}

	fieldDefs := make(map[string]field.Field)
	for _, f := range registerFields {
		fieldDefs[f.Key] = f
	}

	var custom map[string]any
	if custom, err = parseCustomFields(customFields); err != nil {
		return nil, apierr.FromService(span, err)
	}

	validated, fieldsErr := validateFields(custom, fieldDefs, registerFields)
	if len(fieldsErr.FieldErrors) > 0 {
		return nil, apierr.FromService(span, fieldsErr)
	}

	metadata := make(map[string]any)
	schemaPayload := make(map[string]any)

	schemaPayload["schema_id"] = registerSchema.ID.String()
	schemaPayload["schema_version_id"] = registerVersion.ID.String()
	schemaPayload["fields"] = validated

	flowMap := map[string]any{
		flowID: schemaPayload,
	}

	metadata[string(schemaType)] = flowMap

	marshalledMetadata, err := json.Marshal(metadata)
	if err != nil {
		return nil, apierr.FromService(span, project_users.ErrEncodingProjectUserMetadata{Cause: err})
	}

	rawMetadata := json.RawMessage(marshalledMetadata)
	return &rawMetadata, nil
}

func validateFields(custom map[string]any, fieldDefs map[string]field.Field, registerFields []field.Field) (map[string]any, field.ErrFieldsValidation) {
	var ok bool
	var fieldErr field.ErrFieldsValidation
	validated := make(map[string]any)
	for key, value := range custom {
		var f field.Field
		f, ok = fieldDefs[key]
		if !ok {
			fieldErr.FieldErrors = append(fieldErr.FieldErrors, field.ErrFieldNotDefined{Key: key})
			continue
		}

		if !validateFieldValue(f.Type, value) {
			fieldErr.FieldErrors = append(fieldErr.FieldErrors, field.ErrInvalidFieldType{Key: key, Expected: string(f.Type), Got: value})
			continue
		}

		validated[key] = value
	}

	for _, f := range registerFields {
		if !f.Required {
			continue
		}

		if _, ok := validated[f.Key]; !ok {
			fieldErr.FieldErrors = append(fieldErr.FieldErrors, field.ErrMissingRequiredFields{Key: f.Key})
			continue
		}
	}
	return validated, fieldErr
}

func parseCustomFields(customFields *json.RawMessage) (custom map[string]any, err error) {
	if customFields == nil {
		return nil, apierr.ErrMissingCustomFields{}
	}
	if err = json.Unmarshal(*customFields, &custom); err != nil {
		return nil, apierr.ErrInvalidCustomFieldsJSON{Cause: err}
	}
	return custom, nil
}
