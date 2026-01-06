package apierr

const (
	AuthInvalidEmail         ID = "AUTH_001"
	AuthInvalidPassword      ID = "AUTH_002"
	AuthWrongPassword        ID = "AUTH_003"
	AuthUserNotFound         ID = "AUTH_004"
	AuthEmailAlreadyUsed     ID = "AUTH_005"
	AuthAccountDisabled      ID = "AUTH_006"
	AuthTokenInvalid         ID = "AUTH_007"
	AuthTokenExpired         ID = "AUTH_008"
	AuthRefreshInvalid       ID = "AUTH_009"
	AuthRefreshRevoked       ID = "AUTH_010"
	AuthMissingAccessClaims  ID = "AUTH_011"
	AuthInvalidAccessClaims  ID = "AUTH_012"
	AuthMissingRefreshClaims ID = "AUTH_013"
	AuthInvalidRefreshClaims ID = "AUTH_014"
	AuthInvalidCredentials   ID = "AUTH_015"
	AuthInvalidRefreshCookie ID = "AUTH_016"
	AuthInvalidAccessCookie  ID = "AUTH_017"
	AuthMissingRefreshCookie ID = "AUTH_018"
	AuthMissingAccessCookie  ID = "AUTH_019"
	AuthMissingPrincipal     ID = "AUTH_020"
	AuthNotClient            ID = "AUTH_021"
	AuthInvalidPrincipal     ID = "AUTH_022"
)

const (
	UserRequiredID         ID = "USER_001"
	UserNotFound           ID = "USER_002"
	UserAlreadyExists      ID = "USER_003"
	UserEmailNotVerified   ID = "USER_004"
	UserPasswordTooWeak    ID = "USER_005"
	UserInvalidID          ID = "USER_006"
	UserEmailAlreadyExists ID = "USER_003"
)

const (
	SessionInvalidID           ID = "SESS_001"
	SessionExpired             ID = "SESS_002"
	SessionRevoked             ID = "SESS_003"
	SessionNotFound            ID = "SESS_004"
	SessionAlreadyActive       ID = "SESS_005"
	SessionLimitReached        ID = "SESS_006"
	SessionRequiredID          ID = "SESS_007"
	SessionSelfRevokeForbidden ID = "SESS_008"
	SessionUpdateFailed        ID = "SESS_009"
)

const (
	TokenInvalid          ID = "TOKEN_001"
	TokenExpired          ID = "TOKEN_002"
	TokenMalformed        ID = "TOKEN_003"
	TokenSignatureInvalid ID = "TOKEN_004"
	TokenUnsupportedAlg   ID = "TOKEN_005"
	TokenRequiredID       ID = "TOKEN_006"
	TokenCouldNotSign     ID = "TOKEN_007"
	TokenMissingRefresh   ID = "TOKEN_008"
	TokenMissingAccess    ID = "TOKEN_009"
	TokenInvalidClaims    ID = "TOKEN_0010"
	TokenNotYetValid      ID = "TOKEN_0011"
	TokenUsedBeforeIssued ID = "TOKEN_0012"
	TokenInvalidIssuer    ID = "TOKEN_0013"
	TokenInvalidSubject   ID = "TOKEN_0014"
	TokenInvalidAudience  ID = "TOKEN_0015"
	TokenInvalidID        ID = "TOKEN_0016"
	TokenRevoked          ID = "TOKEN_0017"
	TokenRevokeFailed     ID = "TOKEN_0018"
)

const (
	ProjectInvalidID            ID = "PROJ_001"
	ProjectNotFound             ID = "PROJ_002"
	ProjectInactive             ID = "PROJ_003"
	ProjectErrorGeneratingKeys  ID = "PROJ_004"
	ProjectErrorParsingKeys     ID = "PROJ_005"
	ProjectOwnershipCheckFailed ID = "PROJ_006"
	ProjectNotOwnedByPrincipal  ID = "PROJ_007"
)

const (
	ProjectUserInvalidMetadata ID = "PROJ_USR_001"
)

const (
	SchemaFlowIDAlreadyExistsInType ID = "SCHEMA_001"
	SchemaInvalidSchemaType         ID = "SCHEMA_002"
	SchemaInvalidID                 ID = "SCHEMA_003"
	SchemaNotOwnedByPrincipal       ID = "SCHEMA_004"
)

const (
	SchemaVersionDraftAlreadyExists ID = "SCM_VER_001"
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
)

const (
	SystemInternalError  ID = "SYS_001"
	SystemMisconfigured  ID = "SYS_002"
	SystemDependencyDown ID = "SYS_003"
	SystemUnimplemented  ID = "SYS_004"
)
