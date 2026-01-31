package apierr

import "github.com/MintzyG/fail"

var (
	ErrSQLNotFound                 = fail.Form(SQLNotFound, "%s not found", false, map[string]any{"code": 404}, "FORGOT TO SET RESOURCE ON ErrSQLNotFound")
	ErrInternalDBError             = fail.Form(SQLInternalDBError, "internal DB error", true, map[string]any{"code": 500})
	ErrForeignKeyViolation         = fail.Form(SQLForeignKeyViolation, "foreign key violation", false, map[string]any{"code": 409})
	ErrSerializationFailure        = fail.Form(SQLSerializationFailure, "transaction conflict, retry", false, map[string]any{"code": 500})
	ErrNotNULLViolation            = fail.Form(SQLNotNULLViolation, "missing required field", false, map[string]any{"code": 400})
	ErrValueTooLong                = fail.Form(SQLValueTooLong, "value too long", false, map[string]any{"code": 400})
	ErrDBConnectionError           = fail.Form(SQLDBConnectionError, "database connection error", true, map[string]any{"code": 500})
	ErrSQLUnknownError             = fail.Form(SQLUnknownError, "SQL unknown error", true, map[string]any{"code": 500})
	ErrSQLUnmatchedUniqueViolation = fail.Form(SQLUnmatchedUniqueViolation, "resource already exists", false, map[string]any{"code": 400})
	ErrSQLUnmatchedCheckViolation  = fail.Form(SQLUnmatchedCheckViolation, "invalid value, violates a database constraint", false, map[string]any{"code": 400})

	ErrSchemaVersionDraftAlreadyExists = fail.Form(SCHEMAVersionDraftAlreadyExists, "a draft schema version already exists", false, map[string]any{"code": 400})

	ErrSamePositionForMultipleFields = fail.Form(FIELDSamePositionForMultipleFields, "two fields can't occupy the same position", false, map[string]any{"code": 400})
	ErrSameKeyForMultipleFields      = fail.Form(FIELDSameKeyForMultipleFields, "two fields can't have the same key", false, map[string]any{"code": 400})
	ErrInvalidCharacterInFieldKey    = fail.Form(FIELDInvalidCharactersInKey, "field key must start with a lowercase letter and contain only lowercase letters, numbers, or underscores", false, map[string]any{"code": 400})

	ErrScopeDuplicateNameAndExternalID = fail.Form(SCOPEDuplicateNameAndExternalID, "two scopes can't have the same name and external_id", false, map[string]any{"code": 400})
	ErrScopeInvalid                    = fail.Form(SCOPEInvalid, "invalid scope shape: a scope must be one of the following — (1) a global scope with type='global' and no project_id, name, or external_id; (2) a project root scope with type='project_root', a project_id, and no name or external_id; or (3) a project scope with type='project_scope', a project_id, and a name (external_id optional)", false, map[string]any{"code": 400})

	ErrRoleNameAlreadyTaken = fail.Form(ROLENameAlreadyTaken, "role name already taken", false, map[string]any{"code": 400})
	ErrRoleAlreadyGranted   = fail.Form(ROLEAlreadyGranted, "user already has this role in the specified scope", false, map[string]any{"code": 400})

	ErrPermissionAlreadyGranted = fail.Form(PERMISSIONAlreadyGranted, "user already has this permission in the specified scope", false, map[string]any{"code": 400})

	ErrRequestMissingQueryParamValue = fail.Form(RequestMissingQueryParamValue, "missing query parameter value for: %s", false, map[string]any{"code": 400}, "UNDEFINED").
						AddLocalizations(map[string]string{
			"pt-BR": "faltando parâmetro de pesquisa para: %s",
		})
	ErrRequestMissingQueryParam = fail.Form(RequestMissingQueryParam, "missing query parameter: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "faltando parâmetro de pesquisa: %s",
		})
	ErrRequestMissingSchemaCustomFields = fail.Form(RequestMissingSchemaCustomFields, "schema custom fields are required on a schema register", false, map[string]any{"code": 401}, "UNDEFINED").
						AddLocalizations(map[string]string{
			"pt-BR": "Os campos personalizados do schema são obrigatórios no registro do schema.",
		})
	ErrRequestInvalidJSONFormat = fail.Form(RequestInvalidJSONFormat, "Invalid JSON format", false, map[string]any{"code": 401}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "Formato do JSON inválido.",
		})
	ErrRequestValidationError = fail.Form(RequestValidationError, "Validation failed", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "Formato do JSON inválido.",
		})
	ErrRequestNotApplicationJSON = fail.Form(RequestNotApplicationJSON, "Content-Type must be application/json", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "O Content-Type deve ser application/json",
		})
	ErrRequestEmptyCookie = fail.Form(RequestEmptyCookie, "", false, map[string]any{"code": 400}, "UNDEFINED").
				AddLocalizations(map[string]string{
			"pt-BR": "",
		})
	ErrRequestUnknownQueryParam = fail.Form(RequestUnknownQueryParam, "unknown query parameter: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "parâmetro de consulta desconhecido: %s",
		})

	ErrAuthInvalidPrincipal = fail.Form(AuthInvalidPrincipal, "invalid principal", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "principal faltando",
		})
)
