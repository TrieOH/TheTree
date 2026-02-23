package errx

import (
	"github.com/MintzyG/fail/v3"
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
	RequestInvalidPassword           = fail.ID(0, "REQ", 3, true, "REQuestInvalidPassword")

	AuthInvalidAccessCookie = fail.ID(0, "AUTH", 0, false, "AUTHInvalidAccessCookie")
	AuthMissingAccessCookie = fail.ID(0, "AUTH", 1, false, "AUTHMissingAccessCookie")
	AuthSubjectNotInContext = fail.ID(0, "AUTH", 2, false, "AUTHSubjectNotInContext")
	AuthInvalidSubject      = fail.ID(0, "AUTH", 3, false, "AUTHInvalidSubject")

	AuthzInsufficientPermissions = fail.ID(0, "AUTHZ", 0, false, "AuthzInsufficientPermissions")

	TokenInvalidAccessClaims    = fail.ID(0, "TOKEN", 0, false, "TOKENInvalidAccessClaims")
	TokenMissingSubClaim        = fail.ID(0, "TOKEN", 1, false, "TOKENMissingSubClaim")
	TokenSubMarshalFailed       = fail.ID(0, "TOKEN", 2, false, "TOKENSubMarshalFailed")
	TokenSubUnmarshallingFailed = fail.ID(0, "TOKEN", 3, false, "TOKENSubUnmarshallingFailed")

	ValidationUUIDWasNil = fail.ID(1, "VAL", 0, false, "VALidationUUIDWasNil")

	EventSlugAlreadyInUse  = fail.ID(0, "EVENT", 0, false, "EventSlugAlreadyInUse")
	EventPublishNonDraft   = fail.ID(0, "EVENT", 1, false, "EventPublishNonDraft")
	EventCannotAddEditions = fail.ID(0, "EVENT", 2, false, "EventCannotAddEditions")

	EditionInvalidID        = fail.ID(0, "EDITION", 0, false, "EditionInvalidID")
	EditionValidationFailed = fail.ID(0, "EDITION", 1, false, "EditionValidationFailed")

	SQLResourceNotFound         = fail.ID(0, "SQL", 0, false, "SQLResourceNotFound")
	SQLInternalDBError          = fail.ID(9, "SQL", 1, false, "SQLInternalDBError")
	SQLForeignKeyViolation      = fail.ID(0, "SQL", 2, false, "SQLForeignKeyViolation")
	SQLSerializationFailure     = fail.ID(0, "SQL", 3, false, "SQLSerializationFailure")
	SQLNotNULLViolation         = fail.ID(0, "SQL", 4, false, "SQLNotNULLViolation")
	SQLValueTooLong             = fail.ID(0, "SQL", 5, false, "SQLValueTooLong")
	SQLDBConnectionError        = fail.ID(0, "SQL", 6, false, "SQLDBConnectionError")
	SQLDatabaseUnknownError     = fail.ID(0, "SQL", 7, false, "SQLDatabaseUnknownError")
	SQLUnmatchedUniqueViolation = fail.ID(1, "SQL", 8, false, "SQLUnmatchedUniqueViolation")
	SQLUnmatchedCheckViolation  = fail.ID(1, "SQL", 9, false, "SQLUnmatchedCheckViolation")

	SYSDependencyDown        = fail.ID(9, "SYS", 0, false, "SYStemDependencyDown")
	SYSServiceUnavailable    = fail.ID(9, "SYS", 1, false, "SYSServiceUnavailable")
	SYSUUIDV7GenerationError = fail.ID(9, "SYS", 2, false, "SYSUUIDV7GenerationError")

	SYSFunctionalityNotImplemented = fail.ID(9, "SYS", 0, true, "SYSFunctionalityNotImplemented")
	SYSTransactionNilContext       = fail.ID(9, "SYS", 1, true, "SYSTransactionNilContext")

	DBTransactionPanicked     = fail.ID(9, "DB", 0, false, "DBTransactionPanicked")
	DBBeginTransactionFailed  = fail.ID(9, "DB", 1, false, "DBBeginTransactionFailed")
	DBTransactionCommitFailed = fail.ID(9, "DB", 2, false, "DBTransactionCommitFailed")

	DBNestedTransactionNotAllowed = fail.ID(9, "DB", 0, true, "DBNestedTransactionNotAllowed")
)
