package apierr

import (
	"github.com/MintzyG/fail"
)

var (
	RequestMissingQueryParamValue = fail.ID(0, "REQ", 0, false, "REQuestMissingQueryParamValue")
	RequestMissingQueryParam      = fail.ID(0, "REQ", 1, false, "REQuestMissingQueryParam")
	// FIXME create tests for empty cookies
	RequestEmptyCookie = fail.ID(0, "REQ", 2, false, "REQuestEmptyCookie")

	RequestMissingSchemaCustomFields = fail.ID(0, "REQ", 0, true, "REQuestMissingSchemaCustomFields")
	RequestInvalidJSONFormat         = fail.ID(0, "REQ", 1, true, "REQuestInvalidJSONFormat")
	RequestValidationError           = fail.ID(0, "REQ", 2, true, "REQuestValidationError")
	RequestNotApplicationJSON        = fail.ID(0, "REQ", 3, true, "REQquestNotApplicationJSON")
	RequestUnknownQueryParam         = fail.ID(0, "REQ", 5, true, "REQuestUnknownQueryParam")

	AuthEmailAlreadyUsed     = fail.ID(0, "AUTH", 0, false, "AUTHEmailAlreadyUsed")
	AuthInvalidCredentials   = fail.ID(0, "AUTH", 1, false, "AUTHInvalidCredentials")
	AuthInvalidRefreshCookie = fail.ID(1, "AUTH", 2, false, "AUTHInvalidRefreshCookie")
	AuthInvalidAccessCookie  = fail.ID(1, "AUTH", 3, false, "AUTHInvalidAccessCookie")
	AuthMissingRefreshCookie = fail.ID(1, "AUTH", 4, false, "AUTHMissingRefreshCookie")
	AuthMissingAccessCookie  = fail.ID(1, "AUTH", 5, false, "AUTHMissingAccessCookie")

	AuthInvalidPrincipal = fail.ID(1, "AUTH", 0, true, "AUTHInvalidPrincipal")
	AuthInvalidPassword  = fail.ID(0, "AUTH", 1, true, "AUTHInvalidPassword")
	AuthNotClient        = fail.ID(1, "AUTH", 2, true, "AUTHNotClient")
	AuthNotProjectUser   = fail.ID(1, "AUTH", 3, true, "AUTHNotProjectUser")

	SCHEMANoPublishedVersion = fail.ID(0, "SCHEMA", 0, true, "SCHEMANoPublishedVersion")
)

const (
	AuthAlreadyVerified       ID = "AUTH_024"
	AuthPrincipalNotInContext ID = "AUTH_025"
)

const (
	SessionRevoked             ID = "SESS_003"
	SessionNotFound            ID = "SESS_004"
	SessionSelfRevokeForbidden ID = "SESS_008"
	SessionUnauthorized        ID = "SESS_011"
)

const (
	TokenInvalid              ID = "TOKEN_001"
	TokenExpired              ID = "TOKEN_002"
	TokenMalformed            ID = "TOKEN_003"
	TokenSignatureInvalid     ID = "TOKEN_004"
	TokenInvalidAlg           ID = "TOKEN_005"
	TokenCouldNotSign         ID = "TOKEN_007"
	TokenInvalidAccessClaims  ID = "TOKEN_010"
	TokenNotYetValid          ID = "TOKEN_011"
	TokenUsedBeforeIssued     ID = "TOKEN_012"
	TokenInvalidIssuer        ID = "TOKEN_013"
	TokenInvalidSubject       ID = "TOKEN_014"
	TokenInvalidAudience      ID = "TOKEN_015"
	TokenRefreshInvalidID     ID = "TOKEN_016"
	TokenSessionMismatch      ID = "TOKEN_020"
	TokenMismatchDuringAuth   ID = "TOKEN_021"
	TokenAccessInvalidID      ID = "TOKEN_022"
	TokenInvalidKid           ID = "TOKEN_025"
	TokenUnknownKid           ID = "TOKEN_026"
	TokenMissingKid           ID = "TOKEN_027"
	TokenUnverifiable         ID = "TOKEN_029"
	TokenMissingAccessClaims  ID = "TOKEN_030"
	TokenMissingRefreshClaims ID = "TOKEN_031"
	TokenReuseIdentified      ID = "TOKEN_032"
	TokenUserMismatch         ID = "TOKEN_033"
	TokenInvalidFormat        ID = "TOKEN_034"
	TokenUntrusted            ID = "TOKEN_035"
)

const (
	ProjectNotFound            ID = "PROJ_002"
	ProjectErrorGeneratingKeys ID = "PROJ_004"
	ProjectErrorParsingKeys    ID = "PROJ_005"
	ProjectNotOwnedByPrincipal ID = "PROJ_007"
	ProjectFailedToParseKey    ID = "PROJ_009"
)

const (
	ProjectUserRegisterOnSchemaVersionDraft    ID = "PROJ_USR_002"
	ProjectUserRegisterOnSchemaDraft           ID = "PROJ_USR_003"
	ProjectUserRegisterOnSchemaArchived        ID = "PROJ_USR_004"
	ProjectUserRegisterOnSchemaVersionArchived ID = "PROJ_USR_005"
	ProjectUserErrorEncodingMetadata           ID = "PROJ_USR_006"
	ProjectUserNotFromProject                  ID = "PROJ_USR_007"
	ProjectUserRegisterOnSchemaNoVersion       ID = "PROJ_USR_008"
)

const (
	SchemaFlowIDAlreadyExistsInType ID = "SCHEMA_001"
	SchemaInvalidSchemaType         ID = "SCHEMA_002"
	SchemaNotOwnedByPrincipal       ID = "SCHEMA_004"
	SchemaHasOnlyDraftVersion       ID = "SCHEMA_006"
	SchemaHasOnlyArchivedVersion    ID = "SCHEMA_007"
	SchemaTryingToPublishPublished  ID = "SCHEMA_008"
	SchemaTryingToPublishArchived   ID = "SCHEMA_009"
	SchemaNoValidStatus             ID = "SCHEMA_010"
	SchemaInvalidFlowID             ID = "SCHEMA_011"
	SchemaFlowIDIsReserved          ID = "SCHEMA_012"
	SchemaMetadataNotAllowed        ID = "SCHEMA_014"
	SchemaEmptySchemaType           ID = "SCHEMA_016"
	SchemaEmptyFlowID               ID = "SCHEMA_017"
)

const (
	SchemaVersionDraftAlreadyExists         ID = "SCM_VER_001"
	SchemaVersionPublishWithNoFields        ID = "SCM_VER_003"
	SchemaVersionDraftDoesntExist           ID = "SCM_VER_004"
	SchemaVersionTryingToPublishPublished   ID = "SCM_VER_005"
	SchemaVersionTryingToPublishArchived    ID = "SCM_VER_006"
	SchemaVersionMismatch                   ID = "SCM_VER_007"
	SchemaVersionNotDraft                   ID = "SCM_VER_008"
	SchemaVersionNoValidStatus              ID = "SCM_VER_009"
	SchemaVersionDraftOnNonPublished        ID = "SCM_VER_010"
	SchemaVersionNoChanges                  ID = "SCM_VER_011"
	SchemaVersionTryingToPublishNonExistant ID = "SCM_VER_012"
)

const (
	FieldSamePositionForMultipleFields ID = "FIELD_001"
	FieldNoAffectedRowsOnClone         ID = "FIELD_002"
	FieldInvalidCharactersInKey        ID = "FIELD_003"
	FieldSameKeyForMultipleFields      ID = "FIELD_006"
	FieldInvalidType                   ID = "FIELD_008"
	FieldInvalidOwner                  ID = "FIELD_009"
	FieldValidationErrSchemaRegister   ID = "FIELD_010"
	FieldNotFound                      ID = "FIELD_012"
)

const (
	ScopeEmptyName                  ID = "SCP_001"
	ScopeDuplicateNameAndExternalID ID = "SCP_002"
	ScopeInvalid                    ID = "SCP_003"
)

const (
	PermissionInvalidObject            ID = "PERM_001"
	PermissionInvalidAction            ID = "PERM_002"
	PermissionNotOwnedByPrincipal      ID = "PERM_003"
	PermissionAlreadyGranted           ID = "PERM_004"
	PermissionObjectMismatch           ID = "PERM_005"
	PermissionActionMismatch           ID = "PERM_006"
	PermissionInsufficient             ID = "PERM_007"
	PermissionConditionValidationError ID = "PERM_008"
)

const (
	RoleNotOwnedByPrincipal ID = "ROLE_002"
	RoleNameTaken           ID = "ROLE_003"
	RoleAlreadyGranted      ID = "ROLE_004"
)

const (
	DBNotFound             ID = "DB_000"
	DBUniqueViolation      ID = "DB_001"
	DBForeignKeyViolation  ID = "DB_002"
	DBNotNullViolation     ID = "DB_003"
	DBValueTooLong         ID = "DB_004"
	DBSerializationFailure ID = "DB_005"
	DBCommitTXFailed       ID = "DB_006"
	DBBeginTXFailed        ID = "DB_007"
	DBNestedTXNotAllowed   ID = "DB_008"
	DBCheckViolation       ID = "DB_009"
	DBTransactionPanicked  ID = "DB_010"
)

const (
	SystemInternalError            ID = "SYS_001"
	SystemDependencyDown           ID = "SYS_003"
	SystemNotImplemented           ID = "SYS_004"
	SystemTransactionWithNoContext ID = "SYS_005"
	SystemErrorGeneratingUUID      ID = "SYS_006"
	SystemErrorBCryptHashingFailed ID = "SYS_007"
	SystemServiceUnavailable       ID = "SYS_008"
	SystemErrorRenderingEmail      ID = "SYS_009"
	SystemJWKSRetrievalFailed      ID = "SYS_010"
)

var (
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

	SCHEMAVersionDraftAlreadyExists = fail.ID(0, "SCHEMA", 0, false, "SCHEMAVersionDraftAlreadyExists")

	FIELDSamePositionForMultipleFields = fail.ID(0, "FIELD", 0, false, "FIELDSamePositionForMultipleFields")
	FIELDSameKeyForMultipleFields      = fail.ID(0, "FIELD", 1, false, "FIELDSameKeyForMultipleFields")
	FIELDInvalidCharactersInKey        = fail.ID(0, "FIELD", 2, false, "FIELDInvalidCharactersInKey")

	SCOPEDuplicateNameAndExternalID = fail.ID(0, "SCOPE", 0, false, "SCOPEDDuplicateNameAndExternalID")
	SCOPEInvalid                    = fail.ID(0, "SCOPE", 1, false, "SCOPEInvalid")

	ROLENameAlreadyTaken = fail.ID(0, "ROLE", 0, false, "ROLENameAlreadyTaken")
	ROLEAlreadyGranted   = fail.ID(0, "ROLE", 1, false, "ROLEAlreadyGranted")

	PERMISSIONAlreadyGranted = fail.ID(0, "PERM", 0, false, "PERMISSIONAlreadyGranted")
)
