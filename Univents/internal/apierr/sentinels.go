package apierr

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

	// ------ AUTH ------
	ErrAuthEmailAlreadyUsed = fail.Form(AuthEmailAlreadyUsed, "email already in use", false, map[string]any{"code": 409}).
				AddLocalizations(map[string]string{
			"pt-BR": "erro ao registrar o usuário",
		})
	ErrAuthInvalidCredentials = fail.Form(AuthInvalidCredentials, "invalid email or password", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "email ou senha inválidos",
		})
	ErrAuthInvalidRefreshCookie = fail.Form(AuthInvalidRefreshCookie, "invalid refresh_token cookie", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "o cookie refresh_token está inválido",
		})
	ErrAuthInvalidAccessCookie = fail.Form(AuthInvalidAccessCookie, "invalid access_token cookie", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "o cookie access_token está inválido",
		})
	ErrAuthMissingRefreshCookie = fail.Form(AuthMissingRefreshCookie, "missing refresh_token cookie", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "o cookie refresh_token está está faltando",
		})
	ErrAuthMissingAccessCookie = fail.Form(AuthMissingAccessCookie, "missing access_token cookie", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "o cookie access_token está está faltando",
		})

	ErrAuthInvalidPrincipal = fail.Form(AuthInvalidPrincipal, "invalid principal", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "principal faltando",
		})
	ErrAuthInvalidPassword = fail.Form(AuthInvalidPassword, "password length exceeds 72 bytes", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o comprimento da senha excede 72 bytes",
		})
	ErrAuthNotClient = fail.Form(AuthNotClient, "only clients can access this endpoint", false, map[string]any{"code": 403}).
				AddLocalizations(map[string]string{
			"pt-BR": "apenas clientes podem acessar esse endpoint",
		})
	ErrAuthNotProjectUser = fail.Form(AuthNotProjectUser, "only project users can access this endpoint", false, map[string]any{"code": 403}).
				AddLocalizations(map[string]string{
			"pt-BR": "apenas usuários do projeto podem acessar esse endpoint",
		})
	ErrAuthAlreadyVerified = fail.Form(AuthAlreadyVerified, "user already verified", false, map[string]any{"code": 403}).
				AddLocalizations(map[string]string{
			"pt-BR": "o usuário já foi verificado",
		})
	ErrAuthPrincipalNotInContext = fail.Form(AuthPrincipalNotInContext, "missing principal in context", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "está faltando o principal no contexto",
		})
	ErrAuthUserSchemaOutdated = fail.Form(AuthUserSchemaOutdated, "user schema is outdated, please upgrade your metadata", false, map[string]any{"code": 403}).
					AddLocalizations(map[string]string{
			"pt-BR": "o schema do usuário está desatualizado, por favor atualize seus metadados",
		})
	ErrAuthTokenAlreadyUsed = fail.Form(AuthTokenAlreadyUsed, "token already used", false, map[string]any{"code": 403}).
				AddLocalization("pt-BR", "este token já foi usado")
	ErrAuthInvalidApiKey = fail.Form(AuthInvalidApiKey, "invalid api key", false, map[string]any{"code": 401}).
				AddLocalization("pt-BR", "api key inválida")
	ErrAuthInvalidApiKeyShape = fail.Form(AuthInvalidApiKeyShape, "invalid api key shape", false, map[string]any{"code": 401}).
					AddLocalization("pt-BR", "formato de api key inválido")
	ErrAuthApiKeyNotAllowed = fail.Form(AuthApiKeyNotAllowed, "api keys are not allowed for this endpoint", false, map[string]any{"code": 403}).
				AddLocalization("pt-BR", "api keys não são permitidas para este endpoint")

	// ------ VAL ------
	ErrValidationUUIDWasNil = fail.Form(ValidationUUIDWasNil, "%s field is nil", false, map[string]any{"code": 404}).
				AddLocalizations(map[string]string{
			"pt-BR": "%s está nulo",
		})

	// ------ SQL ------
	ErrSQLNotFound = fail.Form(SQLNotFound, "%s not found", false, map[string]any{"code": 404}, "FORGOT TO SET RESOURCE ON ErrSQLNotFound").
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
	ErrSQLUnknownError = fail.Form(SQLUnknownError, "SQL unknown error", true, map[string]any{"code": 500}).
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
	ErrJWKSRetrievalFailed = fail.Form(SYSJWKSRetrievalFailed, "JWKS retrieval failed", true, map[string]any{"code": 500}).
				AddLocalization("pt-BR", "resgate do JWKS falhou")
	ErrUUIDV7GenerationFailed = fail.Form(SYSUUIDV7GenerationError, "error generating UUID V7 at %s", true, map[string]any{"code": 500}, "UNSET LOCATION").
					AddLocalization("pt-BR", "erro gerando UUIDV7 em %s")
	ErrSysCryptoError = fail.Form(SYSCryptoError, "cryptographic error: %s", true, map[string]any{"code": 500}, "UNDEFINED").
				AddLocalization("pt-BR", "erro criptográfico: %s")

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
