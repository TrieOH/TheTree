package apierr

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// FromService This function relies on the invariant that
// only apierr wraps errors. Do not use errors.As here.
func FromService(span trace.Span, err error) *Error {
	switch e := err.(type) {
	case version.ErrRegisterOnVersionDraft:
		httpErr := ErrBadRequest.WithMsg("can't register to a draft schema version").WithID(ProjectUserRegisterOnSchemaVersionDraft)
		RecordDomainError(span, httpErr)
		return httpErr
	case version.ErrRegisterOnVersionArchive:
		httpErr := ErrBadRequest.WithMsg("can't register to an archived schema version").WithID(ProjectUserRegisterOnSchemaVersionArchived)
		RecordDomainError(span, httpErr)
		return httpErr
	case schema.ErrRegisterOnSchemaDraft:
		httpErr := ErrBadRequest.WithMsg("can't register to a draft schema").WithID(ProjectUserRegisterOnSchemaDraft)
		RecordDomainError(span, httpErr)
		return httpErr
	case schema.ErrRegisterOnSchemaArchive:
		httpErr := ErrBadRequest.WithMsg("can't register to an archived schema").WithID(ProjectUserRegisterOnSchemaArchived)
		RecordDomainError(span, httpErr)
		return httpErr
	case schema.ErrSchemaNoPublishedVersion:
		httpErr := ErrBadRequest.WithMsg("schema has no published version").WithID(SchemaNoPublishedVersion)
		RecordDomainError(span, httpErr)
		return httpErr
	case schema.ErrSchemaVersionMismatch:
		httpErr := ErrInternal.WithMsg("schema version and retrieved version mismatch").WithID(SchemaVersionMismatch)
		RecordSystemError(span, httpErr)
		return httpErr
	case ErrInvalidCustomFieldsJSON:
		httpErr := ErrInvalidInput.WithMsg("invalid custom fields JSON").WithID(RequestInvalidJSONFormat).WithCause(e.Cause)
		RecordDomainError(span, httpErr)
		return httpErr
	case ErrMissingCustomFields:
		httpErr := ErrInvalidInput.WithMsg("schema custom fields are required on a schema register").WithID(RequestMissingSchemaCustomFields)
		RecordDomainError(span, httpErr)
		return httpErr
	// FIXME Add a WithMeta to apierr and FUN
	case field.ErrFieldsValidation:
		httpErr := ErrInvalidInput.WithMsg("error validating fields for schema register").WithID(FieldValidationErrSchemaRegister)
		for _, subE := range e.FieldErrors {
			httpErr = httpErr.WithCause(subE)
		}
		RecordDomainError(span, httpErr)
		return httpErr
	case project_users.ErrEncodingProjectUserMetadata:
		httpErr := ErrInternal.WithMsg("error encoding project user metadata").WithID(SystemInternalError).WithCause(e.Cause)
		RecordSystemError(span, httpErr)
		return httpErr
	default:
		httpErr := ErrInternal.WithMsg("unmapped service error").WithCause(err).WithID(SystemInternalError)

		logs.L().Error("unmapped service error",
			zap.Error(err),
			zap.String("span", span.SpanContext().SpanID().String()),
			zap.String("error_id", string(SystemInternalError)),
		)
		return httpErr
	}
}
