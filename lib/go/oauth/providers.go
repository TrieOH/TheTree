package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// Credentials are the per-deployment client credentials for a provider.
// For IdentityX these come from a project's configured provider row or,
// for platform-level logins, from the environment.
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

// redirectFromEnv resolves a provider's redirect URI from the environment.
// The redirect URI is registered in the provider console as an exact URL and
// is the same for every scope (the callback path never carries the project).
func redirectFromEnv(key string) string {
	return os.Getenv(key)
}

// EnvCredentials returns the platform-level credentials for a provider,
// ok=false when they are not configured in the environment. Used for
// IdentityX itself (no project scope).
func EnvCredentials(provider string) (Credentials, bool) {
	switch provider {
	case "google":
		creds := Credentials{ClientID: os.Getenv("GOOGLE_CLIENT_ID"), ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET")}
		return creds, creds.ClientID != "" && creds.ClientSecret != ""
	case "github":
		creds := Credentials{ClientID: os.Getenv("GITHUB_CLIENT_ID"), ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET")}
		return creds, creds.ClientID != "" && creds.ClientSecret != ""
	default:
		return Credentials{}, false
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

// Registry holds the supported providers. Credentials are deliberately not
// here: callers resolve them per scope (project row or env) and build the
// config via Provider.Config.
var Registry = map[string]Provider{
	"google": {
		Endpoint:    google.Endpoint,
		Scopes:      []string{"email", "profile"},
		Userinfo:    "https://www.googleapis.com/oauth2/v2/userinfo",
		RedirectURL: redirectFromEnv("GOOGLE_REDIRECT_URI"),
	},
	"github": {
		Endpoint:    github.Endpoint,
		Scopes:      []string{"user:email"},
		Userinfo:    "https://api.github.com/user",
		RedirectURL: redirectFromEnv("GITHUB_REDIRECT_URI"),
	},
}

type GitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func FetchGitHubEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var emails []GitHubEmail
	err = json.NewDecoder(resp.Body).Decode(&emails)
	if err != nil {
		return "", err
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	return "", errors.New("no verified primary email found")
}
