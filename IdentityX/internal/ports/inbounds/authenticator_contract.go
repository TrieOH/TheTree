package inbounds

type AuthenticateRequestInput struct {
	AccessToken  string
	RefreshToken string
	Issuer       string
}

type ErrInvalidIssuer struct {
	TokenType string
}

func (e ErrInvalidIssuer) Error() string {
	return "invalid " + e.TokenType + " token issuer"
}

type ErrTokenIDMismatch struct{}

func (e ErrTokenIDMismatch) Error() string {
	return "access token does not belong to this refresh token"
}

type ErrTokenSessionMismatch struct{}

func (e ErrTokenSessionMismatch) Error() string {
	return "token/session mismatch"
}

type ErrAuthSessionRevoked struct{}

func (e ErrAuthSessionRevoked) Error() string {
	return "session not found or revoked"
}
