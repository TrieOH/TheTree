package inbounds

type AuthenticateRequestInput struct {
	AccessToken  string
	RefreshToken string
	Issuer       string
}
