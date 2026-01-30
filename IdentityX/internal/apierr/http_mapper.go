package apierr

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/application/validation"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/permissions"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// FromService This function relies on the invariant that
// only apierr wraps errors. Do not use errors.As here.
func FromService(span trace.Span, err error) *Error {
	if err == nil {
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
	case schema.ErrRegisterSchemaNoPublishedVersion:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(ProjectUserRegisterOnSchemaNoVersion)
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
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(ProjectUserErrorEncodingMetadata).WithCause(e.Cause)
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
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(TokenCouldNotSign).WithCause(e.Cause)
		RecordSystemError(span, httpErr)
		return httpErr
	case inbounds.ErrParseProjectKey:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(ProjectFailedToParseKey).WithCause(e.Cause)
		RecordDomainError(span, httpErr)
		return httpErr
	case validation.ErrParseUUID:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(RequestValidationError).WithCause(e.Cause)
		RecordDomainError(span, httpErr)
		return httpErr
	case authz.ErrInvalidPrincipal:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(AuthInvalidPrincipal)
		RecordDomainError(span, httpErr)
		return httpErr
	case authz.ErrPrincipalMissingInContext:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(AuthPrincipalNotInContext)
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
	case inbounds.ErrGeneratingProjectKeys:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(ProjectErrorGeneratingKeys).WithCause(e.Cause)
		RecordSystemError(span, httpErr)
		return httpErr
	case inbounds.ErrParsingProjectPublicKey:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(ProjectErrorParsingKeys).WithCause(e.Cause)
		RecordSystemError(span, httpErr)
		return httpErr
	case inbounds.ErrNotProjectOwner:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(ProjectNotOwnedByPrincipal)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrFlowIDIsReserved:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(SchemaFlowIDIsReserved)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrFlowIDSchemaTypeConflict:
		httpErr := ErrConflict.WithMsg(e.Error()).WithID(SchemaFlowIDAlreadyExistsInType)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrSchemaNotOwned:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(SchemaNotOwnedByPrincipal)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrPublishSchemaPublished:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(SchemaTryingToPublishPublished)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrPublishSchemaArchived:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(SchemaTryingToPublishArchived)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrSchemaInvalidStatus:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(SchemaNoValidStatus)
		RecordSystemError(span, httpErr)
		return httpErr
	case inbounds.ErrSchemaNoPublishedVersions:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(SchemaNoPublishedVersion)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrSchemaOnlyDraft:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(SchemaHasOnlyDraftVersion)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrSchemaOnlyArchived:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(SchemaHasOnlyArchivedVersion)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrSchemaVersionMismatchLatest:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(SchemaVersionMismatch)
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
	case inbounds.ErrPublishVersionNoFields:
		httpErr := ErrBadRequest.WithMsg(e.Error()).WithID(SchemaVersionPublishWithNoFields)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrRevokeCurrentSession:
		httpErr := ErrForbidden.WithMsg(e.Error()).WithID(SessionSelfRevokeForbidden)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrSessionNotFound:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(SessionNotFound)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrSessionUnauthorized:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(SessionUnauthorized)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrInvalidIssuer:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenInvalidIssuer)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrTokenIDMismatch:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenMismatchDuringAuth)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrTokenSessionMismatch:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenSessionMismatch)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrAuthSessionRevoked:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(SessionRevoked)
		RecordSystemError(span, httpErr)
		return httpErr
	case inbounds.ErrEmptyCookie:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(RequestEmptyCookie)
		RecordSystemError(span, httpErr)
		return httpErr
	case inbounds.ErrTokenReuseNotAllowed:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenReuseIdentified)
		RecordDomainError(span, httpErr)
		return httpErr
	case outbounds.ErrServiceUnavailable:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(SystemServiceUnavailable)
		RecordDomainError(span, httpErr)
		return httpErr
	case outbounds.ErrRenderingEmail:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(SystemErrorRenderingEmail)
		RecordSystemError(span, httpErr)
		return httpErr
	case inbounds.ErrTXPanicked:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(DBTransactionPanicked)
		RecordSystemError(span, httpErr)
		logs.L().Error("Transaction panicked", zap.Any("panic", e.Panic))
		return httpErr
	case inbounds.ErrTokenUserMismatch:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenUserMismatch)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrUserAlreadyVerified:
		// Used to block resending verification email only
		httpErr := ErrForbidden.WithMsg(e.Error()).WithID(AuthAlreadyVerified)
		RecordDomainError(span, httpErr)
		return httpErr
	case auth.ErrTokenInvalidAlg:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenInvalidAlg)
		RecordDomainError(span, httpErr)
		return httpErr
	case auth.ErrTokenInvalidFormat:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenInvalidFormat)
		RecordDomainError(span, httpErr)
		return httpErr
	case auth.ErrTokenUntrusted:
		httpErr := ErrUnauthorized.WithMsg(e.Error()).WithID(TokenUntrusted)
		RecordDomainError(span, httpErr)
		return httpErr
	case inbounds.ErrFailedToRetrieveJWKS:
		httpErr := ErrInternal.WithMsg(e.Error()).WithID(SystemJWKSRetrievalFailed).WithCause(e.Cause)
		RecordSystemError(span, httpErr)
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
	case inbounds.ErrProjectUserNotFromProject:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(ProjectUserNotFromProject)
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
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(RequestValidationError).WithCause(e.Cause)
		return httpErr
	case ErrParsingNumber:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(RequestValidationError).WithCause(e.Cause)
		return httpErr
	case ErrMissingParam:
		httpErr := ErrInvalidInput.WithMsg(e.Error()).WithID(RequestValidationError)
		return httpErr
	default:
		httpErr := ErrInternal.WithMsg("unmapped handler error").WithCause(err).WithID(SystemInternalError)
		logs.L().Error("unmapped handler error", zap.Error(httpErr), zap.String("cause", err.Error()))
		return httpErr
	}
}
