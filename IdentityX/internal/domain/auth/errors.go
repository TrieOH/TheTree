package auth

type ErrTokenMissingKID struct {
	TokenType string
}

func (e ErrTokenMissingKID) Error() string {
	return e.TokenType + " token missing kid"
}

type ErrInvalidToken struct {
	TokenType string
}

func (e ErrInvalidToken) Error() string {
	return "invalid " + e.TokenType + " token"
}

type ErrTokenInvalidKID struct {
	TokenType string
}

func (e ErrTokenInvalidKID) Error() string {
	return "invalid " + e.TokenType + " token kid"
}

type ErrTokenUnknownKID struct {
	TokenType string
}

func (e ErrTokenUnknownKID) Error() string {
	return "unknown " + e.TokenType + " token kid"
}

type ErrSigningToken struct {
	TokenType string
	Cause     error
}

func (e ErrSigningToken) Error() string {
	return "error signing " + e.TokenType + " token"
}

type ErrTokenInvalidAlg struct {
	TokenType string
}

func (e ErrTokenInvalidAlg) Error() string {
	return "invalid " + e.TokenType + " token alg"
}

type ErrTokenInvalidFormat struct {
	TokenType string
}

func (e ErrTokenInvalidFormat) Error() string {
	return "invalid " + e.TokenType + " token format"
}

type ErrTokenUntrusted struct {
	TokenType string
}

func (e ErrTokenUntrusted) Error() string {
	return "untrusted " + e.TokenType + " token"
}
