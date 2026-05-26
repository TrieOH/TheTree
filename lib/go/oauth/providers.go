package oauth

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Provider struct {
	Config   *oauth2.Config
	Userinfo string
}

type UserInfo struct {
	Sub   string `json:"id"`
	Email string `json:"email"`
}

var Registry = map[string]Provider{
	"google": {
		Config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
		Userinfo: "https://www.googleapis.com/oauth2/v2/userinfo",
	},
}
