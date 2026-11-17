package errs

import "errors"

// actors
var (
	ErrInvalidActorType   = errors.New("actor type must be one of: human, service, machine")
	ErrInvalidAuthMethod  = errors.New("auth method must be one of: api_key, password, google, github")
	ErrHumanEmailRequired = errors.New("email is required for human actors")
	ErrInvalidExternalIDP = errors.New("external identity provider must be one of: google, github")
)

// organizations and projects
var (
	ErrInvalidOrgRole         = errors.New("organization role must be one of: owner, admin, member")
	ErrInvalidBrandSlugFormat = errors.New("brand slug must contain only lowercase letters")
	ErrInvalidBrandSlugLength = errors.New("brand slug length must be between 3 and 32 characters")
	ErrInvalidProjectRole     = errors.New("project role must be one of: owner, admin, member")
	ErrInvalidPlatformRole    = errors.New("platform role must be one of: super_admin, admin, support")
	ErrInvalidOAuthProvider   = errors.New("OAuth provider must be one of: google, github")
)

// keys and capabilities
var (
	ErrInvalidCryptoKeyType = errors.New("crypto key type must be one of: encryption, signing")
	ErrInvalidCryptoKeyStat = errors.New("crypto key status must be one of: active, retiring, retired, revoked")
)

// emails, terms of service, tokens and blacklist
var (
	ErrInvalidTemplateKindColumn = errors.New("email template kind must be one of: verify, reset")
	ErrTosContentEmpty           = errors.New("terms of service content must not be empty")
	ErrInvalidTosVersion         = errors.New("terms of service version must be at least 1")
	ErrInvalidAcceptanceSource   = errors.New("terms acceptance source must be clickwrap or continued_use")
	ErrInvalidActionTokenPurpose = errors.New("action token purpose must be one of: email_verify, password_reset")
	ErrInvalidBlacklistType      = errors.New("blacklist entry type must be one of: actor, token, api_key, email, ip")
)
