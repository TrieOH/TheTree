package dto

type BeginOAuthRequest struct {
	FinalRedirectURL string `json:"final_redirect_url" validate:"required,url"`
}

type BeginOAuthResponse struct {
	RedirectURL string `json:"redirect_url"`
}
