package inbounds

type TokenIssuer interface {
	NewAccessToken(in NewAccessTokenInput) ([]byte, error)
	NewRefreshToken(in NewRefreshTokenInput) ([]byte, error)
	NewProjectAccessToken(in NewProjectAccessTokenInput) ([]byte, error)
	NewVerificationToken(in NewVerificationTokenInput) ([]byte, error)
	AssembleJWT(payload []byte, sig []byte) string
}
