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
	ErrRequestEmptyCookie = fail.Form(RequestEmptyCookie, "empty %s cookie value", false, map[string]any{"code": 400}, "UNDEFINED").
				AddLocalizations(map[string]string{
			"pt-BR": "O valor do cookie %s está vazio",
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
	ErrRequestUnknownQueryParam = fail.Form(RequestUnknownQueryParam, "unknown query parameter: %s", false, map[string]any{"code": 400}, "UNDEFINED").
					AddLocalizations(map[string]string{
			"pt-BR": "parâmetro de consulta desconhecido: %s",
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

	ErrSchemaNoPublishedVersion = fail.Form(SCHEMANoPublishedVersion, "cannot publish a schema with no versions", false, map[string]any{"code": 400}).
					AddLocalizations(map[string]string{
			"pt-BR": "não é possível publicar um schema sem versões",
		})
)
