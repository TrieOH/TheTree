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
	RequestNotApplicationJSON        = fail.ID(0, "REQ", 3, true, "REQuestNotApplicationJSON")
	RequestUnknownQueryParam         = fail.ID(0, "REQ", 5, true, "REQuestUnknownQueryParam")

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
	SessionSelfRevokeForbidden = fail.ID(1, "SESSION", 2, true, "SESSIONSelfRevokeForbi5dden")
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

	SCHEMANoPublishedVersion = fail.ID(0, "SCHEMA", 0, true, "SCHEMANoPublishedVersion")
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
	DBCheckViolation       ID = "DB_009"
)

const PlaceholderID ID = "PL_000"

var (
	SYSDependencyDown        = fail.ID(9, "SYS", 0, false, "SYStemDependencyDown")
	SYSServiceUnavailable    = fail.ID(9, "SYS", 1, false, "SYSServiceUnavailable")
	SYSJWKSRetrievalFailed   = fail.ID(9, "SYS", 2, false, "SYSJWKSRetrievalFailed")
	SYSRenderingEmailFailed  = fail.ID(9, "SYS", 3, false, "SYSRenderingEmailFailed")
	SYSUUIDV7GenerationError = fail.ID(9, "SYS", 4, false, "SYSUUIDV7GenerationError")

	SYSFunctionalityNotImplemented = fail.ID(9, "SYS", 0, true, "SYSFunctionalityNotImplemented")
	SYSTransactionNilContext       = fail.ID(9, "SYS", 1, true, "SYSTransactionNilContext")

	RequestInvalidPassword = fail.ID(0, "REQ", 6, true, "REQuestInvalidPassword") // FIXME there was a gap from 3 NIL 5, test extensively later

	DBTransactionPanicked     = fail.ID(9, "DB", 0, false, "DBTransactionPanicked")
	DBBeginTransactionFailed  = fail.ID(9, "DB", 1, false, "DBBeginTransactionFailed")
	DBTransactionCommitFailed = fail.ID(9, "DB", 2, false, "DBTransactionCommitFailed")

	DBNestedTransactionNotAllowed = fail.ID(9, "DB", 0, true, "DBNestedTransactionNotAllowed")

	ErrSysDependencyDown = fail.Form(SYSDependencyDown, "system dependency down: %s", true, map[string]any{"code": 500}, "UNNAMED DEPENDENCY").
				AddLocalization("pt-BR", "dependência do sistema está offline: %s")
	ErrServiceUnavailable = fail.Form(SYSServiceUnavailable, "%s is unavailable", true, map[string]any{"code": 500}, "UNNAMED SERVICE").
				AddLocalization("pt-BR", "%s está indisponível")
	ErrJWKSRetrievalFailed = fail.Form(SYSJWKSRetrievalFailed, "JWKS retrieval failed", true, map[string]any{"code": 500}).
				AddLocalization("pt-BR", "resgate do JWKS falhou")
	ErrRenderingEmailFailed = fail.Form(SYSRenderingEmailFailed, "%s email rendering failed", true, map[string]any{"code": 500}, "UNSET EMAIL TYPE").
				AddLocalization("pt-BR", "%s falhou na renderização")
	ErrUUIDV7GenerationFailed = fail.Form(SYSUUIDV7GenerationError, "error generating UUID V7 at %s", true, map[string]any{"code": 500}, "UNSET LOCATION").
					AddLocalization("pt-BR", "erro gerando UUIDV7 em %s")

	ErrSystemFunctionalityNotImplemented = fail.Form(SYSFunctionalityNotImplemented, "this system functionality is not yet implemented", true, map[string]any{"code": 500}).
						AddLocalization("pt-BR", "essa funcionalidade do sistema ainda não foi implementada")
	ErrSystemTransactionNilContext = fail.Form(SYSTransactionNilContext, "cannot start transactions with a nil context", true, map[string]any{"code": 500}).
					AddLocalization("pt-BR", "não é possível começar uma transação com contexto nulo")

	ErrRequestInvalidPassword = fail.Form(RequestInvalidPassword, "invalid password", true, map[string]any{"code": 400}).
					AddLocalization("pt-BR", "senha inválida")

	ErrTransactionPanicked = fail.Form(DBTransactionPanicked, "transaction panicked", true, map[string]any{"code": 500}).
				AddLocalization("pt-BR", "transação resultou em pânico")
	ErrBeginTransactionFailed = fail.Form(DBBeginTransactionFailed, "could not begin transaction", true, map[string]any{"code": 500}).
					AddLocalization("pt-BR", "não foi possível começar a transação")
	ErrTransactionCommitFailed = fail.Form(DBTransactionCommitFailed, "transaction commit failed", true, map[string]any{"code": 500}).
					AddLocalization("pt-BR", "transação falhou no commit")

	ErrNestedTransactionNotAllowed = fail.Form(DBNestedTransactionNotAllowed, "nested transactions are not allowed", true, map[string]any{"code": 500}).
					AddLocalization("pt-BR", "transações aninhadas não são permitidas")
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
