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
	ErrRequestInvalidSubContextJSON = fail.Form(RequestInvalidSubContextJSON, "invalid sub-context JSON", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "o JSON de sub-contexto está inválido",
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

	// ------ SESSION ------
	ErrSessionRevoked = fail.Form(SessionRevoked, "session not found or revoked", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "a sessão não foi encontrado ou revogada",
		})
	ErrSessionNotFound = fail.Form(SessionNotFound, "session not found or revoked", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "a sessão não foi encontrado ou revogada",
		})
	ErrSessionSelfRevokeForbidden = fail.Form(SessionSelfRevokeForbidden, "cannot revoke the currently active session", false, map[string]any{"code": 403}).
					AddLocalizations(map[string]string{
			"pt-BR": "não é possível revogar a sessão atual",
		})
	ErrSessionUnauthorized = fail.Form(SessionUnauthorized, "session not found or revoked", false, map[string]any{"code": 403}).
				AddLocalizations(map[string]string{
			"pt-BR": "a sessão não foi encontrado ou revogada",
		})

	// ------ TOKEN ------
	ErrTokenInvalid = fail.Form(TokenInvalid, "invalid %s token", false, map[string]any{"code": 401}).
			AddLocalizations(map[string]string{
			"pt-BR": "token %s inválido",
		})
	ErrTokenExpired = fail.Form(TokenInvalid, "%s token expired", false, map[string]any{"code": 401}).
			AddLocalizations(map[string]string{
			"pt-BR": "o Token %s está expirado",
		})
	ErrTokenMalformed = fail.Form(TokenMalformed, "malformed %s token", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o Token %s está inválido",
		})
	ErrTokenSignatureInvalid = fail.Form(TokenSignatureInvalid, "invalid %s token signature", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "a assinatura do Token %s está inválida",
		})
	ErrTokenInvalidAlg = fail.Form(TokenInvalidAlg, "invalid %s token alg, expected \"%s\" but got \"%s\"", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "token %s inválido, esperava \"%s\" mas obteve \"%s\"",
		})
	ErrTokenCouldNotSign = fail.Form(TokenCouldNotSign, "error signing %s token", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "erro ao assinar o token $s",
		})
	ErrTokenInvalidAccessClaims = fail.Form(TokenInvalidAccessClaims, "invalid %s token claims", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "token %s com claims inválidas",
		})
	ErrTokenNotYetValid = fail.Form(TokenNotYetValid, "%s token not yet valid", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o Token %s ainda não está válido",
		})
	ErrTokenUsedBeforeIssued = fail.Form(TokenUsedBeforeIssued, "%s token used before issued", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "o Token %s foi usado antes de ser emitido",
		})
	ErrTokenInvalidIssuer = fail.Form(TokenInvalidIssuer, "%s token has invalid issuer", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o Token %s possui um emissor inválido",
		})
	ErrTokenInvalidSubject = fail.Form(TokenInvalidSubject, "%s token has invalid subject", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o Token %s possui um asssunto inválido",
		})
	ErrTokenInvalidAudience = fail.Form(TokenInvalidAudience, "%s token has invalid audience", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o token %s tem um público inválido",
		})
	ErrTokenRefreshInvalidID = fail.Form(TokenRefreshInvalidID, "%s token has invalid id", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "o token %s tem um id inválido",
		})
	ErrTokenAccessInvalidID = fail.Form(TokenAccessInvalidID, "%s token has invalid id", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o token %s tem um id inválido",
		})
	ErrTokenInvalidKid = fail.Form(TokenInvalidKid, "invalid %s token kid", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o token %s kid é inválido",
		})
	ErrTokenUnknownKid = fail.Form(TokenUnknownKid, "unknown %s token kid", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o token %s kid é desconhecido",
		})
	ErrTokenMissingKid = fail.Form(TokenMissingKid, "%s token missing kid", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "está faltando o kid no token %s",
		})
	ErrTokenUnverifiable = fail.Form(TokenUnverifiable, "unverifiable %s token", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o token %s não é verificável",
		})
	ErrTokenReuseIdentified = fail.Form(TokenReuseIdentified, "%s token reuse not allowed", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "não é permitido a reutilização de token %s",
		})
	ErrTokenUserMismatch = fail.Form(TokenUserMismatch, "%s token user mismatch", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o usuário não coincide com o usuário do token %s",
		})
	ErrTokenInvalidFormat = fail.Form(TokenInvalidFormat, "invalid %s token format", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o formato do token %s está inválido",
		})
	ErrTokenUntrusted = fail.Form(TokenUntrusted, "untrusted %s token", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "token de %s não confiável",
		})

	ErrTokenSessionMismatch = fail.Form(TokenSessionMismatch, "token/session mismatch", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "há uma incompatibilidade entre token e sessão",
		})
	ErrTokenMismatchDuringAuth = fail.Form(TokenMismatchDuringAuth, "access token does not belong to this refresh token", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "o token de acesso não pertence a este token de atualização.",
		})
	ErrTokenMissingAccessClaims = fail.Form(TokenMissingAccessClaims, "missing access claims", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "está faltando as claims do access token",
		})
	ErrTokenMissingRefreshClaims = fail.Form(TokenMissingRefreshClaims, "missing refresh claims", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "está faltando as claims do refresh token",
		})

	ErrProjectErrorGeneratingKeys = fail.Form(ProjectErrorGeneratingKeys, "error generating project keys", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "erro ao gerar as chaves do projeto",
		})
	ErrProjectNotOwnedByPrincipal = fail.Form(ProjectNotOwnedByPrincipal, "%s", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "%s",
		})

	// ------ PROJECT ------
	ErrProjectNotFound = fail.Form(ProjectNotFound, "project not found", false, map[string]any{"code": 404}).
				AddLocalizations(map[string]string{
			"pt-BR": "projeto não encontrado",
		})

	// ------ PROJECTUSER ------
	ErrProjectUserErrorEncodingMetadata = fail.Form(ProjectUserErrorEncodingMetadata, "error encoding project user metadata", false, map[string]any{"code": 500}).
						AddLocalizations(map[string]string{
			"pt-BR": "erro ao codificar os metadados do projeto do usuário",
		})
	ErrProjectUserMetadataDecodeFailed = fail.Form(ProjectUserMetadataDecodeFailed, "error decoding project user metadata", false, map[string]any{"code": 500}).
						AddLocalizations(map[string]string{
			"pt-BR": "erro ao decodificar os metadados do projeto do usuário",
		})

	ErrProjectUserRegisterOnSchemaVersionDraft = fail.Form(ProjectUserRegisterOnSchemaVersionDraft, "can't register to a draft schema version", false, map[string]any{"code": 400}).
							AddLocalizations(map[string]string{
			"pt-BR": "você não pode se registrar a um rascunho da versão do schema",
		})
	ErrProjectUserRegisterOnSchemaDraft = fail.Form(ProjectUserRegisterOnSchemaDraft, "can't register to a draft schema", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "você não pode se registrar a um rascunho do schema",
		})
	ErrProjectUserRegisterOnSchemaArchived = fail.Form(ProjectUserRegisterOnSchemaArchived, "can't register to an archived schema", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "você não pode se registrar em um schema arquivado",
		})
	ErrProjectUserRegisterOnSchemaVersionArchived = fail.Form(ProjectUserRegisterOnSchemaVersionArchived, "can't register to an archived schema version", false, map[string]any{"code": 400}).
							AddLocalizations(map[string]string{
			"pt-BR": "você não pode se registrar em uma versão de schema arquivado",
		})
	ErrProjectUserNotFromProject = fail.Form(ProjectUserNotFromProject, "project user not from project", false, map[string]any{"code": 500}).
					AddLocalizations(map[string]string{
			"pt-BR": "o usuário não pertence a esse projeto",
		})
	ErrProjectUserRegisterOnSchemaNoVersion = fail.Form(ProjectUserRegisterOnSchemaNoVersion, "can't register on a schema that has no published version", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "você não pode se registrar em um schema que não possui versão publicada",
		})
	ErrProjectUserRegisterOnNoneProject = fail.Form(ProjectUserRegisterOnNoneProject, "can't register on a non existant project", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "não é possível se registrar em um projeto inexistente",
		})

	// ------ SCHEMA ------
	ErrSchemaNotOwnedByPrincipal = fail.Form(SchemaNotOwnedByPrincipal, "%s", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "%s",
		})
	ErrSchemaNoValidStatus = fail.Form(SchemaNoValidStatus, "CATASTROPHIC: schema found with no valid status: %s", false, map[string]any{"code": 500}).
				AddLocalizations(map[string]string{
			"pt-BR": "CATÁSTROFE: o schema encontrado não possui um status válido: %s",
		})
	ErrSchemaInvalidFlowID = fail.Form(SchemaInvalidFlowID, "invalid flow ID: %s", false, map[string]any{"code": 400}).
				AddLocalizations(map[string]string{
			"pt-BR": "o flow ID é inválido: %s",
		})
	ErrSchemaFlowIDIsReserved = fail.Form(SchemaFlowIDIsReserved, "flow id can't be the reserved keyword '%s'", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "o flow ID não podee ser essa palavra reservada '%s'",
		})

	ErrSchemaNoPublishedVersion = fail.Form(SCHEMANoPublishedVersion, "cannot publish a schema with no versions", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar um schema sem versões",
		})
	ErrSchemaFlowIDAlreadyExistsInType = fail.Form(SchemaFlowIDAlreadyExistsInType, "schema with this flow ID already exists in this type", false, map[string]any{"code": 409}).
						AddLocalizations(map[string]string{
			"pt-BR": "o schema com esse flow ID já existe nesse tipo",
		})
	ErrSchemaInvalidSchemaType = fail.Form(SchemaInvalidSchemaType, "invalid schema type", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "o tipo do schema é inválido",
		})
	ErrSchemaHasOnlyDraftVersion = fail.Form(SchemaHasOnlyDraftVersion, "cannot publish a schema with only draft versions", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar um schema com apenas versões de rascunhos",
		})
	ErrSchemaHasOnlyArchivedVersion = fail.Form(SchemaHasOnlyArchivedVersion, "cannot publish a schema with only archived versions", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar um schema com apenas versões arquivadas",
		})
	ErrSchemaTryingToPublishPublished = fail.Form(SchemaTryingToPublishPublished, "cannot publish a schema that is already published", false, map[string]any{"code": 401}).
						AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar um schema que já está publicado",
		})
	ErrSchemaTryingToPublishArchived = fail.Form(SchemaTryingToPublishArchived, "cannot publish a schema that is archived", false, map[string]any{"code": 401}).
						AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar um schema que está arquivado",
		})
	ErrSchemaMetadataNotAllowed = fail.Form(SchemaMetadataNotAllowed, "custom fields are not allowed for core schema", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "os campos personalizados não são permitidos no esquema principal",
		})
	ErrSchemaEmptySchemaType = fail.Form(SchemaEmptySchemaType, "schema type can't be empty", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "o tipo do schema não poder ser vazio",
		})
	ErrSchemaEmptyFlowID = fail.Form(SchemaEmptyFlowID, "flow id can't be empty", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "o flow ID não pode ser vazio",
		})

	// ------ SCHEMAVERSION ------
	ErrSchemaVersionNotDraft = fail.Form(SchemaVersionNotDraft, "cannot publish a schema version that isn't a draft", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar uma versão do schema que não seja um rascunho",
		})
	ErrSchemaVersionDraftAlreadyExists = fail.Form(SCHEMAVersionDraftAlreadyExists, "a draft schema version already exists", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "já existe um rascunho da versão desse schema",
		})

	ErrSchemaVersionPublishWithNoFields = fail.Form(SchemaVersionPublishWithNoFields, "cannot publish a schema version with no fields", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar uma versão do schema com nenhum campo",
		})
	ErrSchemaVersionDraftDoesntExist = fail.Form(SchemaVersionDraftDoesntExist, "cannot publish a schema with a version draft that doesn't exist", false, map[string]any{"code": 401}).
						AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar uma versão do schema de rascunho que não existe",
		})
	ErrSchemaVersionTryingToPublishPublished = fail.Form(SchemaVersionTryingToPublishPublished, "cannot publish a schema version that is already published", false, map[string]any{"code": 401}).
							AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar uma versão do schema que já está publicada",
		})
	ErrSchemaVersionTryingToPublishArchived = fail.Form(SchemaVersionTryingToPublishArchived, "cannot publish a schema version that is archived", false, map[string]any{"code": 401}).
						AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar uma versão do schema que está arquivada",
		})
	ErrSchemaVersionMismatch = fail.Form(SchemaVersionMismatch, "schema version and supplied version mismatch", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "a versão do schema e a versão fornecida não correspondem",
		})
	ErrSchemaVersionNonDraftAddFieldsNotAllowed = fail.Form(SchemaVersionNonDraftAddFieldsNotAllowed, "cannot add fields to a non-draft version", false, map[string]any{"code": 400}).
							AddLocalizations(map[string]string{
			"pt-BR": "não é possível adicionar campos em uma versão que não seja rascunho",
		})
	ErrSchemaVersionNoValidStatus = fail.Form(SchemaVersionNoValidStatus, "CATASTROPHIC: schema version found with no valid status", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "CATÁSTROFE: a versão do schema encontrada sem status válido",
		})
	ErrSchemaVersionDraftOnNonPublished = fail.Form(SchemaVersionDraftOnNonPublished, "new versions can only be drafted from published versions", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "novas versões só podem virar rascunhos a partir de versões publicadas",
		})
	ErrSchemaVersionNoChanges = fail.Form(SchemaVersionNoChanges, "cannot publish a version with no changes", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar uma versão sem mudanças",
		})
	ErrSchemaVersionTryingToPublishNonExistant = fail.Form(SchemaVersionTryingToPublishNonExistant, "cannot publish a non-existent schema version", false, map[string]any{"code": 400}).
							AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar uma versão inexistente",
		})
	ErrSchemaVersionNotPublished = fail.Form(SchemaVersionNotPublished, "version is not published", false, map[string]any{"code": 400}).
					AddLocalization("pt-BR", "versão não foi publicada")

	// ------ FIELD ------
	ErrFieldValidationErrorOnSchemaRegister = fail.Form(FIELDValidationErrorOnSchemaRegister, "error validating form for schema register", false, map[string]any{"code": 400}).
						AddLocalization("pt-BR", "erro validando formulário para o registro em schema")
	ErrFieldNotFound = fail.Form(FIELDNotFound, "field not found: %s", false, map[string]any{"code": 400}, "FIELD NOT SET").
				AddLocalization("pt-BR", "campo nao encontrado: %s")
	ErrFieldInvalidOwner = fail.Form(FIELDInvalidOwner, "invalid owner type (%s) for field: %s", false, map[string]any{"code": 400}, "UNSET OWNER", "UNSET FIELD KEY").
				AddLocalization("pt-BR", "tipo de dono inválido (%s) para o campo: %s")
	ErrFieldNoAffectedRowsOnClone = fail.Form(FieldNoAffectedRowsOnClone, "no affected rows", false, map[string]any{"code": 404}).
					AddLocalizations(map[string]string{
			"pt-BR": "nenhuma linha afetada",
		})
	ErrFieldInvalidType = fail.Form(FIELDInvalidType, "invalid field type (%s) for field: %s", false, map[string]any{"code": 400}, "UNSET TYPE", "UNSET FIELD KEY").
				AddLocalization("pt-BR", "tipo do campo inválido (%s) para o campo: %s")
	ErrSameKeyForMultipleFields = fail.Form(FIELDSameKeyForMultipleFields, "two fields can't have the same key", false, map[string]any{"code": 409}).
					AddLocalization("pt-BR", "dois campos não podem possuir a mesma chave identificadora")
	ErrFieldSamePositionForMultipleFields = fail.Form(FIELDSamePositionForMultipleFields, "two fields can't occupy the same position", false, map[string]any{"code": 409}).
						AddLocalizations(map[string]string{
			"pt-BR": "dois campos não podem ocupar a mesma posição",
		})
	ErrInvalidCharacterInFieldKey = fail.Form(FIELDInvalidCharactersInKey, "field key must start with a lowercase letter and contain only lowercase letters, numbers, or underscores", false, map[string]any{"code": 400}).
					AddLocalization("pt-BR", "a chave do campo deve começar com uma letra minúscula e conter apenas letras minúsculas, números ou sublinhados.")
	ErrFieldKeyAlreadyExists = fail.Form(FIELDKeyAlreadyExists, "field key '%s' already exists in this version", false, map[string]any{"code": 409}, "KEY NOT SET").
					AddLocalization("pt-BR", "chave do campo '%s' já existe nesta versão")
	ErrFieldHasDependentRules = fail.Form(FIELDHasDependentRules, "cannot delete field: referenced by visibility/required rules of other fields: %v", false, map[string]any{"code": 409}, "DEPENDENT FIELDS NOT SET").
					AddLocalization("pt-BR", "não é possível excluir o campo: referenciado por regras de visibilidade/obrigatoriedade de outros campos: %v")

	// ------ VAL ------
	ErrValidationUUIDWasNil = fail.Form(ValidationUUIDWasNil, "%s field is nil", false, map[string]any{"code": 404}).
				AddLocalizations(map[string]string{
			"pt-BR": "%s está nulo",
		})

	// ------ FORM ------
	ErrFormMissingRequiredField = fail.Form(FORMMissingRequiredField, "form missing required field: %s", false, map[string]any{"code": 400}, "UNSET FIELD KEY").
					AddLocalization("pt-BR", "formulário com campos obrigatorio ausente: %s")
	ErrFormInvalidFieldValue = fail.Form(FORMInvalidFieldValue, "invalid form value for %s: type(%v) value(%v)", false, map[string]any{"code": 400}, "UNSET FIELD KEY", "UNSET FIELD TYPE", "UNSET FIELD VALUE").
					AddLocalization("pt-BR", "campo %s do formulário com valor inválido: tipo(%v) valor(%v)")

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
	ErrRenderingEmailFailed = fail.Form(SYSRenderingEmailFailed, "%s email rendering failed", true, map[string]any{"code": 500}, "UNSET EMAIL TYPE").
				AddLocalization("pt-BR", "%s falhou na renderização")
	ErrUUIDV7GenerationFailed = fail.Form(SYSUUIDV7GenerationError, "error generating UUID V7 at %s", true, map[string]any{"code": 500}, "UNSET LOCATION").
					AddLocalization("pt-BR", "erro gerando UUIDV7 em %s")
	ErrSYSJWKSEncodingFailed = fail.Form(SYSJWKSEncodingFailed, "JWKS enconding failed", true, map[string]any{"code": 500}, "UNSET LOCATION").
					AddLocalization("pt-BR", "falha ao codificar o JWKS")
	ErrSysCryptoError = fail.Form(SYSCryptoError, "cryptographic error: %s", true, map[string]any{"code": 500}, "UNDEFINED").
				AddLocalization("pt-BR", "erro criptográfico: %s")
	ErrSYSInvalidPublicKeyPEM = fail.Form(SYSInvalidPublicKeyPEM, "invalid public key PEM: %s", false, map[string]any{"code": 500}, "UNSET").
					AddLocalization("pt-BR", "PEM da chave pública inválido: %s")
	ErrSYSPublicKeyParseFailed = fail.Form(SYSPublicKeyParseFailed, "failed to parse public key: %s", false, map[string]any{"code": 500}, "UNSET").
					AddLocalization("pt-BR", "falha ao interpretar chave pública: %s")
	ErrSYSInvalidPublicKeyType = fail.Form(SYSInvalidPublicKeyType, "invalid public key type: %s", false, map[string]any{"code": 500}, "UNSET").
					AddLocalization("pt-BR", "tipo de chave pública inválido: %s")
	ErrSYSInvalidPublicKeySize = fail.Form(SYSInvalidPublicKeyByteSize, "invalid public key size: %s", false, map[string]any{"code": 500}, "UNSET").
					AddLocalization("pt-BR", "tamanho da chave pública inválido: %s")
	ErrSYSErrorCreatingKey = fail.Form(SYSErrorCreatingKey, "error generating signing key", true, map[string]any{"code": 500})

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

	// ------ ROLE ------
	ErrRoleNotOwnedByPrincipal = fail.Form(ROLENotOwnedByPrincipal, "role not owned by principal", false, map[string]any{"code": 401}).
					AddLocalization("pt-BR", "principal não é dono do papel")
	ErrRoleAlreadyGranted = fail.Form(ROLEAlreadyGranted, "%s role already granted to user", false, map[string]any{"code": 400}, "ROLE NOT SET").
				AddLocalization("pt-BR", "papel %s já atribuído ao usuário")

	ErrRoleNameAlreadyTaken = fail.Form(ROLENameAlreadyTaken, "role name already taken", false, map[string]any{"code": 409}).
				AddLocalization("pt-BR", "nome do papel já em uso")

	// ------ SCOPE ------
	ErrScopeDuplicateNameAndExternalID = fail.Form(SCOPEDuplicateNameAndExternalID, "scope with name and external id (%s, %s) already exists", false, map[string]any{"code": 409}, "SCOPE NOT SET", "EXTERNAL_ID NOT SET").
						AddLocalization("pt-BR", "escopo com nome e id externo (%s, %s) já existe")

	ScopeInvalidShapeErrorMessage   = "invalid scope shape: a scope must be one of the following — (1) a global scope with type='global' and no project_id, name, or external_id; (2) a project root scope with type='project_root', a project_id, and no name or external_id; or (3) a project scope with type='project_scope', a project_id, and a name (external_id optional)"
	ScopeInvalidShapeErrorMessageBR = "Forma de escopo inválida: um escopo deve ser um dos seguintes — (1) um escopo global com type='global' e sem project_id, name ou external_id; (2) um escopo raiz de projeto com type='project_root', um project_id, e sem name ou external_id; ou (3) um escopo de projeto com type='project_scope', um project_id e um name (external_id opcional)."
	ErrScopeInvalidShape            = fail.Form(SCOPEInvalidShape, ScopeInvalidShapeErrorMessage, false, map[string]any{"code": 400}).
					AddLocalization("pt-BR", ScopeInvalidShapeErrorMessageBR)
	ErrScopeEmptyName = fail.Form(SCOPEEmptyName, "scope name cannot be empty", false, map[string]any{"code": 400}).
				AddLocalization("pt-BR", "nome do escopo não pode estar vazio")
	ErrSCOPEOneGlobal = fail.Form(SCOPEOneGlobal, "only one global scope may exist", true, map[string]any{"code": 409}).
				AddLocalization("pt-BR", "apenas um escopo global deve existir")
	ErrSCOPEOneProjectRootPerProject = fail.Form(SCOPEOneProjectRootPerProject, "only one project_root scope may exist per project", true, map[string]any{"code": 409}).
						AddLocalization("pt-BR", "apenas um escopo raiz de projeto deve existir por projeto")
	ErrScopeParentNotFound = fail.Form(SCOPEParentNotFound, "parent scope not found", false, map[string]any{"code": 404}).
				AddLocalization("pt-BR", "escopo pai não encontrado")
	ErrScopeParentDifferentProject = fail.Form(SCOPEParentDifferentProject, "parent scope belongs to a different project", false, map[string]any{"code": 400}).
					AddLocalization("pt-BR", "escopo pai pertence a um projeto diferente")
	ErrScopeHierarchyError = fail.Form(SCOPEHierarchyError, "error checking scope hierarchy", false, map[string]any{"code": 500}).
				AddLocalization("pt-BR", "erro ao verificar hierarquia de escopo")
	ErrScopeCycleDetected = fail.Form(SCOPECycleDetected, "scope hierarchy cycle detected", false, map[string]any{"code": 400}).
				AddLocalization("pt-BR", "ciclo detectado na hierarquia de escopos")
	ErrScopeRootNotFound = fail.Form(SCOPERootNotFound, "project root scope not found", false, map[string]any{"code": 500}).
				AddLocalization("pt-BR", "escopo raiz do projeto não encontrado")
	ErrScopeDuplicateSibling = fail.Form(SCOPEDuplicateSibling, "scope with same name already exists under this parent", false, map[string]any{"code": 409}).
					AddLocalization("pt-BR", "escopo com o mesmo nome já existe sob este pai")

	// ------ PERM ------
	ErrPermissionLogicalConditionValidationError = fail.Form(PERMissionLogicalConditionValidationError, "%s: %s conditions cannot be empty", false, map[string]any{"code": 400}, "PATH NOT SET", "OPERATOR NOT SET").
							AddLocalization("pt-BR", "%s: condições %s não podem estar vazias")
	ErrPermissionConditionValidationError = fail.Form(PERMissionConditionValidationError, "error validating permission condition at: %s", false, map[string]any{"code": 400}, "PATH NOT SET").
						AddLocalization("pt-BR", "erro validando permissão da condição em: %s")
	ErrPermissionActionMismatch = fail.Form(PERMissionActionMismatch, "action mismatch: permission=%s, request=%s", false, map[string]any{"code": 400}, "EXPECTED NOT SET", "ACTUAL NOT SET").
					AddLocalization("pt-BR", "incompatibilidade de ações: permissão=%s, requisição=%s")
	ErrPermissionObjectMismatch = fail.Form(PERMissionObjectMismatch, "object mismatch: permission=%s, request=%s", false, map[string]any{"code": 400}, "EXPECTED NOT SET", "ACTUAL NOT SET").
					AddLocalization("pt-BR", "incompatibilidade de objetos: permissão=%s, requisição=%s")
	ErrPermissionAlreadyGranted = fail.Form(PERMissionAlreadyGranted, "user already has this permission in the specified scope", false, map[string]any{"code": 409}).
					AddLocalization("pt-BR", "usuário já tem essa permissão no escopo específicado")
	ErrPermissionNotOwnedByPrincipal = fail.Form(PERMissionNotOwnedByPrincipal, "permission not owned by principal", false, map[string]any{"code": 401}).
						AddLocalization("pt-BR", "identidade não é a dona da permissão")
	ErrPermissionInvalidAction = fail.Form(PERMissionInvalidAction, "invalid permission action: (%s)", false, map[string]any{"code": 400}, "ACTION NOT SET").
					AddLocalization("pt-BR", "ação da permissão é invalida: (%s)")
	ErrPermissionInvalidObject = fail.Form(PERMissionInvalidObject, "invalid permission object: (%s)", false, map[string]any{"code": 400}, "OBJECT NOT SET").
					AddLocalization("pt-BR", "objeto da permissão é invalido: (%s)")
	ErrPermissionAlreadyExists = fail.Form(PERMissionAlreadyExists, "permission with object(%s) and action(%s) already exists", false, map[string]any{"code": 409}, "OBJECT NOT SET", "ACTION NOT SET").
					AddLocalization("pt-BR", "permissão com objeto(%s) e ação(%s) já existe")

	ErrInsufficientPermission = fail.Form(PERMissionInsufficient, "insufficient permissions", false, map[string]any{"code": 403}).
					AddLocalization("pt-BR", "permissões insuficientes")
	ErrPermissionNoResource = fail.Form(PERMissionNoResource, "can't check permission(%s) conditions without resource", false, map[string]any{"code": 403}, "PERMISSION ID NOT SET").
				AddLocalization("pt-BR", "não é possível verificar as condições da permissão(%s) sem o recurso")

	// EMAIL
	ErrEMAILTemplateNotFound = fail.Form(EMAILTemplateNotFound, "%s %s template not found", false, map[string]any{"code": 500}, "KEY NOT SET", "TYPE NOT SET").
					AddLocalization("pt-BR", "O template %s %s não foi encontrado")
)
