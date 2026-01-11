package apierr

const (
	RequestMissingQueryParamValue    ID = "REQ_001"
	RequestMissingQueryParam         ID = "REQ_002"
	RequestMissingSchemaCustomFields ID = "REQ_003"
	RequestInvalidJSON               ID = "REQ_004"
	RequestValidationError           ID = "REQ_005"
)

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
	AuthNotProjectUser       ID = "AUTH_023"
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
	SessionMissingID           ID = "SESS_010"
	SessionUnauthorized        ID = "SESS_011"
)

const (
	TokenInvalid              ID = "TOKEN_001"
	TokenExpired              ID = "TOKEN_002"
	TokenMalformed            ID = "TOKEN_003"
	TokenSignatureInvalid     ID = "TOKEN_004"
	TokenUnsupportedAlg       ID = "TOKEN_005"
	TokenRequiredID           ID = "TOKEN_006"
	TokenCouldNotSign         ID = "TOKEN_007"
	TokenMissingRefresh       ID = "TOKEN_008"
	TokenMissingAccess        ID = "TOKEN_009"
	TokenInvalidAccessClaims  ID = "TOKEN_010"
	TokenNotYetValid          ID = "TOKEN_011"
	TokenUsedBeforeIssued     ID = "TOKEN_012"
	TokenInvalidIssuer        ID = "TOKEN_013"
	TokenInvalidSubject       ID = "TOKEN_014"
	TokenInvalidAudience      ID = "TOKEN_015"
	TokenRefreshInvalidID     ID = "TOKEN_016"
	TokenRevoked              ID = "TOKEN_017"
	TokenRevokeFailed         ID = "TOKEN_018"
	TokenRefreshIDMissing     ID = "TOKEN_019"
	TokenSessionMismatch      ID = "TOKEN_020"
	TokenMismatchDuringAuth   ID = "TOKEN_021"
	TokenAccessInvalidID      ID = "TOKEN_022"
	TokenAccessIDMissing      ID = "TOKEN_023"
	TokenAccessIDMatched      ID = "TOKEN_024"
	TokenInvalidKid           ID = "TOKEN_025"
	TokenUnknownKid           ID = "TOKEN_026"
	TokenMissingKid           ID = "TOKEN_027"
	TokenInvalidRefreshClaims ID = "TOKEN_028"
	TokenUnverifiable         ID = "TOKEN_029"
)

const (
	ProjectInvalidID            ID = "PROJ_001"
	ProjectNotFound             ID = "PROJ_002"
	ProjectInactive             ID = "PROJ_003"
	ProjectErrorGeneratingKeys  ID = "PROJ_004"
	ProjectErrorParsingKeys     ID = "PROJ_005"
	ProjectOwnershipCheckFailed ID = "PROJ_006"
	ProjectNotOwnedByPrincipal  ID = "PROJ_007"
	ProjectMissingID            ID = "PROJ_008"
	ProjectFailedToParseKey     ID = "PROJ_009"
)

const (
	ProjectUserInvalidMetadata                 ID = "PROJ_USR_001"
	ProjectUserRegisterOnSchemaVersionDraft    ID = "PROJ_USR_002"
	ProjectUserRegisterOnSchemaDraft           ID = "PROJ_USR_003"
	ProjectUserRegisterOnSchemaArchived        ID = "PROJ_USR_004"
	ProjectUserRegisterOnSchemaVersionArchived ID = "PROJ_USR_005"
)

const (
	SchemaFlowIDAlreadyExistsInType ID = "SCHEMA_001"
	SchemaInvalidSchemaType         ID = "SCHEMA_002"
	SchemaInvalidID                 ID = "SCHEMA_003"
	SchemaNotOwnedByPrincipal       ID = "SCHEMA_004"
	SchemaNoPublishedVersion        ID = "SCHEMA_005"
	SchemaHasOnlyDraftVersion       ID = "SCHEMA_006"
	SchemaHasOnlyArchivedVersion    ID = "SCHEMA_007"
	SchemaTryingToPublishPublished  ID = "SCHEMA_008"
	SchemaTryingToPublishArchived   ID = "SCHEMA_009"
	SchemaNoValidType               ID = "SCHEMA_010"
	SchemaInvalidFlowID             ID = "SCHEMA_011"
	SchemaFlowIDIsReserved          ID = "SCHEMA_012"
	SchemaInvalidMetadata           ID = "SCHEMA_013"
	SchemaMetadataNotAllowed        ID = "SCHEMA_014"
	SchemaMissingID                 ID = "SCHEMA_015"
)

const (
	SchemaVersionDraftAlreadyExists       ID = "SCM_VER_001"
	SchemaVersionInvalidID                ID = "SCM_VER_002"
	SchemaVersionPublishWithNoFields      ID = "SCM_VER_003"
	SchemaVersionDraftDoesntExist         ID = "SCM_VER_004"
	SchemaVersionTryingToPublishPublished ID = "SCM_VER_005"
	SchemaVersionTryingToPublishArchived  ID = "SCM_VER_006"
	SchemaVersionMismatch                 ID = "SCM_VER_007"
	SchemaVersionNotDraft                 ID = "SCM_VER_008"
	SchemaVersionNoValidType              ID = "SCM_VER_009"
	SchemaVersionDraftOnNonPublished      ID = "SCM_VER_010"
	SchemaVersionNoChanges                ID = "SCM_VER_011"
)

const (
	FieldSamePositionForMultipleFields ID = "FIELD_001"
	FieldNoAffectedRowsOnClone         ID = "FIELD_002"
	FieldInvalidCharactersInKey        ID = "FIELD_003"
	FieldNotDefinedInSchema            ID = "FIELD_004"
	FieldTypeMismatch                  ID = "FIELD_005"
	FieldSameKeyForMultipleFields      ID = "FIELD_006"
	FieldRequiredMissing               ID = "FIELD_007"
	FieldInvalidType                   ID = "FIELD_008"
	FieldInvalidOwner                  ID = "FIELD_009"
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
)

const (
	SystemInternalError            ID = "SYS_001"
	SystemMisconfigured            ID = "SYS_002"
	SystemDependencyDown           ID = "SYS_003"
	SystemUnimplemented            ID = "SYS_004"
	SystemTransactionWithNoContext ID = "SYS_005"
	SystemErrorGeneratingUUID      ID = "SYS_006"
)
