package apierr

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/application/validation"

	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/permissions"

	// "GoAuth/internal/domain/schema"
	// "GoAuth/internal/domain/version"
	"GoAuth/internal/ports/inbounds"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// FromService This function relies on the invariant that
// only apierr wraps errors. Do not use errors.As here.
func FromService(span trace.Span, err error) error {
	// var rs any
	if err == nil {
		return nil
	}
	defer func() {
		if err != nil {

		}
	}()
	switch e := err.(type) {
	case ErrInvalidCustomFieldsJSON:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(ID(RequestInvalidJSONFormat.String())).WithCause(e.Cause)
		RecordDomainError(span, httpErr)
		return httpErr
	case ErrMissingCustomFields:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(ID(RequestMissingSchemaCustomFields.String()))
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
	case validation.ErrParseUUID:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(ID(RequestValidationError.String())).WithCause(e.Cause)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrAddFieldsToNonDraftVersion:
		httpErr := ErrConflict.WithMsg(e.Error()).WithID(SchemaVersionNotDraft)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrInvalidFieldType:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(FieldInvalidType)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrInvalidFieldOwner:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(FieldInvalidOwner)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrDraftVersionOnNonPublished:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(SchemaVersionDraftOnNonPublished)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrPublishSchemaNonExistentVersionDraft:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(SchemaVersionDraftDoesntExist)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrPublishVersionPublished:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(SchemaVersionTryingToPublishPublished)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrPublishVersionArchived:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(SchemaVersionTryingToPublishArchived)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrPublishVersionInvalidStatus:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(SchemaVersionNoValidStatus)
		RecordSystemError(span, httpErr)
		return httpErr
	case inbounds.ErrPublishVersionNoChanges:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(SchemaVersionNoChanges)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrEmptyScopeName:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(ScopeEmptyName)
		RecordDomainError(span, httpErr)
		return httpErr
	case permissions.ErrInvalidPermissionObject:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(PermissionInvalidObject)
		RecordDomainError(span, httpErr)
		return httpErr
	case permissions.ErrInvalidPermissionAction:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(PermissionInvalidAction)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrRoleNotOwned:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(RoleNotOwnedByPrincipal)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrPermissionNotOwned:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(PermissionNotOwnedByPrincipal)
		RecordDomainError(span, httpErr)
		return httpErr
	case permissions.ErrActionMismatch:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(PermissionActionMismatch)
		RecordDomainError(span, httpErr)
		return httpErr
	case permissions.ErrObjectMismatch:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(PermissionObjectMismatch)
		RecordDomainError(span, httpErr)
		return httpErr
	case permissions.ErrInsufficientPermissions:
		httpErr := ErrForbidden.WithMsg("Permission Denied").WithID(PermissionInsufficient).WithCause(e)
		RecordDomainError(span, httpErr)
		return httpErr
	case permissions.ErrConditionValidationError:
		httpErr := ErrForbidden.WithMsg(e.Error()).WithID(PermissionConditionValidationError).WithCause(e)
		RecordDomainError(span, httpErr)
		return httpErr
	case permissions.ErrLogicalConditionValidationError:
		httpErr := ErrForbidden.WithMsg(e.Error()).WithID(PermissionConditionValidationError).WithCause(e)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrPublishVersionNotDraft:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(SchemaVersionNotDraft).WithCause(e)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrPublishNonExistentVersion:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(SchemaVersionTryingToPublishNonExistant).WithCause(e)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrFieldNotFound:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(FieldNotFound).WithCause(e)
		RecordDomainError(span, httpErr)
		return httpErr
	default:
		return err
		// httpErr := ErrInternal.WithMsg("unmapped service error").WithCause(err).WithID(SystemInternalError)
		// spanID := ""
		// if span != nil {
		// 	spanID = span.SpanContext().SpanID().String()
		// }
		// logs.L().Error("unmapped service error",
		// 	zap.Error(err),
		// 	zap.String("span", spanID),
		// 	zap.String("error_id", string(SystemInternalError)),
		// )
		// return httpErr
	}
}

func FromHandler(err error) *Error {
	if err == nil {
		logs.L().Warn("apierr.FromHandler called with a nil error")
		return nil
	}
	switch e := err.(type) {
	case validation.ErrParseUUID:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(ID(RequestValidationError.String())).WithCause(e.Cause)
		return httpErr
	case ErrParsingNumber:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(ID(RequestValidationError.String())).WithCause(e.Cause)
		return httpErr
	case ErrMissingParam:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(ID(RequestValidationError.String()))
		return httpErr
	default:
		httpErr := ErrInternal.WithMsg("unmapped handler error").WithCause(err).WithID(PlaceholderID)
		logs.L().Error("unmapped handler error", zap.Error(httpErr), zap.String("cause", err.Error()))
		return httpErr
	}
}
