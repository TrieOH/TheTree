package apierr

import (
	"github.com/MintzyG/fail"
)

var (
	RequestMissingQueryParamValue = fail.ID(0, "REQ", 0, false, "REQuestMissingQueryParamValue")
	RequestMissingQueryParam      = fail.ID(0, "REQ", 1, false, "REQuestMissingQueryParam")
	// FIXME create tests for empty cookies
	RequestEmptyCookie       = fail.ID(0, "REQ", 2, false, "REQuestEmptyCookie")
	RequestUnknownQueryParam = fail.ID(0, "REQ", 3, false, "REQuestUnknownQueryParam")
	RequestValidationError   = fail.ID(0, "REQ", 4, false, "REQuestValidationError")

	RequestMissingSchemaCustomFields = fail.ID(0, "REQ", 0, true, "REQuestMissingSchemaCustomFields")
	RequestInvalidJSONFormat         = fail.ID(0, "REQ", 1, true, "REQuestInvalidJSONFormat")
	RequestNotApplicationJSON        = fail.ID(0, "REQ", 3, true, "REQuestNotApplicationJSON")

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

	// Schema Version
	SchemaVersionDraftAlreadyExists  = fail.ID(1, "SCHEMAVERSION", 0, true, "SCHEMAVERSIONDraftAlreadyExists")
	SchemaVersionPublishWithNoFields = fail.ID(1, "SCHEMAVERSION", 1, true, "SCHEMAVERSIONPublishWithNoFields")

	SchemaVersionMismatch = fail.ID(1, "SCHEMAVERSION", 2, true, "SCHEMAVERSIONMismatch")
)

const (
	SchemaVersionDraftDoesntExist         ID = "SCM_VER_004"
	SchemaVersionTryingToPublishPublished ID = "SCM_VER_005"
	SchemaVersionTryingToPublishArchived  ID = "SCM_VER_006"
	// SchemaVersionMismatch                   ID = "SCM_VER_007"
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
	DBNotFound             ID = "DB_000"
	DBUniqueViolation      ID = "DB_001"
	DBForeignKeyViolation  ID = "DB_002"
	DBNotNullViolation     ID = "DB_003"
	DBValueTooLong         ID = "DB_004"
	DBSerializationFailure ID = "DB_005"
)

const PlaceholderID ID = "PL_000"

const (
	RoleNameTaken                   ID = "ROLE_003"
	RoleAlreadyGranted              ID = "ROLE_004"
	ScopeDuplicateNameAndExternalID ID = "SCP_002"
	ScopeInvalid                    ID = "SCP_003"
	PermissionAlreadyGranted        ID = "PERM_004"
	DBCheckViolation                ID = "DB_009"
)

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

	PERMissionInsufficient = fail.ID(0, "PERM", 0, true, "PERMissionInsufficient")

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

	ErrRoleNotOwnedByPrincipal = fail.Form(ROLENotOwnedByPrincipal, "role not owned by principal", false, map[string]any{"code": 401}).
					AddLocalization("pt-BR", "principal não é dono do papel")

	ErrRoleNameAlreadyTaken = fail.Form(ROLENameAlreadyTaken, "role name already taken", false, map[string]any{"code": 400}).
				AddLocalization("pt-BR", "nome do papel já em uso")
	ErrRoleAlreadyGranted = fail.Form(ROLEAlreadyGranted, "%s role already granted to user", false, map[string]any{"code": 400}, "ROLE NOT SET").
				AddLocalization("pt-BR", "papel %s já atribuído ao usuário")

	ErrScopeDuplicateNameAndExternalID = fail.Form(SCOPEDuplicateNameAndExternalID, "scope with name and external id (%s, %s) already exists", false, map[string]any{"code": 400}, "SCOPE NOT SET", "EXTERNAL_ID NOT SET").
						AddLocalization("pt-BR", "escopo com nome e id externo (%s, %s) já existe")

	ScopeInvalidShapeErrorMessage   = "invalid scope shape: a scope must be one of the following — (1) a global scope with type='global' and no project_id, name, or external_id; (2) a project root scope with type='project_root', a project_id, and no name or external_id; or (3) a project scope with type='project_scope', a project_id, and a name (external_id optional)"
	ScopeInvalidShapeErrorMessageBR = "Forma de escopo inválida: um escopo deve ser um dos seguintes — (1) um escopo global com type='global' e sem project_id, name ou external_id; (2) um escopo raiz de projeto com type='project_root', um project_id, e sem name ou external_id; ou (3) um escopo de projeto com type='project_scope', um project_id e um name (external_id opcional)."
	ErrScopeInvalidShape            = fail.Form(SCOPEInvalidShape, ScopeInvalidShapeErrorMessage, false, map[string]any{"code": 400}).
					AddLocalization("pt-BR", ScopeInvalidShapeErrorMessageBR)
	ErrScopeEmptyName = fail.Form(SCOPEEmptyName, "scope name cannot be empty", false, map[string]any{"code": 400}).
				AddLocalization("pt-BR", "nome do escopo não pode estar vazio")

	ErrPermissionLogicalConditionValidationError = fail.Form(PERMissionLogicalConditionValidationError, "%s: %s conditions cannot be empty", false, map[string]any{"code": 400}, "PATH NOT SET", "OPERATOR NOT SET").
							AddLocalization("pt-BR", "%s: condições %s não podem estar vazias")
	ErrPermissionConditionValidationError = fail.Form(PERMissionConditionValidationError, "error validating permission condition at: %s", false, map[string]any{"code": 400}, "PATH NOT SET").
						AddLocalization("pt-BR", "erro validando permissão da condição em: %s")
	ErrPermissionActionMismatch = fail.Form(PERMissionActionMismatch, "action mismatch: permission=%s, request=%s", false, map[string]any{"code": 400}, "EXPECTED NOT SET", "ACTUAL NOT SET").
					AddLocalization("pt-BR", "incompatibilidade de ações: permissão=%s, requisição=%s")
	ErrPermissionObjectMismatch = fail.Form(PERMissionObjectMismatch, "object mismatch: permission=%s, request=%s", false, map[string]any{"code": 400}, "EXPECTED NOT SET", "ACTUAL NOT SET").
					AddLocalization("pt-BR", "incompatibilidade de objetos: permissão=%s, requisição=%s")
	ErrPermissionNotOwnedByPrincipal = fail.Form(PERMissionNotOwnedByPrincipal, "permission not owned by principal", false, map[string]any{"code": 401}).
						AddLocalization("pt-BR", "identidade não é a dona da permissão")
	ErrPermissionInvalidAction = fail.Form(PERMissionInvalidAction, "invalid permission action: (%s)", false, map[string]any{"code": 400}, "ACTION NOT SET").
					AddLocalization("pt-BR", "ação da permissão é invalida: (%s)")
	ErrPermissionInvalidObject = fail.Form(PERMissionInvalidObject, "invalid permission object: (%s)", false, map[string]any{"code": 400}, "OBJECT NOT SET").
					AddLocalization("pt-BR", "objeto da permissão é invalido: (%s)")

	ErrInsufficientPermission = fail.Form(PERMissionInsufficient, "insufficient permissions", false, map[string]any{"code": 403}).
					AddLocalization("pt-BR", "permissões insuficientes")
	ErrPermissionAlreadyGranted = fail.Form(PERMissionAlreadyGranted, "user already has this permission in the specified scope", false, map[string]any{"code": 400}).
					AddLocalization("pt-BR", "usuário já tem essa permissão no escopo específicado")
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

	FIELDSamePositionForMultipleFields = fail.ID(0, "FIELD", 0, false, "FIELDSamePositionForMultipleFields")
	FIELDSameKeyForMultipleFields      = fail.ID(0, "FIELD", 1, false, "FIELDSameKeyForMultipleFields")
	FIELDInvalidCharactersInKey        = fail.ID(0, "FIELD", 2, false, "FIELDInvalidCharactersInKey")
)
