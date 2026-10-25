package errs

import "errors"

// One error per named constraint in the schema. internal/app's constraint
// registry maps each constraint to its error, and lib/database joins it
// with the violation's class (unique, foreign key, check). Boot fails when
// the schema has a constraint the registry does not list.

// actors
var (
	ErrInvalidActorType      = errors.New("actor type must be one of: human, service, machine")
	ErrInvalidAuthMethod     = errors.New("auth method must be one of: api_key, password, google, github")
	ErrHumanEmailRequired    = errors.New("email is required for human actors")
	ErrEmailTakenForScope    = errors.New("an account with this email already exists for this scope and auth method")
	ErrInvalidExternalIDP    = errors.New("external identity provider must be one of: google, github")
	ErrExternalIdentityTaken = errors.New("this external identity is already linked to this account")
)

// organizations and projects
var (
	ErrOrgSlugTaken            = errors.New("an organization with this slug already exists")
	ErrInvalidOrgRole          = errors.New("organization role must be one of: owner, admin, member")
	ErrProjectSlugTaken        = errors.New("a project with this slug already exists")
	ErrDomainTaken             = errors.New("this domain is already verified by another project")
	ErrInvalidBrandSlugFormat  = errors.New("brand slug must contain only lowercase letters")
	ErrInvalidBrandSlugLength  = errors.New("brand slug length must be between 3 and 32 characters")
	ErrInvalidProjectRole      = errors.New("project role must be one of: owner, admin, member")
	ErrInvalidPlatformRole     = errors.New("platform role must be one of: super_admin, admin, support")
	ErrInvalidOAuthProvider    = errors.New("OAuth provider must be one of: google, github")
	ErrOAuthProviderConfigured = errors.New("this OAuth provider is already configured for this project")
	ErrDuplicateOAuthState     = errors.New("duplicate OAuth login state")
)

// keys and capabilities
var (
	ErrAPIKeyPrefixTaken    = errors.New("an API key with this display prefix already exists")
	ErrCapabilityExists     = errors.New("a capability with this resource and action already exists for this scope")
	ErrInvalidCryptoKeyType = errors.New("crypto key type must be one of: encryption, signing")
	ErrInvalidCryptoKeyStat = errors.New("crypto key status must be one of: active, retiring, retired, revoked")
)

// profiles
var (
	ErrProjectProfileSchemaExists  = errors.New("a profile schema already exists for this project")
	ErrProfileSchemaVersionExists  = errors.New("this profile schema version already exists for the scope")
	ErrPlatformSchemaVersionExists = errors.New("a platform profile schema version already exists")
	ErrHandleTaken                 = errors.New("this handle is already in use by another profile")
)

// emails, terms of service, tokens and blacklist
var (
	ErrInvalidTemplateKindColumn = errors.New("email template kind must be one of: verify, reset")
	ErrEmailTemplateExists       = errors.New("this project already has an email template of this kind")
	ErrTosContentEmpty           = errors.New("terms of service content must not be empty")
	ErrInvalidTosVersion         = errors.New("terms of service version must be at least 1")
	ErrTosAlreadyExists          = errors.New("this project already has terms of service")
	ErrInvalidAcceptanceSource   = errors.New("terms acceptance source must be clickwrap or continued_use")
	ErrTosAlreadyAccepted        = errors.New("this actor has already accepted this terms version")
	ErrInvalidActionTokenPurpose = errors.New("action token purpose must be one of: email_verify, password_reset")
	ErrInvalidBlacklistType      = errors.New("blacklist entry type must be one of: actor, token, api_key, email, ip")
	ErrAlreadyBlacklisted        = errors.New("this target is already blacklisted for this scope")
)
