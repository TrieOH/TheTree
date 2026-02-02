package apierr

import (
	"github.com/MintzyG/fail"
)

var (
	RequestMissingQueryParamValue = fail.ID(0, "REQ", 0, false, "REQuestMissingQueryParamValue")
	RequestMissingQueryParam      = fail.ID(0, "REQ", 1, false, "REQuestMissingQueryParam")
	// FIXME create tests for empty cookies
	RequestEmptyCookie             = fail.ID(0, "REQ", 2, false, "REQuestEmptyCookie")
	RequestUnknownQueryParam       = fail.ID(0, "REQ", 3, false, "REQuestUnknownQueryParam")
	RequestValidationError         = fail.ID(0, "REQ", 4, false, "REQuestValidationError")
	RequestParseUUIDError          = fail.ID(0, "REQ", 5, false, "REQuestParseUUIDError")
	RequestParseNumberError        = fail.ID(0, "REQ", 6, false, "REQuestParseNumberError")
	RequestMissingParamError       = fail.ID(0, "REQ", 7, false, "REQuestMissingParamError")
	RequestInvalidCustomFieldsJSON = fail.ID(0, "REQ", 8, false, "REQuestInvalidCustomFieldsJSON")

	RequestMissingSchemaCustomFields = fail.ID(0, "REQ", 0, true, "REQuestMissingSchemaCustomFields")
	RequestInvalidJSONFormat         = fail.ID(0, "REQ", 1, true, "REQuestInvalidJSONFormat")
	RequestNotApplicationJSON        = fail.ID(0, "REQ", 2, true, "REQuestNotApplicationJSON")
	RequestInvalidPassword           = fail.ID(0, "REQ", 5, true, "REQuestInvalidPassword") // FIXME there was a gap from 2 NIL 5, test extensively later

	AuthEmailAlreadyUsed     = fail.ID(0, "AUTH", 0, false, "AUTHEmailAlreadyUsed")
	AuthInvalidCredentials   = fail.ID(0, "AUTH", 1, false, "AUTHInvalidCredentials")
	AuthInvalidRefreshCookie = fail.ID(1, "AUTH", 2, false, "AUTHInvalidRefreshCookie")
	AuthInvalidAccessCookie  = fail.ID(1, "AUTH", 3, false, "AUTHInvalidAccessCookie")
	AuthMissingRefreshCookie = fail.ID(1, "AUTH", 4, false, "AUTHMissingRefreshCookie")
	AuthMissingAccessCookie  = fail.ID(1, "AUTH", 5, false, "AUTHMissingAccessCookie")

	AuthInvalidPrincipal      = fail.ID(1, "AUTH", 0, true, "AUTHInvalidPrincipal")
	AuthInvalidPassword       = fail.ID(0, "AUTH", 1, true, "AUTHInvalidPassword")
	AuthNotClient             = fail.ID(1, "AUTH", 2, true, "AUTHNotClient")
	AuthNotProjectUser        = fail.ID(1, "AUTH", 3, true, "AUTHNotProjectUser")
	AuthAlreadyVerified       = fail.ID(1, "AUTH", 4, true, "AUTHAlreadyVerified")
	AuthPrincipalNotInContext = fail.ID(1, "AUTH", 5, true, "AUTHPrincipalNotInContext")

	SessionRevoked             = fail.ID(1, "SESSION", 0, true, "SESSIONRevoked")
	SessionNotFound            = fail.ID(1, "SESSION", 1, true, "SESSIONNotFound")
	SessionSelfRevokeForbidden = fail.ID(1, "SESSION", 2, true, "SESSIONSelfRevokeForbidden")
	SessionUnauthorized        = fail.ID(1, "SESSION", 3, true, "SESSIONUnauthorized")

	TokenInvalid             = fail.ID(1, "TOKEN", 0, false, "TOKENInvalid")
	TokenExpired             = fail.ID(1, "TOKEN", 1, false, "TOKENExpired")
	TokenMalformed           = fail.ID(1, "TOKEN", 2, false, "TOKENMalformed")
	TokenSignatureInvalid    = fail.ID(1, "TOKEN", 3, false, "TOKENSignatureInvalid")
	TokenInvalidAlg          = fail.ID(1, "TOKEN", 4, false, "TOKENInvalidAlgorithm")
	TokenCouldNotSign        = fail.ID(1, "TOKEN", 5, false, "TOKENCouldNotSign")
	TokenInvalidAccessClaims = fail.ID(1, "TOKEN", 6, false, "TOKENInvalidAccessClaims")
	TokenNotYetValid         = fail.ID(1, "TOKEN", 7, false, "TOKENNotYetValid")
	TokenUsedBeforeIssued    = fail.ID(1, "TOKEN", 8, false, "TOKENUsedBeforeIssued")
	TokenInvalidIssuer       = fail.ID(1, "TOKEN", 9, false, "TOKENInvalidIssuer")
	TokenInvalidSubject      = fail.ID(1, "TOKEN", 10, false, "TOKENInvalidSubject")
	TokenInvalidAudience     = fail.ID(1, "TOKEN", 11, false, "TOKENInvalidAudience")
	TokenRefreshInvalidID    = fail.ID(1, "TOKEN", 12, false, "TOKENRefreshInvalidID")
	TokenAccessInvalidID     = fail.ID(1, "TOKEN", 13, false, "TOKENAccessInvalidID")
	TokenInvalidKid          = fail.ID(1, "TOKEN", 14, false, "TOKENInvalidKeyID")
	TokenUnknownKid          = fail.ID(1, "TOKEN", 15, false, "TOKENUnknownKeyID")
	TokenMissingKid          = fail.ID(1, "TOKEN", 16, false, "TOKENMissingKeyID")
	TokenUnverifiable        = fail.ID(1, "TOKEN", 17, false, "TOKENUnverifiable")
	TokenReuseIdentified     = fail.ID(1, "TOKEN", 18, false, "TOKENReuseIdentified")
	TokenUserMismatch        = fail.ID(1, "TOKEN", 19, false, "TOKENUserMismatch")
	TokenInvalidFormat       = fail.ID(1, "TOKEN", 20, false, "TOKENInvalidFormat")
	TokenUntrusted           = fail.ID(1, "TOKEN", 21, false, "TOKENUntrusted")

	TokenSessionMismatch      = fail.ID(1, "TOKEN", 0, true, "TOKENSessionMismatch")
	TokenMismatchDuringAuth   = fail.ID(1, "TOKEN", 1, true, "TokenMismatchDuringAuth")
	TokenMissingAccessClaims  = fail.ID(1, "TOKEN", 2, true, "TOKENMissingAccessClaims")
	TokenMissingRefreshClaims = fail.ID(1, "TOKEN", 3, true, "TOKENMissingRefreshClaims")

	ProjectErrorGeneratingKeys = fail.ID(1, "PROJECT", 0, false, "PROJECTErrorGeneratingKeys")
	ProjectNotOwnedByPrincipal = fail.ID(1, "PROJECT", 1, false, "PROJECTNotOwnedByPrincipal")

	ProjectNotFound = fail.ID(1, "PROJECT", 0, true, "PROJECTNotFound")

	ProjectUserErrorEncodingMetadata = fail.ID(1, "PROJECTUSER", 0, false, "PROJECTUSERErrorEncodingMetadata")

	ProjectUserRegisterOnSchemaVersionDraft    = fail.ID(1, "PROJECTUSER", 0, true, "PROJECTUSERRegisterOnSchemaVersionDraft")
	ProjectUserRegisterOnSchemaDraft           = fail.ID(1, "PROJECTUSER", 1, true, "PROJECTUSERRegisterOnSchemaDraft")
	ProjectUserRegisterOnSchemaArchived        = fail.ID(1, "PROJECTUSER", 2, true, "PROJECTUSERRegisterOnSchemaArchived")
	ProjectUserRegisterOnSchemaVersionArchived = fail.ID(1, "PROJECTUSER", 3, true, "PROJECTUSERRegisterOnSchemaVersionArchived")
	ProjectUserNotFromProject                  = fail.ID(1, "PROJECTUSER", 4, true, "PROJECTUSERNotFromProject")
	ProjectUserRegisterOnSchemaNoVersion       = fail.ID(1, "PROJECTUSER", 5, true, "PROJECTUSERRegisterOnSchemaNoVersion")
	ProjectUserRegisterOnNoneProject           = fail.ID(1, "PROJECTUSER", 6, true, "PROJECTUSERRegisterOnNoneProject")

	SchemaNotOwnedByPrincipal = fail.ID(1, "SCHEMA", 0, false, "SCHEMANotOwnedByPrincipal")
	SchemaNoValidStatus       = fail.ID(1, "SCHEMA", 1, false, "SCHEMANoValidStatus")
	SchemaInvalidFlowID       = fail.ID(1, "SCHEMA", 2, false, "SCHEMAInvalidFlowID")
	SchemaFlowIDIsReserved    = fail.ID(1, "SCHEMA", 3, false, "SCHEMAFlowIDIsReserved")

	SCHEMANoPublishedVersion        = fail.ID(1, "SCHEMA", 0, true, "SCHEMANoPublishedVersion")
	SchemaFlowIDAlreadyExistsInType = fail.ID(1, "SCHEMA", 1, true, "SCHEMAFlowIDAlreadyExistsInType")
	SchemaInvalidSchemaType         = fail.ID(1, "SCHEMA", 2, true, "SCHEMAInvalidSchemaType")
	SchemaHasOnlyDraftVersion       = fail.ID(1, "SCHEMA", 3, true, "SCHEMAHasOnlyDraftVersion")
	SchemaHasOnlyArchivedVersion    = fail.ID(1, "SCHEMA", 4, true, "SCHEMAHasOnlyArchivedVersion")
	SchemaTryingToPublishPublished  = fail.ID(1, "SCHEMA", 5, true, "SCHEMATryingToPublishPublished")
	SchemaTryingToPublishArchived   = fail.ID(1, "SCHEMA", 6, true, "SCHEMATryingToPublishArchived")
	SchemaMetadataNotAllowed        = fail.ID(1, "SCHEMA", 7, true, "SCHEMAMetadataNotAllowed")
	SchemaEmptySchemaType           = fail.ID(1, "SCHEMA", 8, true, "SCHEMAEmptySchemaType")
	SchemaEmptyFlowID               = fail.ID(1, "SCHEMA", 9, true, "SCHEMAEmptyFlowID")

	SchemaVersionNotDraft           = fail.ID(1, "SCHEMAVERSION", 0, false, "SCHEMAVERSIONNotDraft")
	SCHEMAVersionDraftAlreadyExists = fail.ID(1, "SCHEMAVERSION", 1, false, "SCHEMAVERSIONDraftAlreadyExists")

	SchemaVersionPublishWithNoFields         = fail.ID(1, "SCHEMAVERSION", 1, true, "SCHEMAVERSIONPublishWithNoFields")
	SchemaVersionDraftDoesntExist            = fail.ID(1, "SCHEMAVERSION", 2, true, "SCHEMAVERSIONDraftDoesntExist")
	SchemaVersionTryingToPublishPublished    = fail.ID(1, "SCHEMAVERSION", 3, true, "SCHEMAVERSIONTryingToPublishPublished")
	SchemaVersionTryingToPublishArchived     = fail.ID(1, "SCHEMAVERSION", 4, true, "SCHEMAVERSIONTryingToPublishArchived")
	SchemaVersionMismatch                    = fail.ID(1, "SCHEMAVERSION", 5, true, "SCHEMAVERSIONMismatch")
	SchemaVersionNonDraftAddFieldsNotAllowed = fail.ID(1, "SCHEMAVERSION", 6, true, "SCHEMAVERSIONNonDraftAddFieldsNotAllowed")
	SchemaVersionNoValidStatus               = fail.ID(1, "SCHEMAVERSION", 7, true, "SCHEMAVERSIONNoValidStatus")
	SchemaVersionDraftOnNonPublished         = fail.ID(1, "SCHEMAVERSION", 8, true, "SCHEMAVERSIONDraftOnNonPublished")
	SchemaVersionNoChanges                   = fail.ID(1, "SCHEMAVERSION", 9, true, "SCHEMAVERSIONNoChanges")
	SchemaVersionTryingToPublishNonExistant  = fail.ID(1, "SCHEMAVERSION", 10, true, "SCHEMAVERSIONTryingToPublishNonExistant")

	FIELDValidationErrorOnSchemaRegister = fail.ID(0, "FIELD", 0, false, "FIELDValidationErrorOnSchemaRegister")
	FIELDNotFound                        = fail.ID(0, "FIELD", 1, false, "FIELDNotFound")
	FIELDInvalidOwner                    = fail.ID(0, "FIELD", 2, false, "FIELDInvalidOwner")
	FieldNoAffectedRowsOnClone           = fail.ID(1, "FIELD", 3, false, "FIELDNoAffectedRowsOnClone")
	FIELDInvalidType                     = fail.ID(0, "FIELD", 4, false, "FIELDInvalidType")
	FIELDSameKeyForMultipleFields        = fail.ID(0, "FIELD", 5, false, "FIELDSameKeyForMultipleFields")
	FIELDSamePositionForMultipleFields   = fail.ID(0, "FIELD", 6, false, "FIELDSamePositionForMultipleFields")
	FIELDInvalidCharactersInKey          = fail.ID(0, "FIELD", 7, false, "FIELDInvalidCharactersInKey")

	ValidationUUIDWasNil = fail.ID(1, "VAL", 0, false, "VALidationUUIDWasNil")

	FORMMissingRequiredField = fail.ID(0, "FORM", 0, false, "FORMMissingRequiredFields")
	FORMInvalidFieldValue    = fail.ID(0, "FORM", 1, false, "FORMInvalidFieldValue")

	SQLNotFound                 = fail.ID(0, "SQL", 0, false, "SQLNotFound")
	SQLInternalDBError          = fail.ID(9, "SQL", 1, false, "SQLInternalDBError")
	SQLForeignKeyViolation      = fail.ID(0, "SQL", 2, false, "SQLForeignKeyViolation")
	SQLSerializationFailure     = fail.ID(0, "SQL", 3, false, "SQLSerializationFailure")
	SQLNotNULLViolation         = fail.ID(0, "SQL", 4, false, "SQLNotNULLViolation")
	SQLValueTooLong             = fail.ID(0, "SQL", 5, false, "SQLValueTooLong")
	SQLDBConnectionError        = fail.ID(0, "SQL", 6, false, "SQLDBConnectionError")
	SQLUnknownError             = fail.ID(0, "SQL", 7, false, "SQLUnknownError")
	SQLUnmatchedUniqueViolation = fail.ID(1, "SQL", 8, false, "SQLUnmatchedUniqueViolation")
	SQLUnmatchedCheckViolation  = fail.ID(1, "SQL", 9, false, "SQLUnmatchedCheckViolation")

	SYSDependencyDown        = fail.ID(9, "SYS", 0, false, "SYStemDependencyDown")
	SYSServiceUnavailable    = fail.ID(9, "SYS", 1, false, "SYSServiceUnavailable")
	SYSJWKSRetrievalFailed   = fail.ID(9, "SYS", 2, false, "SYSJWKSRetrievalFailed")
	SYSRenderingEmailFailed  = fail.ID(9, "SYS", 3, false, "SYSRenderingEmailFailed")
	SYSUUIDV7GenerationError = fail.ID(9, "SYS", 4, false, "SYSUUIDV7GenerationError")
	SYSJWKSEncodingFailed    = fail.ID(9, "SYS", 5, false, "SYSJWKSEncodingFailed")

	SYSFunctionalityNotImplemented = fail.ID(9, "SYS", 0, true, "SYSFunctionalityNotImplemented")
	SYSTransactionNilContext       = fail.ID(9, "SYS", 1, true, "SYSTransactionNilContext")

	DBTransactionPanicked     = fail.ID(9, "DB", 0, false, "DBTransactionPanicked")
	DBBeginTransactionFailed  = fail.ID(9, "DB", 1, false, "DBBeginTransactionFailed")
	DBTransactionCommitFailed = fail.ID(9, "DB", 2, false, "DBTransactionCommitFailed")

	DBNestedTransactionNotAllowed = fail.ID(9, "DB", 0, true, "DBNestedTransactionNotAllowed")

	ROLENotOwnedByPrincipal = fail.ID(0, "ROLE", 0, false, "ROLENotOwnedByPrincipal")
	ROLEAlreadyGranted      = fail.ID(0, "ROLE", 1, false, "ROLEAlreadyGranted")

	ROLENameAlreadyTaken = fail.ID(0, "ROLE", 0, true, "ROLENameAlreadyTaken")

	SCOPEDuplicateNameAndExternalID = fail.ID(0, "SCOPE", 0, false, "SCOPEDDuplicateNameAndExternalID")

	SCOPEInvalidShape = fail.ID(0, "SCOPE", 0, true, "SCOPEInvalidShape")
	SCOPEEmptyName    = fail.ID(0, "SCOPE", 1, true, "SCOPEEmptyName")

	PERMissionLogicalConditionValidationError = fail.ID(0, "PERM", 0, false, "PERMissionLogicalConditionValidationError")
	PERMissionConditionValidationError        = fail.ID(0, "PERM", 1, false, "PERMissionConditionValidationError")
	PERMissionActionMismatch                  = fail.ID(0, "PERM", 2, false, "PERMissionActionMismatch")
	PERMissionObjectMismatch                  = fail.ID(0, "PERM", 3, false, "PERMissionObjectMismatch")
	PERMissionAlreadyGranted                  = fail.ID(0, "PERM", 4, false, "PERMissionAlreadyGranted")
	PERMissionNotOwnedByPrincipal             = fail.ID(0, "PERM", 5, false, "PERMissionNotOwnedByPrincipal")
	PERMissionInvalidAction                   = fail.ID(0, "PERM", 6, false, "PERMissionInvalidAction")
	PERMissionInvalidObject                   = fail.ID(0, "PERM", 7, false, "PERMissionInvalidObject")
	PERMissionAlreadyExists                   = fail.ID(0, "PERM", 8, false, "PERMissionAlreadyExists")

	PERMissionInsufficient = fail.ID(0, "PERM", 0, true, "PERMissionInsufficient")

	EMAILTemplateNotFound = fail.ID(0, "EMAIL", 0, false, "EMAILTemplateNotFound")
)
