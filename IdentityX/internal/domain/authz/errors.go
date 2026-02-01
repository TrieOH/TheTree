package authz

type ErrMissingAccessClaims struct{}

func (ErrMissingAccessClaims) Error() string {
	return "missing access claims"
}

type ErrMissingRefreshClaims struct{}

func (ErrMissingRefreshClaims) Error() string {
	return "missing refresh claims"
}
