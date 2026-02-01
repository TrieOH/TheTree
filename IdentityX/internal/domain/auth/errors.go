package auth

type ErrTokenMissingKID struct {
	TokenType string
}

func (e ErrTokenMissingKID) Error() string {
	return e.TokenType + " token missing kid"
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

type ErrTokenInvalidFormat struct {
	TokenType string
}

func (e ErrTokenInvalidFormat) Error() string {
	return "invalid " + e.TokenType + " token format"
}
