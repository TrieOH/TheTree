package inbounds

type TokenIssuer interface {
	NewAccessToken(in NewAccessTokenInput) (string, error)
	NewRefreshToken(in NewRefreshTokenInput) (string, error)
	NewProjectAccessToken(in NewProjectAccessTokenInput) (string, error)
	NewVerificationToken(in NewVerificationTokenInput) (string, error)
}
