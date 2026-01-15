package apierr

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/application/validation"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/utils"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// FromService This function relies on the invariant that
// only apierr wraps errors. Do not use errors.As here.
func FromService(span trace.Span, err error) *Error {
	if err == nil {
		logs.L().Warn("apierr.FromService called with a nil error")
		return nil
	}
	switch e := err.(type) {
	case version.ErrRegisterOnVersionDraft:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(ProjectUserRegisterOnSchemaVersionDraft)
		RecordDomainError(span, httpErr)
		return httpErr
	case version.ErrRegisterOnVersionArchive:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(ProjectUserRegisterOnSchemaVersionArchived)
		RecordDomainError(span, httpErr)
		return httpErr
	case schema.ErrRegisterOnSchemaDraft:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(ProjectUserRegisterOnSchemaDraft)
		RecordDomainError(span, httpErr)
		return httpErr
	case schema.ErrRegisterOnSchemaArchive:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(ProjectUserRegisterOnSchemaArchived)
		RecordDomainError(span, httpErr)
		return httpErr
	case schema.ErrSchemaNoPublishedVersion:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(SchemaNoPublishedVersion)
		RecordDomainError(span, httpErr)
		return httpErr
	case schema.ErrSchemaVersionMismatch:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(SchemaVersionMismatch)
		RecordSystemError(span, httpErr)
		return httpErr
	case ErrInvalidCustomFieldsJSON:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(RequestInvalidJSONFormat).WithCause(e.Cause)
		RecordDomainError(span, httpErr)
		return httpErr
	case ErrMissingCustomFields:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(RequestMissingSchemaCustomFields)
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
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(SystemInternalError).WithCause(e.Cause)
		RecordSystemError(span, httpErr)
		return httpErr
	case auth.ErrTokenMissingKID:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenMissingKid)
		RecordDomainError(span, httpErr)
		return httpErr
	case auth.ErrInvalidToken:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenInvalid)
		RecordDomainError(span, httpErr)
		return httpErr
	case auth.ErrTokenInvalidKID:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenInvalidKid)
		RecordDomainError(span, httpErr)
		return httpErr
	case auth.ErrTokenUnknownKID:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenUnknownKid)
		RecordDomainError(span, httpErr)
		return httpErr
	case auth.ErrSigningToken:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenCouldNotSign).WithCause(e.Cause)
		RecordSystemError(span, httpErr)
		return httpErr
	case utils.ErrParseProjectKey:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(ProjectFailedToParseKey).WithCause(e.Cause)
		RecordDomainError(span, httpErr)
		return httpErr
	case validation.ErrParseUUID:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(RequestValidationError).WithCause(e.Cause)
		RecordDomainError(span, httpErr)
		return httpErr
	case authz.ErrMissingPrincipal:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(AuthMissingPrincipal)
		RecordDomainError(span, httpErr)
		return httpErr
	case authz.ErrPrincipalMissingInContext:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(AuthMissingPrincipal)
		RecordDomainError(span, httpErr)
		return httpErr
	case authz.ErrMissingAccessClaims:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenMissingAccessClaims)
		RecordDomainError(span, httpErr)
		return httpErr
	case authz.ErrMissingRefreshClaims:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenMissingRefreshClaims)
		RecordDomainError(span, httpErr)
		return httpErr
	case authz.ErrInvalidAccessJTI:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(TokenAccessInvalidID).WithCause(e.Cause)
		RecordDomainError(span, httpErr)
		return httpErr
	case authz.ErrInvalidRefreshJTI:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(TokenRefreshInvalidID).WithCause(e.Cause)
		RecordDomainError(span, httpErr)
		return httpErr
	case ErrPasswordTooLong:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(AuthInvalidPassword)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrHashingPassword:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(SystemErrorBCryptHashingFailed).WithCause(e.Cause)
		RecordSystemError(span, httpErr)
		return httpErr
	case inbounds.ErrEmailAlreadyInUse:
		httpErr := ErrConflict.WithMsg(e.Error()).WithID(AuthEmailAlreadyUsed).WithCause(e.Cause)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrInvalidCredentials:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(AuthInvalidCredentials)
		if e.Cause != nil {
			httpErr = httpErr.WithCause(e.Cause)
		}
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrGeneratingUUID:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(SystemErrorGeneratingUUID).WithCause(e.Cause)
		RecordSystemError(span, httpErr)
		return httpErr
	case inbounds.ErrTokenInvalid:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenInvalid)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrEmptyFlowID:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(SchemaEmptyFlowID)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrEmptySchemaType:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(SchemaEmptySchemaType)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrInvalidSchemaType:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(SchemaInvalidSchemaType)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrInvalidFlowID:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(SchemaInvalidFlowID)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrCustomFieldsNotAllowed:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(SchemaMetadataNotAllowed)
		RecordDomainError(span, httpErr)
		return httpErr
	default:
		httpErr := ErrInternal.WithMsg("unmapped service error").WithCause(err).WithID(SystemInternalError)
		spanID := ""
		if span != nil {
			spanID = span.SpanContext().SpanID().String()
		}
		logs.L().Error("unmapped service error",
			zap.Error(err),
			zap.String("span", spanID),
			zap.String("error_id", string(SystemInternalError)),
		)
		return httpErr
	}
}

func FromHandler(err error) *Error {
	if err == nil {
		logs.L().Warn("apierr.FromHandler called with a nil error")
		return nil
	}
	switch e := err.(type) {
	case validation.ErrParseUUID:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(RequestValidationError)
		return httpErr
	default:
		httpErr := ErrInternal.WithMsg("unmapped service error").WithCause(err).WithID(SystemInternalError)
		return httpErr
	}
}
