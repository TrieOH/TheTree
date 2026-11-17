package oauth

import (
	"encoding/json"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// Credentials are the per-deployment client credentials for a provider.
// For IdentityX these come from a project's configured provider row or,
// for platform-level logins, from the service config.
type Credentials struct {
	ClientID     string
	ClientSecret string
	// RedirectURL overrides the provider's static redirect URI when the
	// scope registers its own (e.g. a project's callback_url). Empty means
	// the provider-level value is used.
	RedirectURL string
}

// Provider is a factory for oauth2.Config values. It holds the static
// provider metadata (endpoints, scopes, userinfo URL, redirect URL) and
// mounts a config from the credentials of the caller's scope — never from
// package-level state.
type Provider struct {
	Endpoint    oauth2.Endpoint
	Scopes      []string
	Userinfo    string
	RedirectURL string
}

func (p Provider) Config(creds Credentials) *oauth2.Config {
	redirectURL := p.RedirectURL
	if creds.RedirectURL != "" {
		redirectURL = creds.RedirectURL
	}
	return &oauth2.Config{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       p.Scopes,
		Endpoint:     p.Endpoint,
	}
}

type UserInfo struct {
	Sub   json.RawMessage `json:"id"`
	Email string          `json:"email"`
}

func (u UserInfo) SubString() string {
	s := strings.Trim(string(u.Sub), "\"")
	return s
}

// NewRegistry returns the supported providers. Credentials are deliberately
// not here: callers resolve them per scope (project row or platform config)
// and build the config via Provider.Config. redirects maps a provider name to
// the redirect URI registered in its console; it is the same for every scope
// (the callback path never carries the project).
func NewRegistry(redirects map[string]string) map[string]Provider {
	return map[string]Provider{
		"google": {
			Endpoint:    google.Endpoint,
			Scopes:      []string{"email", "profile"},
			Userinfo:    "https://www.googleapis.com/oauth2/v2/userinfo",
			RedirectURL: redirects["google"],
		},
		"github": {
			Endpoint:    github.Endpoint,
			Scopes:      []string{"user:email"},
			Userinfo:    "https://api.github.com/user",
			RedirectURL: redirects["github"],
		},
	}
}

// GitHubEmail is one entry of GitHub's /user/emails response.
type GitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}
