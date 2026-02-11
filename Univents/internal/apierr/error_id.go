package apierr

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

	AuthEmailAlreadyUsed     = fail.ID(0, "AUTH", 0, false, "AUTHEmailAlreadyUsed")
	AuthInvalidCredentials   = fail.ID(0, "AUTH", 1, false, "AUTHInvalidCredentials")
	AuthInvalidRefreshCookie = fail.ID(0, "AUTH", 2, false, "AUTHInvalidRefreshCookie")
	AuthInvalidAccessCookie  = fail.ID(0, "AUTH", 3, false, "AUTHInvalidAccessCookie")
	AuthMissingRefreshCookie = fail.ID(0, "AUTH", 4, false, "AUTHMissingRefreshCookie")
	AuthMissingAccessCookie  = fail.ID(0, "AUTH", 5, false, "AUTHMissingAccessCookie")

	AuthInvalidPrincipal      = fail.ID(0, "AUTH", 0, true, "AUTHInvalidPrincipal")
	AuthInvalidPassword       = fail.ID(0, "AUTH", 1, true, "AUTHInvalidPassword")
	AuthNotClient             = fail.ID(0, "AUTH", 2, true, "AUTHNotClient")
	AuthNotProjectUser        = fail.ID(0, "AUTH", 3, true, "AUTHNotProjectUser")
	AuthAlreadyVerified       = fail.ID(0, "AUTH", 4, true, "AUTHAlreadyVerified")
	AuthPrincipalNotInContext = fail.ID(0, "AUTH", 5, true, "AUTHPrincipalNotInContext")
	AuthUserSchemaOutdated    = fail.ID(0, "AUTH", 6, true, "AUTHUserSchemaOutdated")
	AuthTokenAlreadyUsed      = fail.ID(0, "AUTH", 7, true, "AUTHTokenAlreadyUsed")
	AuthApiKeyNotAllowed      = fail.ID(0, "AUTH", 8, true, "AUTHApiKeyNotAllowed")
	AuthInvalidApiKey         = fail.ID(0, "AUTH", 9, true, "AUTHInvalidApiKey")
	AuthInvalidApiKeyShape    = fail.ID(0, "AUTH", 10, true, "AUTHInvalidApiKeyShape")

	ValidationUUIDWasNil = fail.ID(1, "VAL", 0, false, "VALidationUUIDWasNil")

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
	SYSUUIDV7GenerationError = fail.ID(9, "SYS", 4, false, "SYSUUIDV7GenerationError")
	SYSCryptoError           = fail.ID(9, "SYS", 6, false, "SYSCryptoError")

	SYSFunctionalityNotImplemented = fail.ID(9, "SYS", 0, true, "SYSFunctionalityNotImplemented")
	SYSTransactionNilContext       = fail.ID(9, "SYS", 1, true, "SYSTransactionNilContext")

	DBTransactionPanicked     = fail.ID(9, "DB", 0, false, "DBTransactionPanicked")
	DBBeginTransactionFailed  = fail.ID(9, "DB", 1, false, "DBBeginTransactionFailed")
	DBTransactionCommitFailed = fail.ID(9, "DB", 2, false, "DBTransactionCommitFailed")

	DBNestedTransactionNotAllowed = fail.ID(9, "DB", 0, true, "DBNestedTransactionNotAllowed")
)
