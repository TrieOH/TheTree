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

	ErrInvalidCharacterInFieldKey    = fail.Form(FIELDInvalidCharactersInKey, "field key must start with a lowercase letter and contain only lowercase letters, numbers, or underscores", false, map[string]any{"code": 400})
	ErrSamePositionForMultipleFields = fail.Form(FIELDSamePositionForMultipleFields, "two fields can't occupy the same position", false, map[string]any{"code": 400})

	ErrRequestMissingQueryParamValue = fail.Form(RequestMissingQueryParamValue, "missing query parameter value for: %s", false, map[string]any{"code": 400}, "UNDEFINED").
						AddLocalizations(map[string]string{
			"pt-BR": "faltando parâmetro de pesquisa para: %s",
		})
	ErrRequestMissingQueryParam = fail.Form(RequestMissingQueryParam, "missing query parameter: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "faltando parâmetro de pesquisa: %s",
		})
	ErrRequestMissingCustomFields = fail.Form(RequestMissingCustomFields, "schema custom fields are required on a schema register", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "Os campos customizados do schema são obrigatórios no registro de schema",
		})
	ErrRequestEmptyCookie = fail.Form(RequestEmptyCookie, "empty %s cookie value", false, map[string]any{"code": 400}, "UNDEFINED").
				AddLocalizations(map[string]string{
			"pt-BR": "O valor do cookie %s está vazio",
		})
	ErrRequestUnknownQueryParam = fail.Form(RequestUnknownQueryParam, "unknown query parameter: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "parâmetro de consulta desconhecido: %s",
		})
	ErrRequestValidationError = fail.Form(RequestValidationError, "Validation failed", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "Formato do JSON inválido.",
		})
	ErrRequestParseUUIDError = fail.Form(RequestParseUUIDError, "invalid uuid field: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "O campo UUID: %s está inválido",
		})
	ErrRequestParseNumberError = fail.Form(RequestParseNumberError, "error parsing number: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "Erro ao converter o número: %s",
		})
	ErrRequestMissingParamError = fail.Form(RequestMissingParamError, "missing parameter: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "Faltando o parâmetro: %s",
		})
	ErrRequestInvalidCustomFieldsJSON = fail.Form(RequestInvalidCustomFieldsJSON, "invalid custom fields JSON", false, map[string]any{"code": 400}, "UNDEFINED").
						AddLocalizations(map[string]string{
			"pt-BR": "O JSON de campos customizados está inválido",
		})
	ErrRequestMissingSchemaCustomFields = fail.Form(RequestMissingSchemaCustomFields, "schema custom fields are required on a schema register", false, map[string]any{"code": 401}, "UNDEFINED").
						AddLocalizations(map[string]string{
			"pt-BR": "Os campos personalizados do schema são obrigatórios no registro do schema.",
		})
	ErrRequestInvalidJSONFormat = fail.Form(RequestInvalidJSONFormat, "Invalid JSON format", false, map[string]any{"code": 401}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "Formato do JSON inválido.",
		})
	ErrRequestNotApplicationJSON = fail.Form(RequestNotApplicationJSON, "Content-Type must be application/json", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "O Content-Type deve ser application/json",
		})

	ErrAuthEmailAlreadyUsed = fail.Form(AuthEmailAlreadyUsed, "error registering user", false, map[string]any{"code": 409}).
				AddLocalizations(map[string]string{
			"pt-BR": "Erro ao registrar o usuário",
		})
	ErrAuthInvalidCredentials = fail.Form(AuthInvalidCredentials, "invalid email or password", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "Email ou Senha inválida",
		})
	ErrAuthInvalidRefreshCookie = fail.Form(AuthInvalidRefreshCookie, "invalid refresh_token cookie", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "O Cookie refresh_token está inválido",
		})
	ErrAuthInvalidAccessCookie = fail.Form(AuthInvalidAccessCookie, "invalid access_token cookie", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "O Cookie access_token está inválido",
		})
	ErrAuthMissingRefreshCookie = fail.Form(AuthMissingRefreshCookie, "missing refresh_token cookie", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "O Cookie refresh_token está está faltando",
		})
	ErrAuthMissingAccessCookie = fail.Form(AuthMissingAccessCookie, "missing access_token cookie", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "O Cookie access_token está está faltando",
		})
	ErrAuthInvalidPrincipal = fail.Form(AuthInvalidPrincipal, "invalid principal", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "principal faltando",
		})
	ErrAuthInvalidPassword = fail.Form(AuthInvalidPassword, "password length exceeds 72 bytes", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O comprimento da senha excede 72 bytes",
		})
	ErrAuthNotClient = fail.Form(AuthNotClient, "only clients can access this endpoint", false, map[string]any{"code": 403}).
				AddLocalizations(map[string]string{
			"pt-BR": "Apenas clientes podem acessar esse endpoint",
		})
	ErrAuthNotProjectUser = fail.Form(AuthNotProjectUser, "only project users can access this endpoint", false, map[string]any{"code": 403}).
				AddLocalizations(map[string]string{
			"pt-BR": "Apenas usuários do projeto podem acessar esse endpoint",
		})
	ErrAuthAlreadyVerified = fail.Form(AuthAlreadyVerified, "user already verified", false, map[string]any{"code": 403}).
				AddLocalizations(map[string]string{
			"pt-BR": "O usuário já foi verificado",
		})
	ErrAuthPrincipalNotInContext = fail.Form(AuthPrincipalNotInContext, "missing principal in context", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "Está faltando o principal no contexto",
		})

	ErrSessionRevoked = fail.Form(SessionRevoked, "session not found or revoked", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "A sessão não foi encontrado ou revogada",
		})
	ErrSessionNotFound = fail.Form(SessionNotFound, "session not found or revoked", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "A sessão não foi encontrado ou revogada",
		})
	ErrSessionSelfRevokeForbidden = fail.Form(SessionSelfRevokeForbidden, "cannot revoke the currently active session", false, map[string]any{"code": 403}).
					AddLocalizations(map[string]string{
			"pt-BR": "Não é possível revogar a sessão atual",
		})
	ErrSessionUnauthorized = fail.Form(SessionUnauthorized, "session not found or revoked", false, map[string]any{"code": 403}).
				AddLocalizations(map[string]string{
			"pt-BR": "A sessão não foi encontrado ou revogada",
		})

	ErrTokenInvalid = fail.Form(TokenInvalid, "invalid %s token", false, map[string]any{"code": 401}).
			AddLocalizations(map[string]string{
			"pt-BR": "Token %s inválido",
		})
	ErrTokenExpired = fail.Form(TokenInvalid, "%s token expired", false, map[string]any{"code": 401}).
			AddLocalizations(map[string]string{
			"pt-BR": "O Token %s está expirado",
		})
	ErrTokenMalformed = fail.Form(TokenMalformed, "malformed %s token", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O Token %s está inválido",
		})
	ErrTokenSignatureInvalid = fail.Form(TokenSignatureInvalid, "invalid %s token signature", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "A assinatura do Token %s está inválida",
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
			"pt-BR": "Token %s com claims inválidas",
		})
	ErrTokenNotYetValid = fail.Form(TokenNotYetValid, "%s token not yet valid", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O Token %s ainda não está válido",
		})
	ErrTokenUsedBeforeIssued = fail.Form(TokenUsedBeforeIssued, "%s token used before issued", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "O Token %s foi usado antes de ser emitido",
		})
	ErrTokenInvalidIssuer = fail.Form(TokenInvalidIssuer, "%s token has invalid issuer", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O Token %s possui um emissor inválido",
		})
	ErrTokenInvalidSubject = fail.Form(TokenInvalidSubject, "%s token has invalid subject", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O Token %s possui um asssunto inválido",
		})
	ErrTokenInvalidAudience = fail.Form(TokenInvalidAudience, "%s token has invalid audience", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O token %s tem um público inválido",
		})
	ErrTokenRefreshInvalidID = fail.Form(TokenRefreshInvalidID, "%s token has invalid id", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "O token %s tem um id inválido",
		})
	ErrTokenAccessInvalidID = fail.Form(TokenAccessInvalidID, "%s token has invalid id", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O token %s tem um id inválido",
		})
	ErrTokenInvalidKid = fail.Form(TokenInvalidKid, "invalid %s token kid", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O token %s kid é inválido",
		})
	ErrTokenUnknownKid = fail.Form(TokenUnknownKid, "unknown %s token kid", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O token %s kid é desconhecido",
		})
	ErrTokenMissingKid = fail.Form(TokenMissingKid, "%s token missing kid", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "Está faltando o kid no token %s",
		})
	ErrTokenUnverifiable = fail.Form(TokenUnverifiable, "unverifiable %s token", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O token %s não é verificável",
		})
	ErrTokenReuseIdentified = fail.Form(TokenReuseIdentified, "%s token reuse not allowed", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "Não é permitido a reutilização de token %s",
		})
	ErrTokenUserMismatch = fail.Form(TokenUserMismatch, "%s token user mismatch", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O usuário não coincide com o usuário do token %s",
		})
	ErrTokenInvalidFormat = fail.Form(TokenInvalidFormat, "invalid %s token format", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O formato do token %s está inválido",
		})
	ErrTokenUntrusted = fail.Form(TokenUntrusted, "untrusted %s token", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "Token de %s não confiável",
		})
	ErrTokenSessionMismatch = fail.Form(TokenSessionMismatch, "token/session mismatch", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "Há uma incompatibilidade entre token e sessão",
		})
	ErrTokenMismatchDuringAuth = fail.Form(TokenMismatchDuringAuth, "access token does not belong to this refresh token", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "O token de acesso não pertence a este token de atualização.",
		})
	ErrTokenMissingAccessClaims = fail.Form(TokenMissingAccessClaims, "missing access claims", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "Está faltando as claims do access token",
		})
	ErrTokenMissingRefreshClaims = fail.Form(TokenMissingRefreshClaims, "missing refresh claims", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "Está faltando as claims do refresh token",
		})

	ErrProjectErrorGeneratingKeys = fail.Form(ProjectErrorGeneratingKeys, "error generating project keys", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "Erro ao gerar as chaves do projeto",
		})
	ErrProjectNotOwnedByPrincipal = fail.Form(ProjectNotOwnedByPrincipal, "%s", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "%s",
		})
	ErrProjectNotFound = fail.Form(ProjectNotFound, "project not found", false, map[string]any{"code": 404}).
				AddLocalizations(map[string]string{
			"pt-BR": "Projeto não encontrado",
		})

	ErrProjectUserRegisterOnSchemaVersionDraft = fail.Form(ProjectUserRegisterOnSchemaVersionDraft, "can't register to a draft schema version", false, map[string]any{"code": 400}).
							AddLocalizations(map[string]string{
			"pt-BR": "Você não pode se registrar a um rascunho da versão do schema",
		})
	ErrProjectUserRegisterOnSchemaDraft = fail.Form(ProjectUserRegisterOnSchemaDraft, "can't register to a draft schema", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "Você não pode se registrar a um rascunho do schema",
		})
	ErrProjectUserRegisterOnSchemaArchived = fail.Form(ProjectUserRegisterOnSchemaArchived, "can't register to an archived schema", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "Você não pode se registrar em um schema arquivado",
		})
	ErrProjectUserRegisterOnSchemaVersionArchived = fail.Form(ProjectUserRegisterOnSchemaVersionArchived, "can't register to an archived schema version", false, map[string]any{"code": 400}).
							AddLocalizations(map[string]string{
			"pt-BR": "Você não pode se registrar em uma versão de schema arquivado",
		})
	ErrProjectUserErrorEncodingMetadata = fail.Form(ProjectUserErrorEncodingMetadata, "error encoding project user metadata", false, map[string]any{"code": 500}).
						AddLocalizations(map[string]string{
			"pt-BR": "Erro ao codificar os metadados do projeto do usuário",
		})
	ErrProjectUserNotFromProject = fail.Form(ProjectUserNotFromProject, "project user not from project", false, map[string]any{"code": 500}).
					AddLocalizations(map[string]string{
			"pt-BR": "O Usuário não pertence a esse projeto",
		})
	ErrProjectUserRegisterOnSchemaNoVersion = fail.Form(ProjectUserRegisterOnSchemaNoVersion, "can't register on a schema that has no published version", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "Você não pode se registrar em um schema que não possui versão publicada",
		})

	ErrSchemaNotOwnedByPrincipal = fail.Form(SchemaNotOwnedByPrincipal, "%s", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "%s",
		})
	ErrSchemaNoValidStatus = fail.Form(SchemaNoValidStatus, "CATASTROPHIC: schema found with no valid status: %s", false, map[string]any{"code": 500}).
				AddLocalizations(map[string]string{
			"pt-BR": "CATÁSTROFE: O schema encontrado não possui um status válido: %s",
		})
	ErrSchemaInvalidFlowID = fail.Form(SchemaInvalidFlowID, "invalid flow ID: %s", false, map[string]any{"code": 400}).
				AddLocalizations(map[string]string{
			"pt-BR": "O Flow ID é inválido: %s",
		})
	ErrSchemaFlowIDIsReserved = fail.Form(SchemaFlowIDIsReserved, "flow id can't be the reserved keyword '%s'", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "O Flow ID não podee ser essa palavra reservada '%s'",
		})
	ErrSchemaNoPublishedVersion = fail.Form(SCHEMANoPublishedVersion, "cannot publish a schema with no versions", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar um schema sem versões",
		})
	ErrSchemaFlowIDAlreadyExistsInType = fail.Form(SchemaFlowIDAlreadyExistsInType, "schema with this flow ID already exists in this type", false, map[string]any{"code": 409}).
						AddLocalizations(map[string]string{
			"pt-BR": "O schema com esse flow ID já existe nesse tipo",
		})
	ErrSchemaInvalidSchemaType = fail.Form(SchemaInvalidSchemaType, "invalid schema type", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "O tipo do schema é inválido",
		})
	ErrSchemaHasOnlyDraftVersion = fail.Form(SchemaHasOnlyDraftVersion, "cannot publish a schema with only draft versions", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "Não é possível publicar um schema com apenas versões de rascunhos",
		})
	ErrSchemaHasOnlyArchivedVersion = fail.Form(SchemaHasOnlyArchivedVersion, "cannot publish a schema with only archived versions", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "Não é possível publicar um schema com apenas versões arquivadas",
		})
	ErrSchemaTryingToPublishPublished = fail.Form(SchemaTryingToPublishPublished, "cannot publish a schema that is already published", false, map[string]any{"code": 401}).
						AddLocalizations(map[string]string{
			"pt-BR": "Não é possível publicar um schema que já está publicado",
		})
	ErrSchemaTryingToPublishArchived = fail.Form(SchemaTryingToPublishArchived, "cannot publish a schema that is archived", false, map[string]any{"code": 401}).
						AddLocalizations(map[string]string{
			"pt-BR": "Não é possível publicar um schema que está arquivado",
		})
	ErrSchemaMetadataNotAllowed = fail.Form(SchemaMetadataNotAllowed, "custom fields are not allowed for core schema", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "Os campos personalizados não são permitidos no esquema principal",
		})
	ErrSchemaEmptySchemaType = fail.Form(SchemaEmptySchemaType, "schema type can't be empty", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "O tipo do schema não poder ser vazio",
		})
	ErrSchemaEmptyFlowID = fail.Form(SchemaEmptyFlowID, "flow id can't be empty", false, map[string]any{"code": 401}).
				AddLocalizations(map[string]string{
			"pt-BR": "O Flow ID não pode ser vazio",
		})

	ErrSchemaVersionNotDraft = fail.Form(SchemaVersionNotDraft, "cannot publish a schema version that isn't a draft", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "Não é possível publicar uma versão do schema que não seja um rascunho",
		})
	ErrSchemaVersionDraftAlreadyExists = fail.Form(SCHEMAVersionDraftAlreadyExists, "a draft schema version already exists", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "Já existe um rascunho da versão desse schema",
		})
	ErrSchemaVersionPublishWithNoFields = fail.Form(SchemaVersionPublishWithNoFields, "cannot publish a schema version with no fields", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "Não é possível publicar uma versão do schema com nenhum campo",
		})
	ErrSchemaVersionDraftDoesntExist = fail.Form(SchemaVersionDraftDoesntExist, "cannot publish a schema with a version draft that doesn't exist", false, map[string]any{"code": 401}).
						AddLocalizations(map[string]string{
			"pt-BR": "Não é possível publicar uma versão do schema de rascunho que não existe",
		})
	ErrSchemaVersionTryingToPublishPublished = fail.Form(SchemaVersionTryingToPublishPublished, "cannot publish a schema version that is already published", false, map[string]any{"code": 401}).
							AddLocalizations(map[string]string{
			"pt-BR": "Não é possível publicar uma versão do schema que já está publicada",
		})
	ErrSchemaVersionTryingToPublishArchived = fail.Form(SchemaVersionTryingToPublishArchived, "cannot publish a schema version that is archived", false, map[string]any{"code": 401}).
						AddLocalizations(map[string]string{
			"pt-BR": "Não é possível publicar uma versão do schema que está arquivada",
		})
	ErrSchemaVersionMismatch = fail.Form(SchemaVersionMismatch, "schema version and supplied version mismatch", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "A versão do schema e a versão fornecida não correspondem",
		})
	ErrSchemaVersionNonDraftAddFieldsNotAllowed = fail.Form(SchemaVersionNonDraftAddFieldsNotAllowed, "cannot add fields to a non-draft version", false, map[string]any{"code": 400}).
							AddLocalizations(map[string]string{
			"pt-BR": "Não é possível adicionar campos em uma versão que não seja rascunho",
		})
	ErrSchemaVersionNoValidStatus = fail.Form(SchemaVersionNoValidStatus, "CATASTROPHIC: schema version found with no valid status", false, map[string]any{"code": 401}).
					AddLocalizations(map[string]string{
			"pt-BR": "CATÁSTROFE: A versão do schema encontrada sem status válido",
		})
	ErrSchemaVersionDraftOnNonPublished = fail.Form(SchemaVersionDraftOnNonPublished, "new versions can only be drafted from published versions", false, map[string]any{"code": 400}).
						AddLocalizations(map[string]string{
			"pt-BR": "Novas versões só podem virar rascunhos a partir dee versões publicadas",
		})
	ErrSchemaVersionNoChanges = fail.Form(SchemaVersionNoChanges, "cannot publish a version with no changes", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "Não é possível publicar uma versão sem mudanças",
		})
	ErrSchemaVersionTryingToPublishNonExistant = fail.Form(SchemaVersionTryingToPublishNonExistant, "cannot publish a non-existent schema version", false, map[string]any{"code": 400}).
							AddLocalizations(map[string]string{
			"pt-BR": "Não é possível publicar uma versão inexistente",
		})

	ErrFieldSamePositionForMultipleFields = fail.Form(FIELDSamePositionForMultipleFields, "two fields can't occupy the same position", false, map[string]any{"code": 409}).
						AddLocalizations(map[string]string{
			"pt-BR": "Dois campos não podem ocupar a mesma posição",
		})

	ErrFieldNoAffectedRowsOnClone = fail.Form(FieldNoAffectedRowsOnClone, "no affected rows", false, map[string]any{"code": 404}).
					AddLocalizations(map[string]string{
			"pt-BR": "Nenhuma linha afetada",
		})

	ErrValidationUUIDWasNil = fail.Form(ValidationUUIDWasNil, "%s field is nil", false, map[string]any{"code": 404}).
				AddLocalizations(map[string]string{
			"pt-BR": "%s está nulo",
		})
)
