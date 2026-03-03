package errx

import "github.com/MintzyG/fail/v3"

var (
	// ------ REQ ------
	ErrRequestMissingQueryParamValue = fail.Form(RequestMissingQueryParamValue, "missing query parameter value for: %s", false, map[string]any{"code": 400}, "UNDEFINED").
						AddLocalizations(map[string]string{
			"pt-BR": "faltando parâmetro de pesquisa para: %s",
		})
	ErrRequestMissingQueryParam = fail.Form(RequestMissingQueryParam, "missing query parameter: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "faltando parâmetro de pesquisa: %s",
		})
	ErrRequestEmptyCookie = fail.Form(RequestEmptyCookie, "empty %s cookie value", false, map[string]any{"code": 400}, "UNDEFINED").
				AddLocalizations(map[string]string{
			"pt-BR": "o valor do cookie %s está vazio",
		})
	ErrRequestUnknownQueryParam = fail.Form(RequestUnknownQueryParam, "unknown query parameter: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "parâmetro de consulta desconhecido: %s",
		})
	ErrRequestValidationError = fail.Form(RequestValidationError, "Validation failed", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "formato do JSON inválido.",
		})
	ErrRequestParseUUIDError = fail.Form(RequestParseUUIDError, "invalid uuid field: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "o campo UUID: %s está inválido",
		})
	ErrRequestParseNumberError = fail.Form(RequestParseNumberError, "error parsing number: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "erro ao converter o número: %s",
		})
	ErrRequestMissingParamError = fail.Form(RequestMissingParamError, "missing parameter: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "faltando o parâmetro: %s",
		})
	ErrRequestInvalidCustomFieldsJSON = fail.Form(RequestInvalidCustomFieldsJSON, "invalid custom fields JSON", false, map[string]any{"code": 400}, "UNDEFINED").
						AddLocalizations(map[string]string{
			"pt-BR": "o JSON de campos customizados está inválido",
		})

	ErrRequestMissingSchemaCustomFields = fail.Form(RequestMissingSchemaCustomFields, "schema custom fields are required on a schema register", false, map[string]any{"code": 400}, "UNDEFINED").
						AddLocalizations(map[string]string{
			"pt-BR": "os campos personalizados do schema são obrigatórios no registro do schema.",
		})
	ErrRequestInvalidJSONFormat = fail.Form(RequestInvalidJSONFormat, "Invalid JSON format", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "formato do JSON inválido.",
		})
	ErrRequestNotApplicationJSON = fail.Form(RequestNotApplicationJSON, "Content-Type must be application/json", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "o Content-Type deve ser application/json",
		})
	ErrRequestInvalidPassword = fail.Form(RequestInvalidPassword, "invalid password", true, map[string]any{"code": 400}).
					AddLocalization("pt-BR", "senha inválida")

	ErrAuthInvalidAccessCookie = fail.Form(AuthInvalidAccessCookie, "invalid access_token cookie", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "o cookie access_token está inválido",
		})
	ErrAuthMissingAccessCookie = fail.Form(AuthMissingAccessCookie, "missing access_token cookie", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "o cookie access_token está está faltando",
		})
	ErrAuthSubjectNotInContext = fail.Form(AuthSubjectNotInContext, "authentication subject not found in request context", false, map[string]any{"code": 401}).
					AddLocalization("pt-BR", "sujeito de autenticação não encontrado no contexto da requisição")
	ErrAuthInvalidSubject = fail.Form(AuthInvalidSubject, "invalid authentication subject: %s", false, map[string]any{"code": 401}, "UNSET").
				AddLocalization("pt-BR", "sujeito de autenticação inválido: %s")

	ErrAuthzInsufficientPermissions = fail.Form(AuthzInsufficientPermissions, "user has insufficient permissions", false, map[string]any{"code": 403})

	ErrTokenInvalidAccessClaims = fail.Form(TokenInvalidAccessClaims, "invalid %s token claims", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "token %s com claims inválidas",
		})
	ErrTokenMissingSubClaim = fail.Form(TokenMissingSubClaim, "missing access sub claim", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "está faltando as claims do access token",
		})
	ErrTokenSubMarshalFailed   = fail.Form(TokenSubMarshalFailed, "failed to marshal sub claim", false, map[string]any{"code": 400})
	ErrTokenSubUnmarshalFailed = fail.Form(TokenSubUnmarshallingFailed, "failed to unmarshal sub claim into struct", false, map[string]any{"code": 400})

	// ------ VAL ------
	ErrValidationUUIDWasNil = fail.Form(ValidationUUIDWasNil, "%s field is nil", false, map[string]any{"code": 404}).
				AddLocalizations(map[string]string{
			"pt-BR": "%s está nulo",
		})

	// ------ EVENT ------
	ErrEventSlugAlreadyInUser = fail.Form(EventSlugAlreadyInUse, "slug already in use", false, map[string]any{"code": 409}).
					AddLocalization("pt-BR", "slug já está em uso")
	ErrEventPublishNonDraft = fail.Form(EventPublishNonDraft, "can't publish event in non draft status: status(%s)", false, map[string]any{"code": 400}, "UNSET").
				AddLocalization("pt-BR", "não pode se publicar um evento que não seja draft: status(%s)")
	ErrEventCannotAddEditions = fail.Form(EventCannotAddEditions, "cannot add editions to a non is_series event", false, map[string]any{"code": 400})

	ErrEditionInvalidID        = fail.Form(EditionInvalidID, "event_id is %s is invalid", false, map[string]any{"code": 400}, "UNSET")
	ErrEditionValidationFailed = fail.Form(EditionValidationFailed, "edition validation error", false, map[string]any{"code": 400})

	ErrTicketValidationFailed = fail.Form(TicketValidationFailed, "ticket validation error", false, map[string]any{"code": 400})

	// ------ SQL ------
	ErrSQLNotFound = fail.Form(SQLResourceNotFound, "%s not found", false, map[string]any{"code": 404}, "FORGOT TO SET RESOURCE ON ErrSQLNotFound").
			AddLocalization("pt-BR", "%s não foi encontrado")
	ErrInternalDBError = fail.Form(SQLInternalDBError, "internal DB error", true, map[string]any{"code": 500}).
				AddLocalization("pt-BR", "erro interno no banco de dados")
	ErrForeignKeyViolation = fail.Form(SQLForeignKeyViolation, "foreign key violation", false, map[string]any{"code": 409}).
				AddLocalization("pt-BR", "violação de chave estrangeira")
	ErrSerializationFailure = fail.Form(SQLSerializationFailure, "transaction conflict, retry", false, map[string]any{"code": 500}).
				AddLocalization("pt-BR", "conflito na transação, tente novamente")
	ErrNotNULLViolation = fail.Form(SQLNotNULLViolation, "missing required field", false, map[string]any{"code": 400}).
				AddLocalization("pt-BR", "está faltando um campo obrigatório")
	ErrValueTooLong = fail.Form(SQLValueTooLong, "value too long", false, map[string]any{"code": 400}).
			AddLocalization("pt-BR", "valor muito longo")
	ErrDBConnectionError = fail.Form(SQLDBConnectionError, "database connection error", true, map[string]any{"code": 500}).
				AddLocalization("pt-BR", "erro ao conectar no database")
	ErrSQLUnknownError = fail.Form(SQLDatabaseUnknownError, "SQL unknown error", true, map[string]any{"code": 500}).
				AddLocalization("pt-BR", "erro desconhecido no SQL")
	ErrSQLUnmatchedUniqueViolation = fail.Form(SQLUnmatchedUniqueViolation, "resource already exists", false, map[string]any{"code": 400}).
					AddLocalization("pt-BR", "recurso já existe")
	ErrSQLUnmatchedCheckViolation = fail.Form(SQLUnmatchedCheckViolation, "invalid value, violates a database constraint", false, map[string]any{"code": 400}).
					AddLocalization("pt-BR", "valor inválido, viola uma restrição da database")

	// ------ SYS ------
	ErrSysDependencyDown = fail.Form(SYSDependencyDown, "system dependency down: %s", true, map[string]any{"code": 500}, "UNNAMED DEPENDENCY").
				AddLocalization("pt-BR", "dependência do sistema está offline: %s")
	ErrServiceUnavailable = fail.Form(SYSServiceUnavailable, "%s is unavailable", true, map[string]any{"code": 500}, "UNNAMED SERVICE").
				AddLocalization("pt-BR", "%s está indisponível")
	ErrUUIDV7GenerationFailed = fail.Form(SYSUUIDV7GenerationError, "error generating UUID V7 at %s", true, map[string]any{"code": 500}, "UNSET LOCATION").
					AddLocalization("pt-BR", "erro gerando UUIDV7 em %s")

	ErrSystemFunctionalityNotImplemented = fail.Form(SYSFunctionalityNotImplemented, "this system functionality is not yet implemented", true, map[string]any{"code": 500}).
						AddLocalization("pt-BR", "essa funcionalidade do sistema ainda não foi implementada")
	ErrSystemTransactionNilContext = fail.Form(SYSTransactionNilContext, "cannot start transactions with a nil context", true, map[string]any{"code": 500}).
					AddLocalization("pt-BR", "não é possível começar uma transação com contexto nulo")

	// ------ DB ------
	ErrTransactionPanicked = fail.Form(DBTransactionPanicked, "transaction panicked", true, map[string]any{"code": 500}).
				AddLocalization("pt-BR", "transação resultou em pânico")
	ErrBeginTransactionFailed = fail.Form(DBBeginTransactionFailed, "could not begin transaction", true, map[string]any{"code": 500}).
					AddLocalization("pt-BR", "não foi possível começar a transação")
	ErrTransactionCommitFailed = fail.Form(DBTransactionCommitFailed, "transaction commit failed", true, map[string]any{"code": 500}).
					AddLocalization("pt-BR", "transação falhou no commit")

	ErrNestedTransactionNotAllowed = fail.Form(DBNestedTransactionNotAllowed, "nested transactions are not allowed", true, map[string]any{"code": 500}).
					AddLocalization("pt-BR", "transações aninhadas não são permitidas")
)
