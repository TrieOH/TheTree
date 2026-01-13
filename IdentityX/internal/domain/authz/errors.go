package authz

type ErrMissingPrincipal struct{}

func (ErrMissingPrincipal) Error() string {
	return "missing principal"
}

type ErrPrincipalMissingInContext struct{}

func (ErrPrincipalMissingInContext) Error() string {
	return "missing principal in context"
}

type ErrMissingAccessClaims struct{}

func (ErrMissingAccessClaims) Error() string {
	return "missing access claims"
}

type ErrMissingRefreshClaims struct{}

func (ErrMissingRefreshClaims) Error() string {
	return "missing refresh claims"
}

type ErrInvalidAccessJTI struct {
	Cause error
}

func (e ErrInvalidAccessJTI) Error() string {
	return "couldn't parse access token JTI"
}

type ErrInvalidRefreshJTI struct {
	Cause error
}

func (e ErrInvalidRefreshJTI) Error() string {
	return "couldn't parse refresh token JTI"
}
