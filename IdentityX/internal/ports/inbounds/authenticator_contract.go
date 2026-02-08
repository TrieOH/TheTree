package inbounds

type AuthenticateRequestInput struct {
	AccessToken  string
	RefreshToken string
	ApiKey       string
	Issuer       string
}
